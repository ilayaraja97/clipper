package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/ilayaraja97/clipper/ai"
	"github.com/ilayaraja97/clipper/config"
	"github.com/ilayaraja97/clipper/history"
	"github.com/ilayaraja97/clipper/logger"
	"github.com/ilayaraja97/clipper/run"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/viper"
)

type execErrorMsg struct {
	err error
}

type chatErrorMsg struct {
	err error
}

type parallelFallbackMsg struct {
	runOut    run.RunOutput
	engineOut ai.EngineExecOutput
}

// commandFinishedMsg reports a confirmed command that ran via tea.ExecProcess.
// Its output went straight to the terminal, so only the exit status travels here.
type commandFinishedMsg struct {
	command string
	err     error
}

// settingsReloadedMsg carries the config and engine rebuilt after the user
// edited settings. They are built off the main loop, so they arrive as a
// message rather than being assigned to the Ui directly.
type settingsReloadedMsg struct {
	config *config.Config
	engine *ai.Engine
	err    error
}

type UiState struct {
	error       error
	runMode     RunMode
	promptMode  PromptMode
	configuring bool
	querying    bool
	confirming  bool
	executing   bool
	args        string
	pipe        string
	buffer      string
	command     string
}

type UiDimensions struct {
	width  int
	height int
}

type UiComponents struct {
	prompt   *Prompt
	renderer *Renderer
	spinner  *Spinner
}

type Ui struct {
	state      UiState
	dimensions UiDimensions
	components UiComponents
	configFlow configFlow
	config     *config.Config
	engine     *ai.Engine
	history    *history.History
}

func NewUi(input *UiInput) *Ui {
	return &Ui{
		state: UiState{
			error:       nil,
			runMode:     input.GetRunMode(),
			promptMode:  input.GetPromptMode(),
			configuring: false,
			querying:    false,
			confirming:  false,
			executing:   false,
			args:        input.GetArgs(),
			pipe:        input.GetPipe(),
			buffer:      "",
			command:     "",
		},
		dimensions: UiDimensions{
			150,
			150,
		},
		components: UiComponents{
			prompt: NewPrompt(input.GetPromptMode()),
			renderer: NewRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(0),
			),
			spinner: NewSpinner(),
		},
		configFlow: newConfigFlow(),
		history:    history.NewHistory(),
	}
}

func (u *Ui) Init() tea.Cmd {
	logger.Log.Debug().Str("runMode", u.state.runMode.String()).Str("promptMode", u.state.promptMode.String()).Msg("initializing UI")

	config, err := config.NewConfig()
	if err != nil {
		logger.Log.Debug().Err(err).Msg("config error")
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Log.Info().Msg("no config file found, starting configuration")
			if u.state.runMode == ReplMode {
				return tea.Sequence(
					tea.ClearScreen,
					u.startConfig(),
				)
			} else {
				return u.startConfig()
			}
		} else {
			logger.Log.Error().Err(err).Msg("failed to load config")
			return tea.Sequence(
				tea.Println(u.components.renderer.RenderError(err.Error())),
				tea.Quit,
			)
		}
	}

	logger.Log.Debug().Str("runMode", u.state.runMode.String()).Msg("config loaded successfully")

	if u.state.runMode == ReplMode {
		return u.startRepl(config)
	} else {
		return u.startCli(config)
	}
}

