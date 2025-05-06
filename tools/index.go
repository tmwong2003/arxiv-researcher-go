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

type Index struct {
	context context.Context
	store   pinecone.Store
}

// Create a new Index instance and connect to Pinecone.
func GetIndex() (*Index, error) {
	var index = &Index{}
	var err error
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
	return index, nil
}

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

// Perform a similarity search using the supplied query. Returns up to the top n documents similar to the query.
func (index *Index) SearchIndex(query string, n int) ([]schema.Document, error) {
	var err error
	results, err := index.store.SimilaritySearch(index.context, query, n)
	if err != nil {
		return nil, fmt.Errorf("failed while searching: %w", err)
	}
	return results, nil
}
