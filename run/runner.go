package run

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ilayaraja97/clipper/logger"
)

func RunInteractiveCommand(shell, input string) (string, error) {
	out, err := prepareRawShellCommand(shell, input).CombinedOutput()
	if err != nil {
		logger.Log.Debug().Err(err).Str("shell", shell).Str("command", input).Msg("command failed")
	} else {
		logger.Log.Debug().Str("shell", shell).Str("command", input).Msg("command succeeded")
	}
	return string(out), err
}

// PrepareShellCommand builds the command without running it, so the caller can
// hand it to tea.ExecProcess and give the child process the real terminal.
// Commands that prompt for input (sudo, ssh) or take over the screen (vim, top)
// need a TTY; capturing their output instead makes them hang.
func PrepareShellCommand(shell, input string) *exec.Cmd {
	return prepareRawShellCommand(shell, input)
}

func PrepareEditSettingsCommand(shell, input string) *exec.Cmd {
	command := strings.TrimSpace(strings.TrimRight(input, ";"))

	switch getShellKind(shell) {
	case "powershell":
		return exec.Command(
			shell,
			"-NoProfile",
			"-Command",
			fmt.Sprintf("%s; Write-Host \"\"", command),
		)
	case "cmd":
		return exec.Command(
			shell,
			"/C",
			fmt.Sprintf("%s && echo.", command),
		)
	default:
		return exec.Command(
			shell,
			"-c",
			fmt.Sprintf("%s; echo \"\n\";", command),
		)
	}
}

func prepareRawShellCommand(shell, input string) *exec.Cmd {
	command := strings.TrimSpace(strings.TrimRight(input, ";"))

	switch getShellKind(shell) {
	case "powershell":
		return exec.Command(
			shell,
			"-NoProfile",
			"-Command",
			command,
		)
	case "cmd":
		return exec.Command(
			shell,
			"/C",
			command,
		)
	default:
		return exec.Command(
			shell,
			"-c",
			command,
		)
	}
}

// ttyCommands are commands whose first word always needs a real terminal --
// they either prompt for secret input (sudo, ssh, passwd) or take over the
// screen (vim, top). Racing them through RunInteractiveCommand can't work:
// Bubble Tea still owns the terminal at that point, so the child's prompt
// never reaches the screen and it hangs forever waiting for input that can't
// arrive. Callers must route these straight through tea.ExecProcess instead.
var ttyCommands = map[string]bool{
	"sudo": true, "su": true, "ssh": true, "passwd": true, "visudo": true,
	"crontab": true, "gpg": true,
	"vim": true, "vi": true, "nvim": true, "emacs": true, "nano": true, "pico": true,
	"top": true, "htop": true, "btop": true, "watch": true,
	"less": true, "more": true, "man": true,
	"tmux": true, "screen": true,
	"mysql": true, "psql": true, "sqlite3": true,
	"ftp": true, "sftp": true, "telnet": true,
}

// RequiresTTY reports whether input's first word is a command known to need
// a real terminal rather than captured output.
func RequiresTTY(input string) bool {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return false
	}
	return ttyCommands[filepath.Base(fields[0])]
}

func getShellKind(shell string) string {
	name := strings.ToLower(filepath.Base(strings.TrimSpace(shell)))

	switch name {
	case "powershell", "powershell.exe", "pwsh", "pwsh.exe":
		return "powershell"
	case "cmd", "cmd.exe":
		return "cmd"
	default:
		return "posix"
	}
}
