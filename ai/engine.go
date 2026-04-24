package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ilayaraja97/clipper/config"
	"github.com/ilayaraja97/clipper/logger"
	"github.com/ilayaraja97/clipper/system"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const requestTimeout = 60 * time.Second

const noexec = "[noexec]"

type Engine struct {
	mode     EngineMode
	config   *config.Config
	client   llms.Model
	messages []llms.MessageContent
	channel  chan EngineChatStreamOutput
	pipe     string
	running  bool
}

func NewEngine(mode EngineMode, config *config.Config) (*Engine, error) {
	logger.Log.Debug().Str("mode", mode.String()).Msg("creating AI engine")

	opts := []openai.Option{
		openai.WithToken(config.GetAiConfig().GetKey()),
	}
	if config.GetAiConfig().GetBaseURL() != "" {
		logger.Log.Debug().Str("baseURL", config.GetAiConfig().GetBaseURL()).Msg("using custom base URL")
		opts = append(opts, openai.WithBaseURL(config.GetAiConfig().GetBaseURL()))
	}

	if config.GetAiConfig().GetProxy() != "" {
		logger.Log.Debug().Str("proxy", config.GetAiConfig().GetProxy()).Msg("using proxy")
		proxyUrl, err := url.Parse(config.GetAiConfig().GetProxy())
		if err != nil {
			logger.Log.Error().Err(err).Msg("failed to parse proxy URL")
			return nil, err
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}

		httpClient := &http.Client{
			Transport: transport,
		}
		opts = append(opts, openai.WithHTTPClient(httpClient))
	}
	client, err := openai.New(opts...)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to create OpenAI client")
		return nil, err
	}

	logger.Log.Info().Str("mode", mode.String()).Msg("AI engine created")

	return &Engine{
		mode:     mode,
		config:   config,
		client:   client,
		messages: make([]llms.MessageContent, 0),
		channel:  make(chan EngineChatStreamOutput),
		pipe:     "",
		running:  false,
	}, nil
}

func (e *Engine) SetMode(mode EngineMode) *Engine {
	e.mode = mode

	return e
}

func (e *Engine) GetMode() EngineMode {
	return e.mode
}

func (e *Engine) GetChannel() chan EngineChatStreamOutput {
	return e.channel
}

func (e *Engine) SetPipe(pipe string) *Engine {
	e.pipe = pipe

	return e
}

func (e *Engine) Interrupt() *Engine {
	e.channel <- EngineChatStreamOutput{
		content:    "[Interrupt]",
		last:       true,
		interrupt:  true,
		executable: false,
	}

	e.running = false

	return e
}

func (e *Engine) Reset() *Engine {
	e.messages = []llms.MessageContent{}

	return e
}

func (e *Engine) ExecCompletion(input string) (*EngineExecOutput, error) {
	logger.Log.Debug().Str("input", input).Msg("executing completion")
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	e.running = true

	e.appendUserMessage(input)

	resp, err := e.client.GenerateContent(
		ctx,
		e.prepareCompletionMessages(),
		llms.WithModel(e.config.GetAiConfig().GetModel()),
		llms.WithMaxTokens(e.config.GetAiConfig().GetMaxTokens()),
		llms.WithTemperature(e.config.GetAiConfig().GetTemperature()),
	)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			logger.Log.Error().Msg("completion request timed out")
			return nil, fmt.Errorf("request timed out after %v", requestTimeout)
		}
		logger.Log.Error().Err(err).Msg("completion request failed")
		return nil, err
	}
	if len(resp.Choices) == 0 {
		logger.Log.Warn().Msg("empty response from model")
		return nil, fmt.Errorf("empty response from model")
	}

	content := resp.Choices[0].Content
	e.appendAssistantMessage(content)

	var output EngineExecOutput
	err = json.Unmarshal([]byte(content), &output)
	if err != nil {
		logger.Log.Debug().Str("content", content).Msg("JSON unmarshal failed, trying regex extraction")
		re := regexp.MustCompile(`\{.*?\}`)
		match := re.FindString(content)
		if match != "" {
			err = json.Unmarshal([]byte(match), &output)
			if err != nil {
				logger.Log.Error().Err(err).Msg("failed to extract JSON from content")
				return nil, err
			}
		} else {
			logger.Log.Debug().Msg("no JSON found in response, using raw content")
			output = EngineExecOutput{
				Command:     "",
				Explanation: content,
				Executable:  false,
			}
		}
	}

	logger.Log.Debug().
		Bool("executable", output.IsExecutable()).
		Str("cmd", output.GetCommand()).
		Str("exp", output.GetExplanation()).
		Msg("completion result")

	return &output, nil
}

