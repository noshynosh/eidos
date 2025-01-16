package gent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/noshynosh/eidos/kit"
)

type Agent struct {
	name       string
	objective  string     // what the agent is trying to achieve
	background string     // background information about the agent
	traits     []string   // traits of the agent
	tools      []kit.Tool // tools the agent can use

	// Ollama LLM
	llm *llamaClient
}

func NewAgent(name, objective, background string, traits []string, tools []kit.Tool) *Agent {
	return &Agent{
		name:       name,
		objective:  objective,
		background: background,
		traits:     traits,
		tools:      tools,
		llm:        &llamaClient{client: &http.Client{}},
	}
}

func (a *Agent) Run(ctx context.Context) error {
	// Construct the initial LLM prompt message. This involves the following:
	// 1. The agent's name and objective
	// 2. The agent's background
	// 3. The agent's traits
	// 4. The agent's tools (iterating over the tools and calling their description method)
	toolsPrompt, err := a.constructToolsPrompt()
	if err != nil {
		return err
	}

	prompt := fmt.Sprintf(`
		You are %s. Your objective is to %s.
		
		A bit about you:
		  %s
		
		Your traits are:
		  - %s

		You can utilize tools. A tool is a function that you can call to perform a task. In order to use a tool, you must follow the manual for the tool exactly. If you wish to use a tool, you must provide a response that follows the below format:
		{
			"indicator": "<indicator>",
			"<arg1>": "<value1>",
			"<arg2>": "<value2>",
			...
		}

		For example a tools manual may be:
		{"description", "This tool allows you to sleep for a given amount of time in seconds", "indicator": "sleep","duration": "10"}

		If you wish to use the sleep tool, you must provide a response like:
		{"indicator": "sleep", "duration": "10"}

		IMPORTANT: If you wish to use a tool, you must ONLY respond in this format:
		{
			"indicator": "<indicator>",
			"<arg1>": "<value1>",
			...
		}
		DO NOT write anything else outside this format. Do not explain, introduce, or summarize.

		You have the following tools at your disposal:
		%s

		You will either be given a task to complete or you are free to operate how you see fit given the above information.
	`,
		a.name,
		a.objective,
		a.background,
		strings.Join(a.traits, ", "),
		toolsPrompt,
	)

	fmt.Println(prompt)

	response, err := a.llm.Chat(ctx, prompt)
	if err != nil {
		return err
	}

	fmt.Println(response)

	return nil
}

func (a *Agent) constructToolsPrompt() (string, error) {
	type toolInfo struct {
		Description string `json:"description"`
		Manual      string `json:"manual"`
	}
	tools := []toolInfo{}
	for _, tool := range a.tools {
		tools = append(tools, toolInfo{
			Description: tool.Description(),
			Manual:      tool.Manual(),
		})
	}

	toolsJSON, err := json.Marshal(tools)
	if err != nil {
		return "", err
	}

	return string(toolsJSON), nil
}
