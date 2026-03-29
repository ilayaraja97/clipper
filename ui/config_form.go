package ui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/ilayaraja97/clipper/config"
)

type configField struct {
	key          string
	label        string
	description  string
	placeholder  string
	defaultValue string
	secret       bool
}

type configFlow struct {
	fields []configField
	values map[string]string
	index  int
}

func newConfigFlow() configFlow {
	return configFlow{
		fields: []configField{
			{
				key:          "key",
				label:        "API key",
				description:  "Enter the provider API key. Leave it blank to use the local default.",
				placeholder:  fmt.Sprintf("Default: %s", config.DefaultKey),
				defaultValue: config.DefaultKey,
				secret:       true,
			},
			{
				key:          "base_url",
				label:        "Base URL",
				description:  "Set the OpenAI-compatible endpoint that Clipper should call.",
				placeholder:  fmt.Sprintf("Default: %s", config.DefaultBaseURL),
				defaultValue: config.DefaultBaseURL,
			},
			{
				key:          "model",
				label:        "Model",
				description:  "Choose the default model name sent with each request.",
				placeholder:  fmt.Sprintf("Default: %s", config.DefaultModel),
				defaultValue: config.DefaultModel,
			},
			{
				key:          "proxy",
				label:        "Proxy",
				description:  "Optionally route requests through a proxy server.",
				placeholder:  "Optional, default: none",
				defaultValue: config.DefaultProxy,
			},
			{
				key:          "temperature",
				label:        "Temperature",
				description:  "Set response creativity as a decimal number.",
				placeholder:  fmt.Sprintf("Default: %s", strconv.FormatFloat(config.DefaultTemperature, 'f', -1, 64)),
				defaultValue: strconv.FormatFloat(config.DefaultTemperature, 'f', -1, 64),
			},
			{
				key:          "max_tokens",
				label:        "Max tokens",
				description:  "Choose the maximum number of tokens to generate per response.",
				placeholder:  fmt.Sprintf("Default: %d", config.DefaultMaxTokens),
				defaultValue: strconv.Itoa(config.DefaultMaxTokens),
			},
		},
		values: map[string]string{},
	}
}

func (f *configFlow) CurrentField() configField {
	return f.fields[f.index]
}

func (f *configFlow) CurrentValue() string {
	return f.values[f.CurrentField().key]
}

func (f *configFlow) SetCurrentValue(value string) {
	f.values[f.CurrentField().key] = value
}

func (f *configFlow) Next() bool {
	if f.index >= len(f.fields)-1 {
		return false
	}

	f.index++

	return true
}

func (f *configFlow) Move(delta int) {
	next := f.index + delta
	if next < 0 {
		next = 0
	}
	if next >= len(f.fields) {
		next = len(f.fields) - 1
	}

	f.index = next
}

func (f *configFlow) Input() config.ConfigInput {
	return config.ConfigInput{
		Key:         f.values["key"],
		Model:       f.values["model"],
		BaseURL:     f.values["base_url"],
		Proxy:       f.values["proxy"],
		Temperature: f.values["temperature"],
		MaxTokens:   f.values["max_tokens"],
	}
}

func (f *configFlow) ProgressLines() []string {
	progress := make([]string, 0, len(f.fields))
	for index, field := range f.fields {
		prefix := "-"
		if index == f.index {
			prefix = ">"
		}

		progress = append(progress, fmt.Sprintf("%s %s: `%s`", prefix, field.label, f.displayValue(field)))
	}

	return progress
}

func (f *configFlow) displayValue(field configField) string {
	value := f.values[field.key]
	if value == "" {
		if field.defaultValue == "" {
			return "none"
		}

		return fmt.Sprintf("%s (default)", f.maskValue(field, field.defaultValue))
	}

	return f.maskValue(field, value)
}

func (f *configFlow) maskValue(field configField, value string) string {
	if field.secret && value != "" && value != config.DefaultKey {
		return "********"
	}

	return value
}

func (f *configFlow) ApplyToPrompt(prompt *Prompt) {
	field := f.CurrentField()
	prompt.SetMode(ConfigPromptMode)
	prompt.SetPlaceholder(field.placeholder)
	prompt.SetValue(f.CurrentValue())
	if field.secret {
		prompt.SetEchoMode(textinput.EchoPassword)
	} else {
		prompt.SetEchoMode(textinput.EchoNormal)
	}
	prompt.Focus()
}
