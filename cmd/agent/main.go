package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/tools"
	"tmwong.org/arxiv-researcher-go/constants"
	mytools "tmwong.org/arxiv-researcher-go/tools"
)

const (
	queryTemplate = `
Find papers related to '{topic}' in your knowledge database,
and for each paper provide the title, summary, authors, and download link.'
`
)

func run() error {
	agentTools := []tools.Tool{
		mytools.IndexSearcher{},
	}

	agent := agents.NewOneShotAgent(constants.Llm, agentTools, agents.WithMaxIterations(10))
	executor := agents.NewExecutor(agent)

	query := strings.ReplaceAll(queryTemplate, "{topic}", "Quantum field theory")
	fmt.Println("Prompt: ", query)
	answer, err := chains.Run(context.Background(), executor, query)
	fmt.Println(answer)
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal("failed while running chatbot: ", err)
	}
}
