package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/callbacks"
)

// IndexSearcher implements the LangChainGo Tool interface to search a document index for papers relevant to a user
// query.

type IndexSearcher struct {
	CallbacksHandler callbacks.Handler
}

const (
	indexSearcherName        = "IndexSearcher"
	indexSearcherDescription = `
Search the document index for relevant papers to a user keyword query. Invoked with a JSON object containing the query
"query" and the number "n" of results to return. Returns a JSON array of dictionary objects containing the title,
summary, authors, and PDF download link for each paper.
`
)

func (tool IndexSearcher) Name() string {
	return indexSearcherName
}

func (tool IndexSearcher) Description() string {
	return indexSearcherDescription
}

func (tool IndexSearcher) Call(ctx context.Context, input string) (string, error) {
	log.Printf("Calling tool '%s' with input '%s'.\n", tool.Name(), input)
	if tool.CallbacksHandler != nil {
		tool.CallbacksHandler.HandleToolStart(ctx, input)
	}
	index, err := GetIndex()
	if err != nil {
		// Failing to get the index is a fatal error, so propagate it to the caller.
		errMessage := fmt.Sprintf("failed while getting index: %s", err)
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		return errMessage, err
	}
	var args struct {
		Query string `json:"query"`
		N     int    `json:"n"`
	}
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		// Failing to unmarshall is _not_ a fatal error. We have observed an agent iterate through different input
		// JSON formats until it discovers the "right" arguments to pass.
		errMessage := fmt.Sprintf("failed while unmarshalling arguments: %s", err)
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		return errMessage, nil
	}
	rawDocuments, err := index.store.SimilaritySearch(index.context, args.Query, args.N)
	if err != nil {
		// Failing to get a result from the index is _not_ a fatal error, because the calling agent have another tool
		// available for searching. We return the error message to the agent so it can decide what to do with it.
		errMessage := fmt.Sprintf("failed while searching index: %s", err)
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		return errMessage, nil
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
		// Not great, but maybe the agent will have another tool at its disposal to find results.
		errMessage := fmt.Sprintf("failed while marshalling documents: %s", err)
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		return errMessage, nil
	}
	result := string(content)
	log.Printf("Tool '%s' returned with '%d' results.\n", tool.Name(), len(rawDocuments))
	return result, nil
}
