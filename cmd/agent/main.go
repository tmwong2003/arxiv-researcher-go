package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	lcgTools "github.com/tmc/langchaingo/tools"
	"tmwong.org/arxiv-researcher-go/constants"
	"tmwong.org/arxiv-researcher-go/tools"
)

const (
	queryTemplate = `
Find papers related to '{topic}' in your knowledge database. If you find no relevant papers in your database, find
papers in arXiv related to '{topic}'. For each relevant paper you find, provide the title, summary, authors,
and download link, and download the paper. If you find no relevant papers in either the database or arXiv, please say
"No papers found".`
)

func run() error {
	agentTools := []lcgTools.Tool{
		tools.ArxivSearcher,
		tools.IndexSearcher,
		tools.PaperDownloader,
	}

	agent := agents.NewOneShotAgent(
		constants.Llm,
		agentTools,
		agents.WithCallbacksHandler(tools.Logger), // Callbacks for introspection of the agent itself
		agents.WithMaxIterations(10),
	)
	executor := agents.NewExecutor(agent)

	query := strings.ReplaceAll(queryTemplate, "{topic}", "Diffusion Models")
	fmt.Println("Prompt: ", query)
	answer, err := chains.Run(context.Background(), executor, query)
	fmt.Println("Answer: ", answer)
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal("failed while running chatbot: ", err)
	}
}
