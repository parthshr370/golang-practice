package tools

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 1. Define a dummy struct for arguments
type WeatherArgs struct {
	City string `json:"city" description:"The city name to check weather for"`
	Days int    `json:"days" description:"Number of days for forecast"`
}

// 2. Define a dummy function
func GetWeather(args WeatherArgs) string {
	return fmt.Sprintf("Weather in %s for %d days is sunny", args.City, args.Days)
}

func TestRegistry_Register(t *testing.T) {
	// 3. Initialize Registry
	registry := NewRegistry()

	// 4. Register the Tool
	err := registry.Register("get_weather", "Get current weather", GetWeather)
	if err != nil {
		t.Fatalf("Failed to register tool: %v", err)
	}

	// 5. Verify it was stored
	tool, exists := registry.tools["get_weather"]
	if !exists {
		t.Fatal("Tool 'get_weather' was not found in registry")
	}

	// 6. Print the generated schema to verify it looks like JSON Schema
	schemaJSON, _ := json.MarshalIndent(tool.Schema, "", "  ")
	fmt.Printf("Generated Schema:\n%s\n", string(schemaJSON))

	// Basic assertion on schema structure
	props, ok := tool.Schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("Schema missing 'properties' field")
	}

	if _, ok := props["city"]; !ok {
		t.Error("Schema missing 'city' property")
	}
}
