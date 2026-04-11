package run

import (
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("RunInteractiveCommand", testRunInteractiveCommand)
	t.Run("PrepareEditSettingsCommand", testPrepareEditSettingsCommand)
}

func testRunInteractiveCommand(t *testing.T) {
	command := "printf 'Hello, World!'"
	expected := "Hello, World!"
	if runtime.GOOS == "windows" {
		command = "echo Hello, World!"
		expected = "Hello, World!\r\n"
	}

	output, err := RunInteractiveCommand(defaultShellForTests(), command)
	require.NoError(t, err)

	assert.Equal(t, expected, output, "The interactive command output should be captured.")
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

func defaultShellForTests() string {
	if runtime.GOOS == "windows" {
		return "cmd"
	}

	return "bash"
}