func (e *Engine) ChatStreamCompletion(input string) error {
	logger.Log.Debug().Str("input", input).Msg("starting chat stream completion")
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	e.running = true

	e.appendUserMessage(input)

	var output string
	_, err := e.client.GenerateContent(
		ctx,
		e.prepareCompletionMessages(),
		llms.WithModel(e.config.GetAiConfig().GetModel()),
		llms.WithMaxTokens(e.config.GetAiConfig().GetMaxTokens()),
		llms.WithTemperature(e.config.GetAiConfig().GetTemperature()),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if !e.running {
				logger.Log.Debug().Msg("stream interrupted by user")
				cancel()
				return context.Canceled
			}

			select {
			case <-ctx.Done():
				logger.Log.Debug().Msg("stream timed out")
				return ctx.Err()
			default:
			}

			delta := string(chunk)
			output += delta

			e.channel <- EngineChatStreamOutput{
				content: delta,
				last:    false,
			}

			return nil
		}),
	)

	if err != nil && !errors.Is(err, context.Canceled) {
		if ctx.Err() == context.DeadlineExceeded {
			logger.Log.Error().Msg("chat stream request timed out")
			e.running = false
			return fmt.Errorf("request timed out after %v", requestTimeout)
		}
		logger.Log.Error().Err(err).Msg("chat stream request failed")
		e.running = false
		return err
	}

	executable := false
	if e.mode == ExecEngineMode {
		if !strings.HasPrefix(output, noexec) && !strings.Contains(output, "\n") {
			executable = true
		}
	}

	logger.Log.Debug().
		Str("output", output).
		Bool("executable", executable).
		Msg("chat stream completed")

	e.channel <- EngineChatStreamOutput{
		content:    "",
		last:       true,
		executable: executable,
	}
	e.running = false
	e.appendAssistantMessage(output)

	return nil
}

func (e *Engine) appendUserMessage(content string) *Engine {
	e.messages = append(e.messages, llms.TextParts(llms.ChatMessageTypeHuman, content))

	return e
}

func (e *Engine) appendAssistantMessage(content string) *Engine {
	e.messages = append(e.messages, llms.TextParts(llms.ChatMessageTypeAI, content))

	return e
}

func (e *Engine) AppendAssistantMessage(content string) *Engine {
	return e.appendAssistantMessage(content)
}

func (e *Engine) prepareCompletionMessages() []llms.MessageContent {
	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, e.prepareSystemPrompt()),
	}

	if e.pipe != "" {
		messages = append(
			messages,
			llms.TextParts(llms.ChatMessageTypeHuman, e.preparePipePrompt()),
		)
	}

	messages = append(messages, e.messages...)

	return messages
}

func (e *Engine) preparePipePrompt() string {
	return fmt.Sprintf("I will work on the following input: %s", e.pipe)
}

func (e *Engine) prepareSystemPrompt() string {
	var bodyPart string
	if e.mode == ExecEngineMode {
		bodyPart = e.prepareSystemPromptExecPart()
	} else {
		bodyPart = e.prepareSystemPromptChatPart()
	}

	return fmt.Sprintf("%s\n%s", bodyPart, e.prepareSystemPromptContextPart())
}

func (e *Engine) prepareSystemPromptExecPart() string {
	return "Your are Clipper, a powerful terminal assistant generating a JSON containing a command line for my input.\n" +
		"You will always reply using the following json structure: {\"cmd\":\"the command\", \"exp\": \"some explanation\", \"exec\": true}.\n" +
		"Your answer will always only contain the json structure, never add any advice or supplementary detail or information, even if I asked the same question before.\n" +
		"The field cmd will contain a single line command (don't use new lines, use separators like && and ; instead).\n" +
		"The field exp will contain an short explanation of the command if you managed to generate an executable command, otherwise it will contain the reason of your failure.\n" +
		"The field exec will contain true if you managed to generate an executable command, false otherwise." +
		"\n" +
		"Examples:\n" +
		"Me: list all files in my home dir\n" +
		"Clipper: {\"cmd\":\"ls ~\", \"exp\": \"list all files in your home dir\", \"exec\\: true}\n" +
		"Me: list all pods of all namespaces\n" +
		"Clipper: {\"cmd\":\"kubectl get pods --all-namespaces\", \"exp\": \"list pods form all k8s namespaces\", \"exec\": true}\n" +
		"Me: how are you ?\n" +
		"Clipper: {\"cmd\":\"\", \"exp\": \"I'm good thanks but I cannot generate a command for this. Use the chat mode to discuss.\", \"exec\": false}"
}

func (e *Engine) prepareSystemPromptChatPart() string {
	return "You are Clipper, a powerful terminal assistant created by github.com/ilayaraja97.\n" +
		"You will answer in the most helpful possible way.\n" +
		"Always format your answer in markdown format.\n\n" +
		"For example:\n" +
		"Me: What is 2+2 ?\n" +
		"Clipper: The answer for `2+2` is `4`\n" +
		"Me: +2 again ?\n" +
		"Clipper: The answer is `6`\n"
}

func (e *Engine) prepareSystemPromptContextPart() string {
	part := "My context: "

	if e.config.GetSystemConfig().GetOperatingSystem() != system.UnknownOperatingSystem {
		part += fmt.Sprintf("my operating system is %s, ", e.config.GetSystemConfig().GetOperatingSystem().String())
	}
	if e.config.GetSystemConfig().GetDistribution() != "" {
		part += fmt.Sprintf("my distribution is %s, ", e.config.GetSystemConfig().GetDistribution())
	}
	if e.config.GetSystemConfig().GetHomeDirectory() != "" {
		part += fmt.Sprintf("my home directory is %s, ", e.config.GetSystemConfig().GetHomeDirectory())
	}
	if e.config.GetSystemConfig().GetShell() != "" {
		part += fmt.Sprintf("my shell is %s, ", e.config.GetSystemConfig().GetShell())
	}
	if e.config.GetSystemConfig().GetShell() != "" {
		part += fmt.Sprintf("my editor is %s, ", e.config.GetSystemConfig().GetEditor())
	}
	part += "take this into account. "

	if e.config.GetUserConfig().GetPreferences() != "" {
		part += fmt.Sprintf("Also, %s.", e.config.GetUserConfig().GetPreferences())
	}

	return part
}