func (u *Ui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if u.state.querying {
			spinner, cmd := u.components.spinner.Update(msg)
			u.components.spinner = spinner
			return u, cmd
		}

	case tea.WindowSizeMsg:
		u.dimensions.width = msg.Width
		u.dimensions.height = msg.Height

	case tea.KeyMsg:
		return u.handleKeyPress(msg)

	case ai.EngineExecOutput:
		var output string
		if msg.IsExecutable() {
			u.state.confirming = true
			u.state.command = msg.GetCommand()
			output = u.components.renderer.RenderContent(fmt.Sprintf("`%s`", u.state.command))
			output += fmt.Sprintf("  %s\n\n  confirm execution? [y/N]", u.components.renderer.RenderHelp(msg.GetExplanation()))
			u.components.prompt.Blur()
		} else {
			output = u.components.renderer.RenderContent(msg.GetExplanation())
			u.components.prompt.Focus()
			if u.state.runMode == CliMode {
				return u, tea.Sequence(tea.Println(output), tea.Quit)
			}
		}
		prompt, cmd := u.components.prompt.Update(msg)
		u.components.prompt = prompt
		return u, tea.Sequence(cmd, textinput.Blink, tea.Println(output))

	case parallelFallbackMsg:
		u.state.querying = false
		prompt, cmd := u.components.prompt.Update(msg)
		u.components.prompt = prompt

		var output string
		runOutput := strings.TrimSpace(msg.runOut.GetContent())
		if runOutput != "" {
			output += u.components.renderer.RenderContent(runOutput) + "\n"
		}
		if msg.runOut.HasError() {
			errOutput := u.components.renderer.RenderError(fmt.Sprintf("\n%s\n", msg.runOut.GetErrorMessage()))
			output += errOutput + "\n"
		}

		if msg.engineOut.IsExecutable() {
			u.state.confirming = true
			u.state.command = msg.engineOut.GetCommand()
			output += u.components.renderer.RenderContent(fmt.Sprintf("`%s`", u.state.command))
			output += fmt.Sprintf("  %s\n\n  confirm execution? [y/N]", u.components.renderer.RenderHelp(msg.engineOut.GetExplanation()))
			u.components.prompt.Blur()
		} else {
			output += u.components.renderer.RenderContent(msg.engineOut.GetExplanation())
			u.components.prompt.Focus()
			if u.state.runMode == CliMode {
				return u, tea.Sequence(tea.Println(output), tea.Quit)
			}
		}
		return u, tea.Sequence(cmd, textinput.Blink, tea.Println(output))

	case commandFinishedMsg:
		u.state.executing = false
		u.state.command = ""
		u.components.prompt.Focus()

		if msg.err != nil {
			logger.Log.Error().Err(msg.err).Str("command", msg.command).Msg("command execution failed")
		} else {
			logger.Log.Info().Str("command", msg.command).Msg("command executed successfully")
		}

		if u.state.runMode == ReplMode {
			result := "succeeded"
			if msg.err != nil {
				result = fmt.Sprintf("failed (%s)", msg.err.Error())
			}
			u.engine.AppendAssistantMessage(
				fmt.Sprintf("Command: %s\nResult: %s", msg.command, result),
			)
		}

		var output string
		if msg.err != nil {
			output = u.components.renderer.RenderError(fmt.Sprintf("\n[error]: %s\n", msg.err.Error()))
		} else {
			output = u.components.renderer.RenderSuccess("\n[ok]\n")
		}

		if u.state.runMode == CliMode {
			return u, tea.Sequence(tea.Println(output), tea.Quit)
		}
		return u, tea.Sequence(tea.Println(output), textinput.Blink)

	case settingsReloadedMsg:
		u.state.executing = false
		u.state.command = ""
		u.components.prompt.Focus()

		if msg.err != nil {
			errOutput := u.components.renderer.RenderError(
				fmt.Sprintf("\n[settings error]: %s\n", msg.err.Error()),
			)
			if u.state.runMode == CliMode {
				return u, tea.Sequence(tea.Println(errOutput), tea.Quit)
			}
			return u, tea.Sequence(tea.Println(errOutput), textinput.Blink)
		}

		u.config = msg.config
		u.engine = msg.engine

		output := u.components.renderer.RenderSuccess("\n[settings ok]\n")
		if u.state.runMode == CliMode {
			return u, tea.Sequence(tea.Println(output), tea.Quit)
		}
		return u, tea.Sequence(tea.Println(output), textinput.Blink)

	case run.RunOutput:
		u.state.querying = false
		prompt, cmd := u.components.prompt.Update(msg)
		u.components.prompt = prompt
		u.components.prompt.Focus()
		output := strings.TrimSpace(msg.GetContent())
		if output != "" {
			output = u.components.renderer.RenderContent(output)
		}
		if msg.HasError() {
			errOutput := u.components.renderer.RenderError(fmt.Sprintf("\n%s\n", msg.GetErrorMessage()))
			if output != "" {
				output = fmt.Sprintf("%s\n%s", output, errOutput)
			} else {
				output = errOutput
			}
		} else if msg.GetSuccessMessage() != "" {
			successOutput := u.components.renderer.RenderSuccess(fmt.Sprintf("\n%s\n", msg.GetSuccessMessage()))
			if output != "" {
				output = fmt.Sprintf("%s\n%s", output, successOutput)
			} else {
				output = successOutput
			}
		}
		if u.state.runMode == CliMode {
			return u, tea.Sequence(tea.Println(output), tea.Quit)
		}
		return u, tea.Sequence(tea.Println(output), cmd, textinput.Blink)

	case execErrorMsg:
		u.state.querying = false
		u.components.prompt.Focus()
		errOutput := u.components.renderer.RenderError(fmt.Sprintf("\nexec error: %s\n", msg.err.Error()))
		if u.state.runMode == CliMode {
			return u, tea.Sequence(tea.Println(errOutput), tea.Quit)
		}
		return u, tea.Sequence(tea.Println(errOutput), textinput.Blink)

	case chatErrorMsg:
		u.state.querying = false
		u.components.prompt.Focus()
		errOutput := u.components.renderer.RenderError(fmt.Sprintf("\nchat error: %s\n", msg.err.Error()))
		if u.state.runMode == CliMode {
			return u, tea.Sequence(tea.Println(errOutput), tea.Quit)
		}
		return u, tea.Sequence(tea.Println(errOutput), textinput.Blink)

	case error:
		u.state.error = msg
	}

	return u, nil
}

