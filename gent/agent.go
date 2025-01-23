package gent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"

	"github.com/noshynosh/eidos/kit"
	"github.com/noshynosh/eidos/kit/research"
	"github.com/noshynosh/eidos/kit/sleep"
	"github.com/noshynosh/eidos/kit/twitter"
	"github.com/noshynosh/eidos/worm"
)

type Agent struct {
	llm  llms.Model // Changed to langchaingo Model interface
	worm worm.MockWorm
}

func NewAgent() (*Agent, error) {
	// Initialize Ollama using langchaingo
	ollamaLLM, err := ollama.New(
		ollama.WithModel("llama3.2"),
		ollama.WithSystemPrompt(seedPrompt),
		ollama.WithFormat("json"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}

	return &Agent{llm: ollamaLLM}, nil
}

// Run starts the agent's execution. It will first seed the LLM with an initial
// prompt and then continue to run the agent until the context is cancelled.
func (a *Agent) Run(ctx context.Context) error {
	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, seedPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, "What would you like to do?"),
	}

	response, err := a.llm.GenerateContent(ctx, messages)
	if err != nil {
		return fmt.Errorf("failed to seed LLM: %w", err)
	}

	if len(response.Choices) == 0 {
		return fmt.Errorf("no choices in response")
	}

	cnt := response.Choices[0].Content

	// Add the response to the messages
	messages = append(messages, llms.TextParts(llms.ChatMessageTypeAI, cnt))

	type llmResp struct {
		Actions   []Action `json:"actions"`
		Reasoning string   `json:"reasoning"`
	}

	var resp llmResp
	if err := json.Unmarshal([]byte(cnt), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal actions: %w", err)
	}

	fmt.Println("first resp", resp)

	bestAction := a.WormSelection(resp.Actions)
	fmt.Println("best action", bestAction.Action)

	var tool kit.Tool
	switch bestAction.Action {
	case actionTweet:
		tool = &twitter.TweetTool{}
	case actionReply:
		tool = &twitter.ReplyTool{}
	case actionResearch:
		tool = &research.Researcher{}
	case actionSleep:
		tool = &sleep.Sleeper{}
	default:
		return fmt.Errorf("unknown action: %s", bestAction.Action)
	}
	tool = &twitter.TweetTool{}

	toolPrompt := kit.BuildToolPrompt(tool)
	fmt.Println("tool prompt", toolPrompt)

	messages = append(messages, llms.TextParts(llms.ChatMessageTypeHuman, toolPrompt))

	newInstructions := fmt.Sprintf(
		"You are a helpful assistant. You are given the following instructions: %s",
		toolPrompt,
	)
	messages = []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, newInstructions),
	}

	response, err = a.llm.GenerateContent(ctx, messages)
	if err != nil {
		return fmt.Errorf("failed to use tool: %w", err)
	}

	fmt.Println("response len", len(response.Choices))
	content := response.Choices[0].Content
	fmt.Println(content)

	// Complete the action
	// if err := tool.Use(ctx, content); err != nil {
	// 	return fmt.Errorf("failed to use tool: %w", err)
	// }

	return nil
}

type Action struct {
	Action     string `json:"action"`
	Chemotaxis int    `json:"chemotaxis"`
	NoseTouch  int    `json:"nose_touch"`
}

func (a *Agent) WormSelection(actions []Action) Action {
	var bestAction Action
	var bestScore int
	for _, action := range actions {
		chemoResp := a.worm.Chemotaxis(action.Chemotaxis)
		noseResp := a.worm.NoseTouch(action.NoseTouch)

		newScore := chemoResp + noseResp
		if newScore > bestScore {
			bestAction = action
			bestScore = newScore
		}
	}

	return bestAction
}
