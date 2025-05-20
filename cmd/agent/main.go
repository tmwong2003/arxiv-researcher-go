/*
Search for research papers related to a given topic.
The agent takes a short topic keyword phrase as input and searches for relevant papers in its knowledge database.
If it finds no papers in its database, it expands its search to arXiv.
After completing its search, the agent will display a list of any relevant papers it found
and download the papers to the local file system.

Usage:

	$ go run cmd/agent/main.go <topic keyword>

where <topic keyword> is a short phrase describing the topic of interest.
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	lcgTools "github.com/tmc/langchaingo/tools"
	"tmwong.org/arxiv-researcher-go/constants"
	"tmwong.org/arxiv-researcher-go/tools"
)

// Prompt templates for a zero-shot agent that searches for research papers related to a given topic keyword. The
// LangChainGo [agents.OneShotZeroAgent] prepares a prompt for the LLM using a prefix, a set of format instructions,
// and a suffix.
const (
	// Agents use the prefix template to include common execution context for all of its prompts to the the LLM.
	// The prefix includes natural-language instructions to describe the desired task/behavior of the agent, and a
	// placeholder (.tool_descriptions) for the set of available tools for accessing external data sources.
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
	// The agent uses the format instructions to declare to the LLM how it expects to receive responses from the LLM.
	// Our agent does not use LLM responses directly, so we do not override the default used by
	// [agents.OneShotZeroAgent].
	//  formatInstructions
	// Agents use the suffix template to pass the user query and scratchpad (i.e., the agent memory) to the LLM. Even
	// though the user invokes the agent only once (as opposed to having an interactive conversation with the agent),
	// the agent itself may have multiple iterations in its own internal conversation with the LLM, and thus uses its
	// scratchpad to pass the record of its conversation back and forth with the LLM.
	suffix = `Begin!
Topic keyword: {{.input}}
{{.agent_scratchpad}}`
)

func run() error {
	// Declare the tools that the agent can use to access external data sources.
	agentTools := []lcgTools.Tool{
		tools.ArxivSearcher,
		tools.IndexSearcher,
		tools.PaperDownloader,
	}
	// Create a new one-shot agent that uses our custom prompt templates.
	agent := agents.NewOneShotAgent(
		constants.Llm,
		agentTools,
		// Callbacks for introspection of agent execution, as opposed to callbacks for tool execution.
		agents.WithCallbacksHandler(tools.Logger),
		agents.WithPromptPrefix(prefix),
		agents.WithPromptSuffix(suffix),
	)
	executor := agents.NewExecutor(
		agent,
		agents.WithMaxIterations(25),
	)

	query := ""
	if len(os.Args) > 1 {
		query = strings.Join(os.Args[1:], " ")
	} else {
		query = "one-shot agents"
	}
	fmt.Println("Query: ", query)
	answer, err := chains.Run(context.Background(), executor, query)
	fmt.Println("Answer: ", answer)
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal("Failed while running chatbot: ", err)
	}
}
