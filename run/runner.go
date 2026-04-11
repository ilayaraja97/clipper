package run

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunInteractiveCommand(shell, input string) (string, error) {
	out, err := prepareRawShellCommand(shell, input).CombinedOutput()
	return string(out), err
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