func (u *Ui) View() string {
	if u.state.error != nil {
		if u.components.renderer != nil {
			return u.components.renderer.RenderError(fmt.Sprintf("[error] %s", u.state.error))
		}
		return fmt.Sprintf("[error] %s\n", u.state.error)
	}

	if u.state.configuring {
		if u.components.renderer != nil {
			return fmt.Sprintf(
				"%s\n%s",
				u.components.renderer.RenderContent(u.state.buffer),
				u.components.prompt.View(),
			)
		}
		return u.components.prompt.View()
	}

	if !u.state.querying && !u.state.confirming && !u.state.executing {
		return u.components.prompt.View()
	}

	if u.state.promptMode == ChatPromptMode {
		if u.state.querying {
			return u.components.spinner.View()
		}
		if u.components.renderer != nil {
			return u.components.renderer.RenderContent(u.state.buffer)
		}
		return u.state.buffer
	} else {
		if u.state.querying {
			return u.components.spinner.View()
		} else {
			if u.state.executing {
				if u.components.renderer != nil {
					return u.components.renderer.RenderHelp("\n  executing command...")
				}
				return "\n  executing command...\n"
			}

			if u.components.renderer != nil {
				return u.components.renderer.RenderContent(u.state.buffer)
			}
			return u.state.buffer
		}
	}
}

func (u *Ui) startRepl(config *config.Config) tea.Cmd {
	u.config = config

	if u.state.promptMode == DefaultPromptMode {
		u.state.promptMode = GetPromptModeFromString(config.GetUserConfig().GetDefaultPromptMode())
	}

	engineMode := ai.ExecEngineMode
	if u.state.promptMode == ChatPromptMode {
		engineMode = ai.ChatEngineMode
	}

	engine, err := ai.NewEngine(engineMode, config)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to create engine in REPL mode")
		u.state.error = err
		return nil
	}

	if u.state.pipe != "" {
		engine.SetPipe(u.state.pipe)
	}

	u.engine = engine
	u.state.buffer = "Welcome \n\n"
	u.state.command = ""
	u.components.prompt = NewPrompt(u.state.promptMode)

	logger.Log.Info().Str("mode", engineMode.String()).Str("pipe", u.state.pipe).Msg("REPL started")

	return tea.Sequence(
		tea.ClearScreen,
		tea.Println(u.components.renderer.RenderContent(u.components.renderer.RenderHelpMessage())),
		textinput.Blink,
	)
}

func (u *Ui) startCli(config *config.Config) tea.Cmd {
	logger.Log.Debug().Str("args", u.state.args).Str("pipe", u.state.pipe).Msg("starting CLI mode")

	u.config = config

	if u.state.promptMode == DefaultPromptMode {
		u.state.promptMode = GetPromptModeFromString(config.GetUserConfig().GetDefaultPromptMode())
	}

	engineMode := ai.ExecEngineMode
	if u.state.promptMode == ChatPromptMode {
		engineMode = ai.ChatEngineMode
	}

	engine, err := ai.NewEngine(engineMode, config)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to create engine in CLI mode")
		u.state.error = err
		return nil
	}

	if u.state.pipe != "" {
		engine.SetPipe(u.state.pipe)
	}

	u.engine = engine
	u.state.querying = true
	u.state.confirming = false
	u.state.buffer = ""
	u.state.command = ""

	logger.Log.Info().Str("mode", engineMode.String()).Str("args", u.state.args).Msg("CLI mode started")

	if u.state.promptMode == ExecPromptMode {
		if run.RequiresTTY(u.state.args) {
			return u.execCommand(u.state.args)
		}
		return tea.Batch(
			u.components.spinner.Tick,
			u.startExec(u.state.args),
		)
	} else {
		return tea.Batch(
			u.components.spinner.Tick,
			u.startChat(u.state.args),
		)
	}
}

