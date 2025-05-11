package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"tmwong.org/arxiv-researcher-go/constants"
	"tmwong.org/arxiv-researcher-go/tools"
)

const (
	queryTemplate = `
Find papers related to '{topic}' in your knowledge database,
and for each paper provide the title, summary, authors, and download link.'
`
)

func appendLlmResponse(messageHistory []llms.MessageContent, response *llms.ContentResponse) []llms.MessageContent {
	toolCalls := make([]string, len(response.Choices[0].ToolCalls))
	for i, choice := range response.Choices {
		for j, toolCall := range choice.ToolCalls {
			toolCalls[j] = fmt.Sprintf("%s(%s)", toolCall.FunctionCall.Name, toolCall.FunctionCall.Arguments)
		}
		log.Printf("Choice '%d': content '%s', tool calls '%s'.\n", i, choice.Content, strings.Join(toolCalls, ", "))
	}
	// Select the first choice and append the tool calls to the message history.
	choice := response.Choices[0]
	responseText := llms.TextParts(llms.ChatMessageTypeAI, choice.Content)
	for _, toolCall := range choice.ToolCalls {
		responseText.Parts = append(responseText.Parts, toolCall)
	}
	return append(messageHistory, responseText)
}

func executeToolCalls(ctx context.Context, messageHistory []llms.MessageContent, response *llms.ContentResponse) []llms.MessageContent {
	messageHistory = appendLlmResponse(messageHistory, response)
	messageHistory, responseCount, err := tools.ExecuteTools(ctx, constants.Llm, messageHistory, response)
	if err != nil {
		log.Fatal("failed while executing tool calls for the LLM: ", err)
	}
	for i := responseCount; i > 0; i-- {
		fmt.Println(messageHistory[len(messageHistory)-i].Parts[0])
	}
	return messageHistory
}

func main() {
	ctx := context.Background()

	query := strings.ReplaceAll(queryTemplate, "{topic}", "Language Models")
	fmt.Println("Prompt: ", query)
	messageHistory := []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeHuman, query)}
	response, err := constants.Llm.GenerateContent(ctx, messageHistory, llms.WithTools(tools.Tools))
	if err != nil {
		log.Fatal("failed while invoking the LLM: ", err)
	}
	messageHistory = executeToolCalls(ctx, messageHistory, response)

	command := "Download any papers that you find interesting."
	messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeHuman, command))
	response, err = constants.Llm.GenerateContent(ctx, messageHistory, llms.WithTools(tools.Tools))
	if err != nil {
		log.Fatal("failed while invoking the LLM: ", err)
	}
	executeToolCalls(ctx, messageHistory, response)
}
