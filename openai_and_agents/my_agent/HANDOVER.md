# Handover Protocol: Go Agentic Framework

## 1. Project Context
**Goal:** Build a minimalist, reflection-based Agentic Framework in Go (inspired by Pydantic-AI) from first principles.
**Constraint:** No official SDKs. Pure `net/http`, `encoding/json`, and `reflect`.
**Current State:** Phase 2 (Agent Orchestrator) Complete. Phase 3 (Tooling) In Progress.

## 2. Architecture Overview
*   **Layer 1: Provider ("Mouth")** (`llm/`)
    *   `client.go`: Handles HTTP transport, Auth, Context/Timeouts.
    *   `types.go`: Strict Struct definitions for OpenRouter/OpenAI API.
    *   `messages.go`: Factory functions for type-safe message creation (`NewUserMessage`, etc.).
*   **Layer 2: Agent ("Brain")** (`agent/`)
    *   `agent.go`: Orchestrator holding `History` (Memory) and `Client`.
    *   Pattern: Functional Options (`New(client, model, WithSystemPrompt(...))`).
*   **Layer 3: Tools ("Hands")** (`tools/`)
    *   `registry.go`: Map-based registry for storing functions (`map[string]Tool`).
    *   `jsonschema/schema.go`: Reflection logic to auto-generate JSON Schemas from Go structs.

## 3. Implementation Status

### ✅ Completed
*   **Client:** Fully functional `CreateChat` with context support.
*   **Agent Loop:** `Run` method handles Append -> Send -> Append cycle.
*   **Refactoring:** Migrated from manual struct literals to `llm.New...` helpers.
*   **Registry Foundation:** `Registry` struct and `Register` method are written (but not yet integrated into Agent).
*   **Schema Generation:** Basic reflection logic (`GenerateSchema`) is implemented for Structs/Ints/Strings.

### 🚧 In Progress (The "Hard Part")
*   **Tool Execution:** The `Registry` can *store* tools, but cannot *execute* them yet.
*   **Agent Integration:** The `Agent` struct is not aware of the `Registry`.

## 4. Key Files & Resources
*   **Main Logic:** `/home/parthshr370/Downloads/golang-practice/openai_and_agents/my_agent/agent/agent.go`
*   **Tool Registry:** `/home/parthshr370/Downloads/golang-practice/openai_and_agents/my_agent/tools/registry.go`
*   **Schema Helper:** `/home/parthshr370/Downloads/golang-practice/openai_and_agents/my_agent/tools/jsonschema/schema.go`
*   **Manifesto:** `/home/parthshr370/Downloads/golang-practice/openai_and_agents/my_agent/AGENT.md` (Read this first! Strict "No Spoonfeeding" rules).

## 5. Immediate Next Steps (For Next Agent)
1.  **Implement `Execute`:** Add a method to `Registry` that takes a tool name + JSON args, unmarshals them into a new struct instance (via reflection), and calls the function.
2.  **Integrate Registry:** Add `Registry *tools.Registry` to the `Agent` struct.
3.  **Update Loop:** Modify `agent.Run` to check `resp.Choices[0].FinishReason == "tool_calls"`. If true -> Execute Tool -> Send Result.

**Critical Note:** The user prefers "Mental Model" explanations over code dumps. Explain *how* `reflect.Value.Call` works before writing the implementation.