func (u *Ui) startConfig() tea.Cmd {
	u.state.configuring = true
	u.state.querying = false
	u.state.confirming = false
	u.state.executing = false

	u.state.command = ""
	u.configFlow = newConfigFlow()
	u.components.prompt = NewPrompt(ConfigPromptMode)
	u.refreshConfigScreen()

	return textinput.Blink
}

func (u *Ui) finishConfig() tea.Cmd {
	u.state.configuring = false

	config, err := config.WriteConfig(u.configFlow.Input(), true)
	if err != nil {
		u.state.configuring = true
		u.refreshConfigScreen()
		u.state.buffer += fmt.Sprintf("\n\n%s", u.components.renderer.RenderWarning(err.Error()))

		return textinput.Blink
	}

	u.config = config
	engineMode := ai.ExecEngineMode
	if u.state.promptMode == ChatPromptMode {
		engineMode = ai.ChatEngineMode
	}

	engine, err := ai.NewEngine(engineMode, config)
	if err != nil {
		u.state.error = err
		return nil
	}

	if u.state.pipe != "" {
		engine.SetPipe(u.state.pipe)
	}

	u.engine = engine

	if u.state.runMode == ReplMode {
		nextPromptMode := resolvePromptMode(u.state.promptMode, config)
		u.state.buffer = ""
		u.state.command = ""
		u.state.promptMode = nextPromptMode
		u.components.prompt = NewPrompt(nextPromptMode)

		return tea.Sequence(
			tea.ClearScreen,
			tea.Println(u.components.renderer.RenderSuccess("\n[settings ok]\n")),
			textinput.Blink,
		)
	} else {
		if u.state.promptMode == ExecPromptMode {
			u.state.configuring = false
			u.state.buffer = ""
			if run.RequiresTTY(u.state.args) {
				return tea.Sequence(
					tea.Println(u.components.renderer.RenderSuccess("\n[settings ok]")),
					u.execCommand(u.state.args),
				)
			}
			u.state.querying = true
			return tea.Sequence(
				tea.Println(u.components.renderer.RenderSuccess("\n[settings ok]")),
				u.components.spinner.Tick,
				u.startExec(u.state.args),
			)
		} else {
			u.state.querying = true
			return tea.Batch(
				u.components.spinner.Tick,
				u.startChat(u.state.args),
			)
		}
	}
}

func resolvePromptMode(current PromptMode, config *config.Config) PromptMode {
	if current != DefaultPromptMode {
		return current
	}

	return GetPromptModeFromString(config.GetUserConfig().GetDefaultPromptMode())
}

func (u *Ui) advanceConfig() tea.Cmd {
	u.configFlow.SetCurrentValue(u.components.prompt.GetValue())
	if u.configFlow.Next() {
		u.refreshConfigScreen()

		return textinput.Blink
	}

	return u.finishConfig()
}

func (u *Ui) refreshConfigScreen() {
	field := u.configFlow.CurrentField()
	u.state.buffer = u.components.renderer.RenderConfigMessage(
		field.label,
		field.description,
		u.configFlow.ProgressLines(),
	)
	u.configFlow.ApplyToPrompt(u.components.prompt)
}

