// Utility functions for querying arXiv and downloading papers.
package tools

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

type Paper struct {
	Id               string
	Title            string
	Summary          string
	Published        string
	JournalReference string
	Doi              string
	PrimaryCategory  string
	Categories       []string
	PdfUrl           string
	ArxivUrl         string
	Authors          []string
}

const (
	PapersDirectory = "papers"
)

func getOptionalField(key string, fields map[string][]ext.Extension) string {
	if field, ok := fields[key]; ok {
		return field[0].Value
	} else {
		return ""
	}
}

func DownloadPaper(fileName string, url string) {
	// Download the paper from the specified URL to the papers directory with the specified file name.
	if err := os.MkdirAll(PapersDirectory, 0755); err != nil {
		log.Fatal(err)
	}
	filePath := filepath.Join(PapersDirectory, fileName)
	if response, err := http.Get(url); err != nil {
		log.Output(2, fmt.Sprintf("Got error downloading paper: %s", err))
	} else {
		defer response.Body.Close()
		if file, err := os.Create(filePath); err != nil {
			log.Output(2, fmt.Sprintf("Got error creating file: %s", err))
		} else {
			defer file.Close()
			if _, err := io.Copy(file, response.Body); err != nil {
				log.Output(2, fmt.Sprintf("Got error downloading paper: %s", err))
			}
		}
	}
}

func FetchPapers(keyword string, count int) []Paper {
	// Query arXiv for papers containing the specified keyword and return a list of Paper objects.
	// We parse out the returned data with the help of the arXiv entry metadata specification:
	// https://info.arxiv.org/help/api/user-manual.html#_entry_metadata
	arxivParser := gofeed.NewParser()
	keywordEscaped := url.QueryEscape(keyword)
	queryUrl := fmt.Sprintf("http://export.arxiv.org/api/query?search_query=all:%s&start=0&max_results=%d", keywordEscaped, count)
	queryResults, _ := arxivParser.ParseURL(queryUrl)
	var papers []Paper
	// The results come back as an RSS feed but with some additional arXiv-specific fields in the extensions. We
	// extract the relevant fields and create a Paper object for each result.
	for i := range queryResults.Items {
		arxivFields := queryResults.Items[i].Extensions["arxiv"]
		var authors []string
		for _, author := range queryResults.Items[i].Authors {
			authors = append(authors, author.Name)
		}
		paper := Paper{
			Id: strings.Replace(queryResults.Items[i].GUID, "http://arxiv.org/abs/", "", 1),
			// Remove the injected newlines from the title
			Title:            strings.Replace(queryResults.Items[i].Title, "\n", "", -1),
			Summary:          queryResults.Items[i].Description,
			Published:        queryResults.Items[i].Published,
			JournalReference: getOptionalField("journal_ref", arxivFields),
			Doi:              getOptionalField("doi", arxivFields),
			PrimaryCategory:  arxivFields["primary_category"][0].Value,
			Categories:       queryResults.Items[i].Categories,
			// Annoyingly arXiv doesn't appear to populate the Links field with the PDF link, but according th the
			// arXiv API specification we can construct the link.
			PdfUrl:   strings.Replace(queryResults.Items[i].Link, "abs", "pdf", 1),
			ArxivUrl: queryResults.Items[i].Link,
			Authors:  authors,
		}
		papers = append(papers, paper)
	}
	return papers
}
