# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```shell
make build        # go install ./...  (binary lands in $GOBIN)
make test         # go test -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...
make lint         # golangci-lint run   (v2 config in .golangci.yml; gofumpt as formatter)
make fmt          # gofumpt -w .
make check        # test + lint  — what CI runs
make install      # build, then move the binary to ~/.local/bin (or %LOCALAPPDATA%\Programs\clipper on Windows)
```

Single test / package:

```shell
go test ./ui -run TestNewUIInput -v
go test -race ./ai/...
```

`run_benchmarks.sh` is **not** a Go benchmark suite — it is a manual smoke script that shells out to an already-installed `clipper -c` with ~19 prompts against `testdata/`. It hits a real LLM endpoint; do not wire it into `make check`.

## Architecture

Single Bubble Tea program. `main.go` initializes the zerolog file logger, parses input, and runs one `tea.Model`.

**Input → mode selection** (`ui/input.go`, `ui/enum.go`): positional args ⇒ `CliMode` (one-shot, quits after the first result); no args ⇒ `ReplMode`. `-e`/`-c` pick `ExecPromptMode`/`ChatPromptMode`; otherwise `DefaultPromptMode` is resolved from `USER_DEFAULT_PROMPT_MODE` in config. Stdin is drained into `pipe` when it is a named pipe or non-empty, and prepended to the LLM conversation as a synthetic user message.

**`ui/ui.go` is the state machine.** `Ui.state` holds mutually-exclusive-ish flags (`configuring`/`querying`/`confirming`/`executing`) that both `View()` and `handleKeyPress` branch on. All async work returns a `tea.Msg`; the message type determines the transition. When adding behavior, add a message type and a case in `Update` rather than doing work inline — and remember every branch needs a `runMode == CliMode` path that ends in `tea.Quit`, or the one-shot CLI will hang.

**Exec mode races the shell against the LLM** (`Ui.startExec`). The user's raw input is executed directly via `run.RunInteractiveCommand` *while* `Engine.ExecCompletion` runs in a goroutine. If the literal shell command succeeds, the LLM result is dropped and only `run.RunOutput` is emitted. If it fails, both results travel together in `parallelFallbackMsg` and the generated command is offered for `y/N` confirmation. This is the core UX of the tool — typing a real command just runs it, typing English falls through to the model.

**`ai/engine.go`** talks to any OpenAI-compatible endpoint through `langchaingo`'s `llms/openai` (base URL + optional proxy from config). Both exec and chat system prompts instruct the model to reply with `{"cmd":…,"exp":…,"exec":…}`, unmarshalled into `EngineExecOutput`. Parsing degrades in three steps: direct `json.Unmarshal` → regex-extract a JSON object from the prose → fall back to treating the whole response as `Explanation` with `exec:false`. Local/small models frequently need step 2 or 3, so keep that ladder intact. `requestTimeout` is a hard 15s; on error the trailing user message is popped so the conversation isn't poisoned. `e.messages` is mutex-guarded because of the exec-mode goroutine.

Note: despite the names, `ChatStreamCompletion` does not stream — it is a single `GenerateContent` call that returns via `parallelFallbackMsg`. `awaitChatStream`/`EngineChatStreamOutput`/`engine.channel` are leftovers from a streaming implementation; only `Interrupt()` ever writes to the channel. Real streaming would mean feeding `engine.channel` and keeping the `EngineChatStreamOutput` case in `Update`.

**System context** (`system/analyzer.go`) is snapshotted at config load (OS, Linux distro, shell, home dir, username, `$EDITOR`) and appended to every system prompt by `prepareSystemPromptContextPart`, together with free-form `USER_PREFERENCES`. The detected shell is also what `run/runner.go` executes commands with; `getShellKind` dispatches between posix (`-c`), powershell (`-NoProfile -Command`), and cmd (`/C`).

**Config** (`config/`) is `~/.config/clipper.json` via viper. Keys are the uppercase constants in `config/ai.go` and `config/user.go` (`KEY`, `MODEL`, `BASE_URL`, `PROXY`, `TEMPERATURE`, `MAX_TOKENS`, `USER_*`). Defaults target a **local Ollama** server (`http://127.0.0.1:11434/v1`, model `gemma3n:e4b`, key `local`), not a cloud provider. A `viper.ConfigFileNotFoundError` on startup drops the user into the interactive setup flow in `ui/config_form.go` (`configFlow` — an ordered field list navigable with ↑/↓, secrets masked); any other config error prints and quits. `ctrl+s` opens the config file in `$EDITOR` via `tea.ExecProcess` and rebuilds the engine afterwards.

Config uses the **global** viper singleton, so tests must `viper.Reset()` in setup and teardown (see `config/config_test.go`) and cannot run in parallel with each other.

**Logging** (`logger/logger.go`) goes to a per-session file — `~/.local/share/clipper/clipper-<timestamp>.log`, `%LOCALAPPDATA%\clipper\` on Windows — never to stdout, because the TUI owns the terminal. Use `logger.Log`, never `fmt.Print`, for diagnostics.

## Conventions

- Getters are hand-written and fields are unexported across `config`, `system`, `ai`, `history` — follow that when adding state rather than exporting struct fields.
- Tests use testify (`assert`/`require`), live in the same package as the code, and are table-driven where there is more than one case.
- `testdata/` holds synthetic fixtures (JSON, CSV, INI, SQL, logs, Unicode edge cases) used by the manual smoke script, not by the Go tests.
- Releases are goreleaser-driven from tags (`.goreleaser.yaml`); `CHANGELOG.md` is maintained by hand.
