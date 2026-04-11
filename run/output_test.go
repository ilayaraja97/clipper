package run

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunOutput(t *testing.T) {
	t.Run("HasError", testHasError)
	t.Run("GetErrorMessage", testGetErrorMessage)
	t.Run("GetSuccessMessage", testGetSuccessMessage)
	t.Run("GetContent", testGetContent)
}

func testHasError(t *testing.T) {
	err := errors.New("test error")
	runOutputWithError := NewRunOutput(err, "Error occurred", "Success", "output")
	runOutputWithoutError := NewRunOutput(nil, "Error occurred", "Success", "output")

	assert.True(t, runOutputWithError.HasError(), "RunOutput should have an error.")
	assert.False(t, runOutputWithoutError.HasError(), "RunOutput should not have an error.")
}

func testGetErrorMessage(t *testing.T) {
	err := errors.New("test error")
	runOutput := NewRunOutput(err, "Error occurred", "Success", "output")

	expectedErrorMessage := "Error occurred: test error"
	actualErrorMessage := runOutput.GetErrorMessage()

	assert.Equal(t, expectedErrorMessage, actualErrorMessage, "The error messages should be the same.")
}

func testGetSuccessMessage(t *testing.T) {
	runOutput := NewRunOutput(nil, "Error occurred", "Success", "output")

	expectedSuccessMessage := "Success"
	actualSuccessMessage := runOutput.GetSuccessMessage()

	assert.Equal(t, expectedSuccessMessage, actualSuccessMessage, "The success messages should be the same.")
}

func testGetContent(t *testing.T) {
	runOutput := NewRunOutput(nil, "Error occurred", "Success", "command output")

	assert.Equal(t, "command output", runOutput.GetContent(), "The command output should be preserved.")
}
