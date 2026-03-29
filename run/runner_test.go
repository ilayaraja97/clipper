package run

import (
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("RunCommand", testRunCommand)
	t.Run("PrepareInteractiveCommand", testPrepareInteractiveCommand)
	t.Run("PrepareEditSettingsCommand", testPrepareEditSettingsCommand)
}

func testRunCommand(t *testing.T) {
	command := "echo"
	args := []string{"Hello, World!"}
	expected := "Hello, World!\n"

	if runtime.GOOS == "windows" {
		command = "cmd"
		args = []string{"/C", "echo Hello, World!"}
		expected = "Hello, World!\r\n"
	}

	output, err := RunCommand(command, args...)
	require.NoError(t, err)

	assert.Equal(t, expected, output, "The command output should be the same.")
}

func testPrepareInteractiveCommand(t *testing.T) {
	cmd := PrepareInteractiveCommand("bash", "echo 'Hello, World!'")

	expectedCmd := exec.Command(
		"bash",
		"-c",
		"echo \"\n\";echo 'Hello, World!'; echo \"\n\";",
	)

	assert.Equal(t, expectedCmd.Args, cmd.Args, "The command arguments should be the same.")
}

func testPrepareEditSettingsCommand(t *testing.T) {
	cmd := PrepareEditSettingsCommand("bash", "nano yo.json")

	expectedCmd := exec.Command(
		"bash",
		"-c",
		"nano yo.json; echo \"\n\";",
	)

	assert.Equal(t, expectedCmd.Args, cmd.Args, "The command arguments should be the same.")
}

func TestPrepareInteractiveCommandPowerShell(t *testing.T) {
	cmd := PrepareInteractiveCommand("powershell", "Get-NetIPConfiguration")

	expectedCmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-Command",
		"Write-Host \"\"; Get-NetIPConfiguration; Write-Host \"\"",
	)

	assert.Equal(t, expectedCmd.Args, cmd.Args, "The PowerShell command arguments should be the same.")
}
