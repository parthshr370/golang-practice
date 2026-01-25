package agent

import (
	"context"
	"my_agent/llm"
)

// the main agent file that sees and takes care of all the things for us
// from message handling to agent initialising to tool calling all is taken care here
// initially we just define the structs of things the agent can take and then we move forward from there
type Agent struct {
	//import from our openrouter client
	client *llm.Client

	// some of our new structs here
	SystemPrompt string
	MaxRetries   int
	Model        string

	// state in the agent something that keeps on passing with each loop
	History []llm.Message
}

type Option func(*Agent)

// we use variadic params here ( the ... do thing , which assigns the var opts to a slice of n values ) these ... tell that this slice can grow , so opts is essentially a slice at its core
func New(client *llm.Client, model string, opts ...Option) *Agent {
	// Default values which is changed eventually if we perform the .Options there and append the value in memory with these pointer ops
	a := &Agent{
		client:     client,
		Model:      model,
		MaxRetries: 1,
		History:    make([]llm.Message, 0),
	}

	// Apply Options
	// range returns ( index , value ) we remove the index w _ since we dont care and just go on w it
	for _, opt := range opts {
		opt(a) // opt a is just a variable holding a function ie a = Agent here
	}

	// Init History with System Prompt if present
	if a.SystemPrompt != "" {
		a.History = append(a.History, llm.Message{Role: "system", Content: a.SystemPrompt})
	}

	return a
}

// now defining our opt functions
// returns kinda like a nested fun where we return the Option variable but inside it only the system prompt that is what we care about when triggering this function
func WithSystemPrompts(prompt string) Option {
	return func(a *Agent) {
		a.SystemPrompt = prompt
	}
}

// same here Option holds the agent but this function will trigger the max retries part only for us which is a lot cleaner
func WithMaxRetries(n int) Option {
	return func(a *Agent) {
		a.MaxRetries = n
	}
}

func (a *Agent) Run(ctx context.Context, usrMsg string) (string, error) {

}
