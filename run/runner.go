package run

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunCommand(cmd string, arg ...string) (string, error) {
	out, err := exec.Command(cmd, arg...).Output()
	if err != nil {
		return fmt.Sprintf("error: %v", err), err
	}

	return string(out), nil
}

func PrepareInteractiveCommand(shell, input string) *exec.Cmd {
	command := strings.TrimSpace(strings.TrimRight(input, ";"))

	switch getShellKind(shell) {
	case "powershell":
		return exec.Command(
			shell,
			"-NoProfile",
			"-Command",
			fmt.Sprintf("Write-Host \"\"; %s; Write-Host \"\"", command),
		)
	case "cmd":
		return exec.Command(
			shell,
			"/C",
			fmt.Sprintf("echo. && %s && echo.", command),
		)
	default:
		return exec.Command(
			shell,
			"-c",
			fmt.Sprintf("echo \"\n\";%s; echo \"\n\";", command),
		)
	}
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
