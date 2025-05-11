package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
)

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

// ArxivSearcher is a tool that implements the LangChainGo Tool interface to search a arXiv for papers relevant to a
// user query.

func (arxivSearcher ArxivSearcher) Name() string {
	return arxivSearcherName
}

func (arxivSearcher ArxivSearcher) Description() string {
	return arxivSearcherDescription
}

func (arxivSearcher ArxivSearcher) Call(ctx context.Context, input string) (string, error) {
	log.Printf("Calling tool '%s' with input '%s'.\n", arxivSearcher.Name(), input)
	if arxivSearcher.CallbacksHandler != nil {
		arxivSearcher.CallbacksHandler.HandleToolStart(ctx, input)
	}
	var args struct {
		Query string `json:"query"`
		N     int    `json:"n"`
	}
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		// Failing to unmarshall is _not_ a fatal error. We have observed the agent iterate through different input
		// JSON formats until it discovers the "right" arguments to pass.
		return fmt.Sprintf("failed while unmarshalling arguments: %s", err), nil
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
		return fmt.Sprint("failed while marshalling documents: ", err), nil
	}
	result := string(content)
	return result, nil
}
