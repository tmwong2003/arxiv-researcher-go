package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
)

const (
	downloadPaperToolName = "DownloadPaper"
	searchIndexToolName   = "SearchIndex"
)

// Define the paper download tool for use by an LLM.
var downloadPaperTool = llms.Tool{
	Type: "function",
	Function: &llms.FunctionDefinition{
		Name:        downloadPaperToolName,
		Description: "Download the paper from the specified URL.",
		Parameters: json.RawMessage(`{
			"type": "object",
			"properties": {
				"fileName": {
					"type": "string",
					"description": "The file name to save the paper as."
				},
				"url": {
					"type": "string",
					"description": "The URL of the paper."
				}
			},
			"required": ["fileName", "url"]
		}`),
	},
}

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
	downloadPaperTool,
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
	log.Printf("Executing '%d' tool calls\n", len(resp.Choices[0].ToolCalls))
	for _, toolCall := range resp.Choices[0].ToolCalls {
		log.Printf("Calling tool '%s'\n.", toolCall.FunctionCall.Name)
		var response llms.MessageContent
		switch toolCall.FunctionCall.Name {
		case downloadPaperToolName:
			var args struct {
				FileName string `json:"fileName"`
				URL      string `json:"url"`
			}
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				return messageHistory, fmt.Errorf("failed while unmarshalling arguments: %w", err)
			}
			if err := DownloadPaper(args.FileName, args.URL); err != nil {
				return messageHistory, fmt.Errorf("failed while downloading paper: %w", err)
			}
			content := fmt.Sprintf("OK, I have downloaded the paper '%s' from '%s'.", args.FileName, args.URL)
			response = llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: toolCall.ID,
						Name:       toolCall.FunctionCall.Name,
						Content:    content,
					},
				},
			}
		case searchIndexToolName:
			var args struct {
				Query string `json:"query"`
				N     int    `json:"n"`
			}
			if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
				return messageHistory, fmt.Errorf("failed while unmarshalling arguments: %w", err)
			}
			rawDocuments, err := index.SearchIndex(args.Query, args.N)
			if err != nil {
				return messageHistory, fmt.Errorf("failed while searching index: %w", err)
			}
			cookedDocuments := make([]map[string]string, len(rawDocuments))
			for i, document := range rawDocuments {
				cookedDocuments[i] = map[string]string{
					"Title":   fmt.Sprintf("%s", document.Metadata["Title"]),
					"Authors": fmt.Sprintf("%s", document.Metadata["Authors"]),
					"PDF URL": fmt.Sprintf("%s", document.Metadata["PDF URL"]),
					"Summary": document.PageContent,
				}
			}
			content, err := json.MarshalIndent(cookedDocuments, "", "  ")
			if err != nil {
				return messageHistory, fmt.Errorf("failed while marshalling documents: %w", err)
			}
			response = llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: toolCall.ID,
						Name:       toolCall.FunctionCall.Name,
						Content:    string(content),
					},
				},
			}
		default:
			return messageHistory, fmt.Errorf("got unsupported tool call: '%s'", toolCall.FunctionCall.Name)
		}
		messageHistory = append(messageHistory, response)
	}
	return messageHistory, nil
}
