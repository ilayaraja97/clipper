package ui

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/ilayaraja97/clipper/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigFlow(t *testing.T) {
	t.Run("NewConfigFlow", testNewConfigFlow)
	t.Run("MoveBounds", testConfigFlowMoveBounds)
	t.Run("Next", testConfigFlowNext)
	t.Run("Input", testConfigFlowInput)
	t.Run("DisplayValue", testConfigFlowDisplayValue)
	t.Run("ProgressLines", testConfigFlowProgressLines)
	t.Run("ApplyToPrompt", testConfigFlowApplyToPrompt)
}

func testNewConfigFlow(t *testing.T) {
	flow := newConfigFlow()

	assert.Len(t, flow.fields, 6, "The config flow should include all setup fields.")
	assert.Equal(t, 0, flow.index, "The config flow should start at the first field.")
	assert.Empty(t, flow.values, "The config flow should start with no user values.")
	assert.Equal(t, "key", flow.CurrentField().key, "The first config field should be the API key.")
}

func testConfigFlowMoveBounds(t *testing.T) {
	flow := newConfigFlow()

	flow.Move(-1)
	assert.Equal(t, 0, flow.index, "Moving before the first field should clamp to zero.")

	flow.Move(len(flow.fields) + 10)
	assert.Equal(t, len(flow.fields)-1, flow.index, "Moving past the last field should clamp to the final index.")
}

func testConfigFlowNext(t *testing.T) {
	flow := newConfigFlow()

	for i := 0; i < len(flow.fields)-1; i++ {
		assert.True(t, flow.Next(), "Next should advance while there are remaining fields.")
	}

	assert.Equal(t, len(flow.fields)-1, flow.index, "The flow should stop at the final field.")
	assert.False(t, flow.Next(), "Next should return false once the final field is reached.")
}

func testConfigFlowInput(t *testing.T) {
	flow := newConfigFlow()
	flow.values["key"] = "test-key"
	flow.values["base_url"] = "http://localhost:11434/v1"
	flow.values["model"] = "test-model"
	flow.values["proxy"] = "http://proxy"
	flow.values["temperature"] = "0.3"
	flow.values["max_tokens"] = "4096"

	input := flow.Input()

	assert.Equal(t, "test-key", input.Key, "The API key should be copied into the config input.")
	assert.Equal(t, "http://localhost:11434/v1", input.BaseURL, "The base URL should be copied into the config input.")
	assert.Equal(t, "test-model", input.Model, "The model should be copied into the config input.")
	assert.Equal(t, "http://proxy", input.Proxy, "The proxy should be copied into the config input.")
	assert.Equal(t, "0.3", input.Temperature, "The temperature should be copied into the config input.")
	assert.Equal(t, "4096", input.MaxTokens, "The max tokens should be copied into the config input.")
}

func testConfigFlowDisplayValue(t *testing.T) {
	flow := newConfigFlow()

	assert.Equal(
		t,
		fmt.Sprintf("%s (default)", config.DefaultKey),
		flow.displayValue(flow.fields[0]),
		"The default local API key should remain visible instead of masked.",
	)

	flow.values["key"] = "sk-secret"
	assert.Equal(t, "********", flow.displayValue(flow.fields[0]), "Custom secret values should be masked.")

	assert.Equal(t, "none", flow.displayValue(flow.fields[3]), "An empty optional field without a default should display as none.")

	assert.Equal(
		t,
		fmt.Sprintf("%d (default)", config.DefaultMaxTokens),
		flow.displayValue(flow.fields[5]),
		"The max tokens field should display its default value when empty.",
	)
}

func testConfigFlowProgressLines(t *testing.T) {
	flow := newConfigFlow()
	flow.values["key"] = "sk-secret"
	flow.values["max_tokens"] = "2048"
	flow.Move(5)

	progress := flow.ProgressLines()

	assert.Len(t, progress, len(flow.fields), "The progress output should contain one line per field.")
	assert.Contains(t, progress[0], "API key: `********`", "Progress should mask secret values.")
	assert.Contains(t, progress[5], "> Max tokens: `2048`", "The active field should be prefixed as the current step.")
}

func testConfigFlowApplyToPrompt(t *testing.T) {
	flow := newConfigFlow()
	prompt := NewPrompt(ExecPromptMode)

	flow.ApplyToPrompt(prompt)
	assert.Equal(t, ConfigPromptMode, prompt.GetMode(), "Applying config flow should switch the prompt into config mode.")
	assert.Equal(t, flow.fields[0].placeholder, prompt.input.Placeholder, "The prompt placeholder should match the current field.")
	assert.Equal(t, textinput.EchoPassword, prompt.input.EchoMode, "Secret fields should use password echo mode.")
	assert.True(t, prompt.input.Focused(), "The prompt should be focused after applying config flow.")

	flow.Move(1)
	flow.values["base_url"] = "http://localhost:11434/v1"
	flow.ApplyToPrompt(prompt)

	assert.Equal(t, flow.fields[1].placeholder, prompt.input.Placeholder, "The prompt placeholder should update when the current field changes.")
	assert.Equal(t, "http://localhost:11434/v1", prompt.GetValue(), "The prompt value should match the current field value.")
	assert.Equal(t, textinput.EchoNormal, prompt.input.EchoMode, "Non-secret fields should use normal echo mode.")

	flow.Move(4)
	flow.values["max_tokens"] = strconv.Itoa(config.DefaultMaxTokens)
	flow.ApplyToPrompt(prompt)
	assert.Equal(t, strconv.Itoa(config.DefaultMaxTokens), prompt.GetValue(), "The prompt should show the saved value for later fields too.")
}
