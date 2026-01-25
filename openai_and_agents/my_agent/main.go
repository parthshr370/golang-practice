package main

import (
	"context"
	"fmt"
	"log"
	"my_agent/llm"
	"os"
	"time"
)

func main() {

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please provide the api key ")
	}

	newAgent := llm.NewClient(apiKey)

	req := llm.ChatRequest{
		Model: "openai/gpt-5.2", // use the model of your liking .
		Messages: []llm.Message{ // current way we are managing the response quite primitve , will soon build helper functions
			{
				Role:    "user",
				Content: "Give me list of top 5 movies of 2025 in terms of overall buzz and reviews",
			},
		},
		Temperature: 0.7,
	}
	fmt.Println("Starting our agent Carol Sturka , here have a look ")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Release resources when done

	resp, err := newAgent.CreateChat(ctx, req)

	// if it fails then we return the error with the message
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	// 5. Handle the Response (Safety First)
	if len(resp.Choices) == 0 {
		log.Fatal("OpenRouter returned zero choices. Check your model name or balance.")
	}

	responseReply := resp.Choices[0].Message.Content
	if responseReply == "" {
		log.Fatal("Assistant returned an empty message.")
	}

	fmt.Printf("Response: %s\n", responseReply)
}
