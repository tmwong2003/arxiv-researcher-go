package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Singleton [Tool] instance to search arXiv for relevant papers to a user keyword query.
var ArxivSearcher = Tool[arxivSearcherArgs]{
	name:                   arxivSearcherName,
	description:            arxivSearcherDescription,
	Callback:               searchArxiv,
	introspectionCallbacks: Logger,
}

const (
	arxivSearcherName        = "ArxivSearcher"
	arxivSearcherDescription = `
Search arXiv for relevant papers to a user keyword query.

JSON input format: { "query": "<user query>", "n": <number of results> }

Success: Returns a JSON array of dictionary objects containing the title, summary, authors, and PDF download link for
each paper

Failure: Returns an error message.
`
)

// The arguments for the [ArxivSearcher] tool. The structure and the [ArxivSearcher] tool description must remain in
// sync with each other to ensure that agents call the tool with the correct JSON argument keys.
type arxivSearcherArgs struct {
	Query string `json:"query"`
	N     int    `json:"n"`
}

// Search arXiv for relevant papers to a user keyword query.
//
// Returns a JSON array of dictionary objects containing the title, summary, authors, and PDF download link for each
// paper if the search is successful, otherwise returns an error message.
func searchArxiv(_ context.Context, args arxivSearcherArgs) (string, error) {
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
		return "", fmt.Errorf("failed while marshalling documents: %w", err)
	}
	result := string(content)
	log.Printf("Tool returned with '%d' results.\n", len(rawPapers))
	return result, nil
}
