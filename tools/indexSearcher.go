package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// Singleton [Tool] instance to search the document index for relevant papers to a user keyword query.
var IndexSearcher = Tool[indexSearcherArgs]{
	name:                   indesSearcherName,
	description:            indexSearcherDescription,
	Callback:               searchIndex,
	introspectionCallbacks: Logger,
}

const (
	indesSearcherName        = "IndexSearcher"
	indexSearcherDescription = `
Search the document index for relevant papers to a user keyword query.

JSON input format: { "query": "<user query>", "n": <number of results> }

Success: Returns a JSON array of dictionary objects containing the title, summary, authors, and PDF download link for
each paper

Failure: Returns an error message.
`
)

// The arguments for the [IndexSearcher] tool. The structure and the [IndexSearcher] tool description must remain in
// sync with each other to ensure that agents call the tool with the correct JSON argument keys.
type indexSearcherArgs struct {
	Query string `json:"query"`
	N     int    `json:"n"`
}

// Search a document index for relevant papers to a user keyword query.
//
// Returns a JSON array of dictionary objects containing the title, summary, authors, and PDF download link for each
// paper if the search is successful, otherwise returns an error message.
func searchIndex(_ context.Context, args indexSearcherArgs) (string, error) {
	index, err := GetIndex()
	if err != nil {
		return "", fmt.Errorf("failed while getting index: %s", err)
	}
	rawDocuments, err := index.store.SimilaritySearch(index.context, args.Query, args.N)
	if err != nil {
		return fmt.Sprintf("failed while searching index: %s", err), nil
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
		return "", fmt.Errorf("failed while marshalling documents: %w", err)
	}
	result := string(content)
	log.Printf("Tool returned with '%d' results.\n", len(rawDocuments))
	return result, nil
}
