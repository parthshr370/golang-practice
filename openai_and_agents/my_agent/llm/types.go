package llm

// omitempty is essentially a way to tell the whole struct that you can just omit mentioning this whole thing when wanting to work with
type ChatRequest struct {
	// Required
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`

	// Optional Configuration
	Temperature      float64         `json:"temperature,omitempty"`
	TopP             float64         `json:"top_p,omitempty"`
	Stream           bool            `json:"stream,omitempty"`
	Stop             []string        `json:"stop,omitempty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	PresencePenalty  float64         `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64         `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int  `json:"logit_bias,omitempty"`
	User             string          `json:"user,omitempty"`
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`
	Seed             int             `json:"seed,omitempty"`

	// Tool Calling Configuration
	// interface{} is essentially way of saying that " Put anything inside of this {} and we will accept it "
	Tools      []Tool      `json:"tools,omitempty"`
	ToolChoice interface{} `json:"tool_choice,omitempty"` // Can be "auto", "none", or a specific tool object
}

// another struct for message passing with its corresponding json
type Message struct {
	Role       string     `json:"role"`    // "user", "assistant", "system", "tool"
	Content    string     `json:"content"` // Keep as string. If null, it will be empty string.
	Name       string     `json:"name,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // Only for Assistant messages
	ToolCallID string     `json:"tool_call_id,omitempty"` // Only for Tool messages
}

// Tool Definitions (For Request) & Tool Calls (For Response)
// essentially defining and establishing the function to perform tool calls

type Tool struct {
	Type     string              `json:"type"` // Usually "function"
	Function FunctionDescription `json:"function"`
}

type FunctionDescription struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Parameters  interface{} `json:"parameters"` // JSON Schema object (map[string]any)
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // Usually "function"
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // This is a JSON STRING, not an object!
}

// Phase 2
// This is where we need to establish what we get back from the API
// this is where we do the rest of the stuff of return like response and choice ??
type ChatResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	FinishReason string      `json:"finish_reason"`      // "stop", "length", "tool_calls", "content_filter"
	Logprobs     interface{} `json:"logprobs,omitempty"` // Null unless top_logprobs set
}

// typical telemetry about what all was consumed in the process
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ResponseFormat struct {
	Type string `json:"type"` // text of json object
}
