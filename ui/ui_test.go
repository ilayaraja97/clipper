package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewShowsExecutionMessageWhileExecuting(t *testing.T) {
	ui := &Ui{
		state: UiState{
			promptMode: ExecPromptMode,
			executing:  true,
		},
		components: UiComponents{
			renderer: NewRenderer(),
		},
	}

	view := ui.View()

	assert.Contains(t, view, "executing command")
}
