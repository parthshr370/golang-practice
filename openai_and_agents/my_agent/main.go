package main

import (
	"context"
	"fmt"
	"log"
	"my_agent/agent" // Import your new Agent package
	"my_agent/llm"
	"os"
	"time"
)

func main() {
	// basic env setup
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please provide OPENROUTER_API_KEY")
	}

	// Initialize Client this directly calls the llm NewClient to populate w api key
	client := llm.NewClient(apiKey)

	// Initialize Agent (The Brain)
	// We use the Functional Options pattern here unlike before where we were doing  it with structs/json
	// its also prefered we also init the system prompt before but this is simple example to start with
	myAgent := agent.New(client, "google/gemini-3-flash-preview",
		agent.WithSystemPrompts("You are a helpful assistant who speaks like a pirate."),
		agent.WithMaxRetries(3),
	)

	fmt.Println("Starting Agent Carol Sturka...")

	// Create Context , this is essentiall timeout to make sure we dont keep on waiting indefinitely for a dead api
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel() // close the ctx

	// try it out !
	fmt.Println("User: Give me a list of top 10 movies of 2025. overall")
	reply, err := myAgent.Run(ctx, "Give me a list of top 3 movies of 2024.")
	if err != nil {
		log.Fatalf("Turn 1 failed: %v", err)
	}
	fmt.Printf("Agent: %s\n\n", reply)

}

/*
	// --- LEGACY: Manual Client Usage (Phase 1) ---
	// This shows how we used to manually construct requests before we built the Agent abstraction.

	newAgent := llm.NewClient(apiKey)

	req := llm.ChatRequest{
		Model: "openai/gpt-5.2",
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: "Give me list of top 5 movies of 2025 in terms of overall buzz and reviews",
			},
		},
		Temperature: 0.7,
	}

	resp, err := newAgent.CreateChat(ctx, req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	if len(resp.Choices) == 0 {
		log.Fatal("OpenRouter returned zero choices.")
	}
	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
*/
