package tools

import (
	"context"
	"fmt"
	"log"
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
		log.Println("Failed while creating embedder:", err)
		return nil, err
	}
	index.store, err = pinecone.New(
		pinecone.WithAPIKey(os.Getenv(("PINECONE_API_KEY"))),
		pinecone.WithHost(os.Getenv("PINECONE_HOST_NAME")),
		pinecone.WithEmbedder(embedder),
		pinecone.WithNameSpace(os.Getenv("PINECONE_NAME_SPACE")),
	)
	if err != nil {
		log.Println("Failed while connecting to Pinecone:", err)
		return nil, err
	}
	index.context = context.Background()
	return index, nil
}

func (index *Index) AddPapers(papers []Paper) error {
	documents := make([]schema.Document, len(papers))
	for i, paper := range papers {
		content := make([]string, 9)
		content[0] = fmt.Sprintf("Title: %s", paper.Title)
		content[1] = fmt.Sprintf("Authors: %s", strings.Join(paper.Authors, ", "))
		content[2] = fmt.Sprintf("Summary: %s", paper.Summary)
		content[3] = fmt.Sprintf("Published: %s", paper.Published)
		content[4] = fmt.Sprintf("DOI: %s", paper.Doi)
		content[5] = fmt.Sprintf("Primary Category: %s", paper.PrimaryCategory)
		content[6] = fmt.Sprintf("Categories: %s", strings.Join(paper.Categories, ", "))
		content[7] = fmt.Sprintf("PDF URL: %s", paper.PdfUrl)
		content[8] = fmt.Sprintf("arxiv URL: %s", paper.ArxivUrl)
		documents[i] = schema.Document{
			PageContent: strings.Join(content, "\n"),
			Metadata: map[string]any{
				"Title":   paper.Title,
				"Authors": strings.Join(paper.Authors, ", "),
				"DOI":     paper.Doi,
				"PDF URL": paper.PdfUrl,
			},
		}
	}
	_, err := index.store.AddDocuments(index.context, documents)
	return err
}

// Perform a similarity search using the supplied query. Returns up to the top n documents similary to the query.
func (index *Index) SearchIndex(query string, n int) ([]schema.Document, error) {
	var err error
	results, err := index.store.SimilaritySearch(index.context, query, n)
	if err != nil {
		log.Println("Failed while searching:", err)
		return nil, err
	}
	return results, nil
}
