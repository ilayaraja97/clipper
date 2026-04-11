package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
)

func TestEngineExecOutputGetCommand(t *testing.T) {
	eo := EngineExecOutput{Command: "testCommand"}
	result := eo.GetCommand()

	assert.Equal(t, "testCommand", result)
}

func TestEngineExecOutputGetExplanation(t *testing.T) {
	eo := EngineExecOutput{Explanation: "testExplanation"}
	result := eo.GetExplanation()

	assert.Equal(t, "testExplanation", result)
}

func TestEngineExecOutputIsExecutable(t *testing.T) {
	eo := EngineExecOutput{Executable: true}
	result := eo.IsExecutable()

	assert.True(t, result)
}

func TestEngineChatStreamOutputGetContent(t *testing.T) {
	co := EngineChatStreamOutput{content: "testContent"}
	result := co.GetContent()

	assert.Equal(t, "testContent", result)
}

func TestEngineChatStreamOutputIsLast(t *testing.T) {
	co := EngineChatStreamOutput{last: true}
	result := co.IsLast()

	assert.True(t, result)
}

func TestEngineChatStreamOutputIsInterrupt(t *testing.T) {
	co := EngineChatStreamOutput{interrupt: true}
	result := co.IsInterrupt()

	assert.True(t, result)
}

func TestEngineChatStreamOutputIsExecutable(t *testing.T) {
	co := EngineChatStreamOutput{executable: true}
	result := co.IsExecutable()

	assert.True(t, result)
}

func TestAppendFunctionMessageUsesFunctionRole(t *testing.T) {
	engine := &Engine{
		mode:     ExecEngineMode,
		messages: make([]llms.MessageContent, 0),
	}

	engine.AppendFunctionMessage("Command: ls\nOutput:\nfile.txt")

	if assert.Len(t, engine.messages, 1) {
		assert.Equal(t, llms.ChatMessageTypeFunction, engine.messages[0].Role)
	}
}

func TestReset(t *testing.T) {
	engine := &Engine{
		messages: []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, "test"),
		},
	}

	engine.Reset()

	assert.Empty(t, engine.messages)
}

func TestAppendUserMessage(t *testing.T) {
	engine := &Engine{
		messages: make([]llms.MessageContent, 0),
	}

	engine.appendUserMessage("list files")

	assert.Len(t, engine.messages, 1)
	assert.Equal(t, llms.ChatMessageTypeHuman, engine.messages[0].Role)
}

func TestAppendAssistantMessage(t *testing.T) {
	engine := &Engine{
		messages: make([]llms.MessageContent, 0),
	}

	engine.appendAssistantMessage("Hello!")

	assert.Len(t, engine.messages, 1)
	assert.Equal(t, llms.ChatMessageTypeAI, engine.messages[0].Role)
}

func TestPrepareCompletionMessagesIncludesMessages(t *testing.T) {
	engine := &Engine{
		messages: []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, "user input"),
			llms.TextParts(llms.ChatMessageTypeAI, "assistant response"),
		},
	}

	allMessages := engine.messages

	assert.Len(t, allMessages, 2)
	assert.Equal(t, llms.ChatMessageTypeHuman, allMessages[0].Role)
	assert.Equal(t, llms.ChatMessageTypeAI, allMessages[1].Role)
}

func TestMessageAccumulation(t *testing.T) {
	engine := &Engine{
		messages: make([]llms.MessageContent, 0),
	}

	engine.appendUserMessage("first")
	engine.appendAssistantMessage("response1")
	engine.AppendFunctionMessage("Command: ls")
	engine.appendUserMessage("second")
	engine.appendAssistantMessage("response2")

	assert.Len(t, engine.messages, 5)
	assert.Equal(t, llms.ChatMessageTypeHuman, engine.messages[0].Role)
	assert.Equal(t, llms.ChatMessageTypeAI, engine.messages[1].Role)
	assert.Equal(t, llms.ChatMessageTypeFunction, engine.messages[2].Role)
	assert.Equal(t, llms.ChatMessageTypeHuman, engine.messages[3].Role)
	assert.Equal(t, llms.ChatMessageTypeAI, engine.messages[4].Role)
}