func (u *Ui) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return u, tea.Quit

	case tea.KeyUp, tea.KeyDown:
		if u.state.configuring {
			u.configFlow.SetCurrentValue(u.components.prompt.GetValue())
			if msg.Type == tea.KeyUp {
				u.configFlow.Move(-1)
			} else {
				u.configFlow.Move(1)
			}
			u.refreshConfigScreen()
			return u, textinput.Blink
		}
		if !u.state.querying && !u.state.confirming {
			var input *string
			if msg.Type == tea.KeyUp {
				input = u.history.GetPrevious()
			} else {
				input = u.history.GetNext()
			}
			if input != nil {
				u.components.prompt.SetValue(*input)
				prompt, cmd := u.components.prompt.Update(msg)
				u.components.prompt = prompt
				return u, cmd
			}
		}

	case tea.KeyTab:
		if !u.state.querying && !u.state.confirming && !u.state.configuring {
			if u.state.promptMode == ChatPromptMode {
				u.state.promptMode = ExecPromptMode
				u.components.prompt.SetMode(ExecPromptMode)
				u.engine.SetMode(ai.ExecEngineMode)
			} else {
				u.state.promptMode = ChatPromptMode
				u.components.prompt.SetMode(ChatPromptMode)
				u.engine.SetMode(ai.ChatEngineMode)
			}
			prompt, cmd := u.components.prompt.Update(msg)
			u.components.prompt = prompt
			return u, tea.Batch(cmd, textinput.Blink)
		}

	case tea.KeyEnter:
		if u.state.configuring {
			return u, u.advanceConfig()
		}
		if !u.state.querying && !u.state.confirming {
			input := u.components.prompt.GetValue()
			if input != "" {
				inputPrint := u.components.prompt.AsString()
				u.history.Add(input)
				u.components.prompt.SetValue("")
				u.components.prompt.Blur()
				prompt, cmd := u.components.prompt.Update(msg)
				u.components.prompt = prompt
				if u.state.promptMode == ChatPromptMode {
					u.state.querying = true
					return u, tea.Batch(
						cmd,
						tea.Println(inputPrint),
						u.components.spinner.Tick,
						u.startChat(input),
					)
				}
				if run.RequiresTTY(input) {
					return u, tea.Sequence(
						cmd,
						tea.Println(inputPrint),
						u.execCommand(input),
					)
				}
				u.state.querying = true
				return u, tea.Batch(
					cmd,
					tea.Println(inputPrint),
					u.startExec(input),
					u.components.spinner.Tick,
				)
			}
		}

	case tea.KeyCtrlH:
		if !u.state.configuring && !u.state.querying && !u.state.confirming {
			prompt, cmd := u.components.prompt.Update(msg)
			u.components.prompt = prompt
			return u, tea.Batch(
				cmd,
				tea.Println(u.components.renderer.RenderContent(u.components.renderer.RenderHelpMessage())),
				textinput.Blink,
			)
		}

	case tea.KeyCtrlL:
		if !u.state.querying && !u.state.confirming {
			prompt, cmd := u.components.prompt.Update(msg)
			u.components.prompt = prompt
			return u, tea.Batch(cmd, tea.ClearScreen, textinput.Blink)
		}

	case tea.KeyCtrlR:
		if !u.state.querying && !u.state.confirming {
			u.history.Reset()
			u.engine.Reset()
			u.state.buffer = ""
			u.state.command = ""
			u.components.prompt.SetValue("")
			prompt, cmd := u.components.prompt.Update(msg)
			u.components.prompt = prompt
			return u, tea.Batch(
				cmd,
				tea.ClearScreen,
				tea.Println(u.components.renderer.RenderContent(u.components.renderer.RenderHelpMessage())),
				textinput.Blink,
			)
		}

	case tea.KeyCtrlS:
		if !u.state.querying && !u.state.confirming && !u.state.configuring && !u.state.executing {
			u.state.executing = true
			u.state.buffer = ""
			u.state.command = ""
			u.components.prompt.Blur()
			prompt, cmd := u.components.prompt.Update(msg)
			u.components.prompt = prompt
			return u, tea.Batch(cmd, u.editSettings())
		}
	}
	if u.state.confirming {
		if strings.ToLower(msg.String()) == "y" {
			u.state.confirming = false
			u.state.executing = true
			u.state.buffer = ""
			u.components.prompt.SetValue("")
			return u, u.execCommand(u.state.command)
		}
		u.state.confirming = false
		u.state.executing = false
		u.state.buffer = ""
		prompt, cmd := u.components.prompt.Update(msg)
		u.components.prompt = prompt
		u.components.prompt.SetValue("")
		u.components.prompt.Focus()
		if u.state.runMode == ReplMode {
			return u, tea.Batch(
				cmd,
				tea.Println(fmt.Sprintf("\n%s\n", u.components.renderer.RenderWarning("[cancel]"))),
				textinput.Blink,
			)
		}
		return u, tea.Sequence(
			cmd,
			tea.Println(fmt.Sprintf("\n%s\n", u.components.renderer.RenderWarning("[cancel]"))),
			tea.Quit,
		)
	}
	u.components.prompt.Focus()
	prompt, cmd := u.components.prompt.Update(msg)
	u.components.prompt = prompt
	return u, tea.Batch(cmd, textinput.Blink)
}

