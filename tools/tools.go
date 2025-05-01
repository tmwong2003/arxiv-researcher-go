package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

const (
	searchIndexToolName = "SearchIndex"
)

// Define the search index tool for use by an LLM.
var searchIndexTool = llms.Tool{
	Type: "function",
	Function: &llms.FunctionDefinition{
		Name:        searchIndexToolName,
		Description: "Search the index for relevant papers about language models.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"query": {
					"type": "string",
					"description": "The query."
				},
				"n": {
					"type": "integer",
					"description": "The number of results to return."
				}
			},
			"required": ["query", "n"]
		}`),
	},
}

// Define the tools avaialble for use by an LLM.
var Tools = []llms.Tool{
	searchIndexTool,
}

// Dispatch tool calls from an LLM to the appropriate handler function. Returns the supplied message history with the
// responses from successful tool calls appended.
func ExecuteTools(ctx context.Context, llm llms.Model, messageHistory []llms.MessageContent, resp *llms.ContentResponse) ([]llms.MessageContent, error) {
	var err error
	index, err := GetIndex()
	if err != nil {
		return messageHistory, err
	}
	log.Println("Executing", len(resp.Choices[0].ToolCalls), "tool calls")
	for _, toolCall := range resp.Choices[0].ToolCalls {
		log.Printf("Calling tool '%s'\n.", toolCall.FunctionCall.Name)
		switch toolCall.FunctionCall.Name {
		case searchIndexToolName:
			var args struct {
				Query string `json:"query"`
				N     int    `json:"n"`
			}
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				return messageHistory, fmt.Errorf("failed while unmarshalling arguments: %w", err)
			}

			documents, err := index.SearchIndex(args.Query, args.N)
			if err != nil {
				return messageHistory, fmt.Errorf("failed while searching index: %w", err)
			}

			titles := make([]string, len(documents))
			for i, document := range documents {
				titles[i] = fmt.Sprintf("'%s'", document.Metadata["Title"])
			}

			searchIndexResponse := llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: toolCall.ID,
						Name:       toolCall.FunctionCall.Name,
						Content:    strings.Join(titles, ", "),
					},
				},
			}
			messageHistory = append(messageHistory, searchIndexResponse)
		default:
			return messageHistory, fmt.Errorf("got unsupported tool call: '%s'", toolCall.FunctionCall.Name)
		}
	}

	return messageHistory, nil
}
