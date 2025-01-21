package gent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/noshynosh/eidos/kit"
	"github.com/noshynosh/eidos/kit/twitter"
	"github.com/noshynosh/eidos/llm"
	"github.com/noshynosh/eidos/worm"
)

type Agent struct {
	llm  *llm.OllamaClient // LLM client
	worm worm.MockWorm     // Worm client
}

func NewAgent() *Agent {
	return &Agent{
		llm: &llm.OllamaClient{Client: &http.Client{}},
	}
}

// Run starts the agent's execution. It will first seed the LLM with an initial
// prompt and then continue to run the agent until the context is cancelled.
func (a *Agent) Run(ctx context.Context) error {
	messages := make([]llm.ChatMessage, 0)
	seedResp, err := a.seedLLM(ctx)
	if err != nil {
		return err
	}
	messages = append(messages, seedResp)

	actions := make([]Action, 0)
	if err := json.Unmarshal([]byte(seedResp.Content), &actions); err != nil {
		return fmt.Errorf("failed to unmarshal actions: %w", err)
	}

	fmt.Println(actions)

	bestAction := a.WormSelection(actions)
	fmt.Println(bestAction.Action)

	// temp
	bestAction.Action = actionTweet

	var tool kit.Tool
	switch bestAction.Action {
	case actionTweet:
		tool = &twitter.TweetTool{}
	// case actionReply:
	// 	replyTool := twitter.ReplyTool{}
	// 	replyTool.Use(ctx, bestAction.Action)
	// case actionResearch:
	// 	researchTool := twitter.ResearchTool{}
	// 	researchTool.Use(ctx, bestAction.Action)
	default:
		return fmt.Errorf("unknown action: %s", bestAction.Action)
	}

	toolPrompt := kit.BuildToolPrompt(tool)
	fmt.Println(toolPrompt)

	response, err := a.llm.Chat(ctx, toolPrompt, messages...)
	if err != nil {
		return fmt.Errorf("failed to use tool: %w", err)
	}

	fmt.Println(response.Content)

	// Complete the action
	if err := tool.Use(ctx, response.Content); err != nil {
		return fmt.Errorf("failed to use tool: %w", err)
	}

	return nil
}

func (a *Agent) seedLLM(ctx context.Context) (llm.ChatMessage, error) {
	response, err := a.llm.Chat(ctx, seedPrompt)
	if err != nil {
		return llm.ChatMessage{}, fmt.Errorf("failed to seed LLM: %w", err)
	}

	return response, nil
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
