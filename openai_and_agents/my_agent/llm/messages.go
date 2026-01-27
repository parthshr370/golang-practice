package llm

import "fmt"

func NewSystemMessage(content string) Message {

	return Message{
		Role:    "system",
		Content: content,
	}
}

func NewUserMessage(content string) Message {

	return Message{
		Role:    "user",
		Content: content,
	}
}

// Standard reply
func NewAssistantMessage(content string) Message {
	return Message{
		Role:    "assistant",
		Content: content,
	}
}

// When the LLM decides to call tools (Pydantic's `LLMToolCalls`)
// We rarely create this manually (we receive it), but useful for testing/mocks
func NewToolCallMessage(calls []ToolCall) Message {
	return Message{
		Role:      "assistant",
		ToolCalls: calls,
		// Content must be empty for tool calls in strict OpenAI standards
	}
}

// Crucial: Forces you to pass the ID so you don't break the chain
func NewToolResult(toolCallID string, output string) Message {
	return Message{
		Role:       "tool",
		ToolCallID: toolCallID,
		Content:    output,
	}
}

func NewToolError(toolCallID string, err error) Message {
	return Message{
		Role:       "tool",
		ToolCallID: toolCallID,
		Content:    fmt.Sprintf("Error executing tool: %v. Please fix your arguments.", err),
	}
}