func (u *Ui) startExec(input string) tea.Cmd {
	logger.Log.Debug().Str("input", input).Msg("starting exec")

	// Mutate state here, on the main loop. The returned tea.Cmd runs on its own
	// goroutine and must not touch u.* -- it may only read the locals captured below.
	u.state.confirming = false
	u.state.buffer = ""
	u.state.command = ""

	shell := u.config.GetSystemConfig().GetShell()
	engine := u.engine

	return func() tea.Msg {
		type engineResult struct {
			out *ai.EngineExecOutput
			err error
		}
		engineCh := make(chan engineResult, 1)

		go func() {
			out, err := engine.ExecCompletion(context.Background(), input)
			engineCh <- engineResult{out, err}
		}()

		cmdOut, cmdErr := run.RunInteractiveCommand(shell, input)

		if cmdErr == nil {
			logger.Log.Debug().Str("cmd", input).Msg("direct exec succeeded")
			return run.NewRunOutput(nil, "", "", cmdOut)
		}

		logger.Log.Debug().Err(cmdErr).Str("cmd", input).Msg("direct exec failed, waiting for engine")
		res := <-engineCh
		if res.err != nil {
			logger.Log.Error().Err(res.err).Msg("exec completion failed too")
			return run.NewRunOutput(fmt.Errorf("command failed: %w. engine failed: %v", cmdErr, res.err), "[error]", "", cmdOut)
		}

		logger.Log.Debug().Str("cmd", res.out.GetCommand()).Bool("executable", res.out.IsExecutable()).Msg("exec completed fallback")
		return parallelFallbackMsg{
			runOut:    run.NewRunOutput(cmdErr, "[error]", "", cmdOut),
			engineOut: *res.out,
		}
	}
}

func (u *Ui) startChat(input string) tea.Cmd {
	logger.Log.Debug().Str("input", input).Msg("starting chat")

	u.state.executing = false
	u.state.confirming = false
	u.state.buffer = ""
	u.state.command = ""

	engine := u.engine

	return func() tea.Msg {
		res, err := engine.ChatCompletion(context.Background(), input)
		if err != nil {
			logger.Log.Error().Err(err).Msg("chat failed")
			return chatErrorMsg{err}
		}

		return parallelFallbackMsg{
			engineOut: *res,
		}
	}
}

func (u *Ui) execCommand(input string) tea.Cmd {
	u.state.querying = false
	u.state.confirming = false
	u.state.executing = true
	u.state.buffer = ""

	shell := u.config.GetSystemConfig().GetShell()
	logger.Log.Info().Str("shell", shell).Str("command", input).Msg("executing command")

	// tea.ExecProcess releases the terminal to the child process, so commands
	// that prompt (sudo) or take the screen (vim, top) work. Their output goes
	// straight to the terminal rather than being captured.
	return tea.ExecProcess(run.PrepareShellCommand(shell, input), func(err error) tea.Msg {
		return commandFinishedMsg{command: input, err: err}
	})
}

func (u *Ui) editSettings() tea.Cmd {
	u.state.querying = false
	u.state.confirming = false
	u.state.executing = true

	systemConfig := u.config.GetSystemConfig()
	c := run.PrepareEditSettingsCommand(
		systemConfig.GetShell(),
		fmt.Sprintf(
			"%s %s",
			systemConfig.GetEditor(),
			systemConfig.GetConfigFile(),
		),
	)

	engineMode := ai.ExecEngineMode
	if u.state.promptMode == ChatPromptMode {
		engineMode = ai.ChatEngineMode
	}
	pipe := u.state.pipe

	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return settingsReloadedMsg{err: err}
		}

		newConfig, err := config.NewConfig()
		if err != nil {
			return settingsReloadedMsg{err: err}
		}

		engine, err := ai.NewEngine(engineMode, newConfig)
		if err != nil {
			return settingsReloadedMsg{err: err}
		}

		if pipe != "" {
			engine.SetPipe(pipe)
		}

		return settingsReloadedMsg{config: newConfig, engine: engine}
	})
}
