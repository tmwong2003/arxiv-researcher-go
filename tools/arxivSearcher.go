package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
)

// ArxivSearcher implements the LangChainGo Tool interface to search a arXiv for papers relevant to a user query.

type ArxivSearcher struct {
	CallbacksHandler callbacks.Handler
}

const (
	arxivSearcherName        = "ArxivSearcher"
	arxivSearcherDescription = `
Search arXiv for relevant papers to a user keyword query. Invoked with a JSON object containing the query "query" and
the number "n" of results to return. Returns a JSON array of dictionary objects containing the title, summary,
authors, and PDF download link for each paper.
`
)

func (tool ArxivSearcher) Name() string {
	return arxivSearcherName
}

func (tool ArxivSearcher) Description() string {
	return arxivSearcherDescription
}

func (tool ArxivSearcher) Call(ctx context.Context, input string) (string, error) {
	log.Printf("Calling tool '%s' with input '%s'.\n", tool.Name(), input)
	if tool.CallbacksHandler != nil {
		tool.CallbacksHandler.HandleToolStart(ctx, input)
	}
	var args struct {
		Query string `json:"query"`
		N     int    `json:"n"`
	}
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		// Failing to unmarshall is _not_ a fatal error. We have observed the agent iterate through different input
		// JSON formats until it discovers the "right" arguments to pass.
		errMessage := fmt.Sprintf("failed while unmarshalling arguments: %s", err)
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		return errMessage, nil
	}
	rawPapers := FetchPapers(args.Query, args.N)
	cookedPapers := make([]map[string]string, len(rawPapers))
	for i, paper := range rawPapers {
		cookedPapers[i] = map[string]string{
			"Title":   paper.Title,
			"Authors": strings.Join(paper.Authors, ", "),
			"PDF URL": paper.PdfUrl,
			"Summary": paper.Summary,
		}
	}
	content, err := json.MarshalIndent(cookedPapers, "", "  ")
	if err != nil {
		// Not great, but maybe the agent will have another tool at its disposal to find results.
		errMessage := fmt.Sprintf("failed while marshalling documents: %s", err)
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		return errMessage, nil
	}
	result := string(content)
	return result, nil
}
