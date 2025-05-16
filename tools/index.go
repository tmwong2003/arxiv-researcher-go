package tools

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/pinecone"
	"tmwong.org/arxiv-researcher-go/constants"
)

// Represents a connection to a Pinecone index that holds documents for a RAG-based chatbot agent.
type Index struct {
	context context.Context
	store   pinecone.Store
}

// Singleton [Index] connection instance used by a chatbot agent.
var index *Index = nil

// Get the singleton [Index] connection. On the first call, we open a new connection to Pinecone and attach a document
// embedder that computes vector representations of (text) documents for use when indexing documents for storage and
// retrieval. On subsequent calls, we return the existing connection.
//
// Returns the singleton connection if it exists, otherwise returns an error.
func GetIndex() (*Index, error) {
	if index == nil {
		var err error
		index = &Index{}
		embedder, err := embeddings.NewEmbedder(constants.Llm)
		if err != nil {
			return nil, fmt.Errorf("failed while creating embedder: %w", err)
		}
		index.store, err = pinecone.New(
			pinecone.WithAPIKey(os.Getenv(("PINECONE_API_KEY"))),
			pinecone.WithHost(os.Getenv("PINECONE_HOST_NAME")),
			pinecone.WithEmbedder(embedder),
			pinecone.WithNameSpace(os.Getenv("PINECONE_NAME_SPACE")),
		)
		if err != nil {
			return nil, fmt.Errorf("failed while connecting to Pinecone: %w", err)
		}
		index.context = context.Background()
	}
	return index, nil
}

// Add a set of papers to the document index. We treat the concatenated title and summary of each paper as the
// document to index, and the metadata of each paper as the metadata of that document.
//
// Returns nil if we add the papers successfully, otherwise returns an error.
func (index *Index) AddPapers(papers []Paper) error {
	documents := make([]schema.Document, len(papers))
	for i, paper := range papers {
		content := make([]string, 2)
		content[0] = fmt.Sprintf("Title: {%s}", paper.Title)
		content[1] = fmt.Sprintf("Summary: {%s}", paper.Summary)
		documents[i] = schema.Document{
			Metadata: map[string]any{
				"Title":             paper.Title,
				"Authors":           strings.Join(paper.Authors, ", "),
				"Published":         paper.Published,
				"Journal Reference": paper.JournalReference,
				"DOI":               paper.Doi,
				"Primary Category":  paper.PrimaryCategory,
				"Categories":        strings.Join(paper.Categories, ", "),
				"PDF URL":           paper.PdfUrl,
				"arxiv URL":         paper.ArxivUrl,
			},
			PageContent: strings.Join(content, "\n"),
		}
	}
	_, err := index.store.AddDocuments(index.context, documents)
	return err
}
