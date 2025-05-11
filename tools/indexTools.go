package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/callbacks"
)

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

// IndexSearcher is a tool that implements the LangChainGo Tool interface to search a document index for papers
// relevant to a user query.

func (indexSearcher IndexSearcher) Name() string {
	return indexSearcherName
}

func (indexSearcher IndexSearcher) Description() string {
	return indexSearcherDescription
}

func (indexSearcher IndexSearcher) Call(ctx context.Context, input string) (string, error) {
	log.Printf("Calling tool '%s' with input '%s'.\n", indexSearcher.Name(), input)
	if indexSearcher.CallbacksHandler != nil {
		indexSearcher.CallbacksHandler.HandleToolStart(ctx, input)
	}
	index, err := GetIndex()
	if err != nil {
		// Failing to get the index is a fatal error, so propagate it to the caller.
		errMessage := fmt.Sprintf("failed while getting index: %s", err)
		if indexSearcher.CallbacksHandler != nil {
			indexSearcher.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		return errMessage, err
	}
	resultCount, result, err := index.SearchIndex(input)
	if err != nil {
		// Failing to get a result is _not_ a fatal error, because the calling agent can attempt to call the tool
		// again with a different input format. We have observed the agent iterate through different input formats
		// until it discovers the "right" arguments to pass.
		errMessage := fmt.Sprintf("failed while searching index: %s", err)
		if indexSearcher.CallbacksHandler != nil {
			indexSearcher.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		log.Println(errMessage)
		return errMessage, nil
	}
	if indexSearcher.CallbacksHandler != nil {
		indexSearcher.CallbacksHandler.HandleToolEnd(ctx, result)
	}
	log.Printf("Tool '%s' returned with '%d' results.\n", indexSearcher.Name(), resultCount)
	return result, nil
}
