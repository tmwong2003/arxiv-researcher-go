package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	lcgTools "github.com/tmc/langchaingo/tools"
	"tmwong.org/arxiv-researcher-go/constants"
	"tmwong.org/arxiv-researcher-go/tools"
)

const (
	prefix = `Today is {{.today}}.
You are a research assistant. You have access to a database of research papers and the arXiv database. When asked
for papers relevant to given topic keyword, you should search for related to the topic in your knowledge database. If
you find no relevant papers in your database, find papers in arXiv related to the topic. For each relevant paper you
find, provide the title, summary, authors, and download link. If you find relevant papers, you should download the
papers to the local file system. If you find no relevant papers in either the database or arXiv, please say "No papers
found".

You have access to the following tools:
{{.tool_descriptions}}
`
	suffix = `Begin!
Topic keyword: {{.input}}
{{.agent_scratchpad}}`
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
		agents.WithPromptPrefix(prefix),
		agents.WithPromptSuffix(suffix),
	)
	executor := agents.NewExecutor(
		agent,
		agents.WithMaxIterations(25),
	)

	query := "Diffusion Models"
	fmt.Println("Query: ", query)
	answer, err := chains.Run(context.Background(), executor, query)
	fmt.Println("Answer: ", answer)
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal("failed while running chatbot: ", err)
	}
}
