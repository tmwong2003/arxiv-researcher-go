package tools

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

// Represents a paper held by arXiv. Each field corresponds to an equivalent field in the
// [arXiv entry metadata specification]. Note that some fields are optional and may not be present in the metadata
// returned by arXiv.
//
// [arXiv entry metadata specification]: https://info.arxiv.org/help/api/user-manual.html#_entry_metadata
type Paper struct {
	Id               string
	Title            string
	Authors          []string
	Summary          string
	Published        string
	JournalReference string // Optional
	Doi              string // Optional
	PrimaryCategory  string
	Categories       []string
	PdfUrl           string
	ArxivUrl         string
}

// The directory in the local filesystem in which the download tool saves papers. This directory is relative to the
// current working directory of the process running the tool.
const (
	PapersDirectory = "papers"
)

// Get the value of an optional field from the an arXiv metadata.
//
// Returns the value of the field if it exists, otherwise returns an empty string.
func getOptionalField(key string, fields map[string][]ext.Extension) string {
	if field, ok := fields[key]; ok {
		return field[0].Value
	} else {
		return ""
	}
}

// Download a paper from a URL to the local file system. The caller should ensure that the file name is a valid file
// name for the local file system.
//
// Returns nil if the paper is downloaded successfully, otherwise returns an error.
func DownloadPaper(fileName string, url string) error {
	if err := os.MkdirAll(PapersDirectory, 0755); err != nil {
		return fmt.Errorf("failed while downloading from '%s': %w", url, err)
	}
	filePath := filepath.Join(PapersDirectory, fileName)
	if response, err := http.Get(url); err != nil {
		return fmt.Errorf("failed while downloading from '%s': %w", url, err)
	} else {
		defer response.Body.Close()
		if file, err := os.Create(filePath); err != nil {
			return fmt.Errorf("failed while creating file '%s': %w", filePath, err)
		} else {
			defer file.Close()
			if _, err := io.Copy(file, response.Body); err != nil {
				return fmt.Errorf("failed while downloading from '%s': %w", url, err)
			}
		}
	}
	return nil
}

// Query arXiv for papers relevant to a given topic keyword. We parse out the returned data with the help of the
// [arXiv entry metadata specification].
//
// Returns a list of zero or more [Paper] objects corresponding to each relevant paper found.
//
// [arXiv entry metadata specification]: https://info.arxiv.org/help/api/user-manual.html#_entry_metadata
func FetchPapers(keyword string, count int) []Paper {
	arxivParser := gofeed.NewParser()
	keywordEscaped := url.QueryEscape(keyword)
	queryUrl := fmt.Sprintf(
		"http://export.arxiv.org/api/query?search_query=all:%s&start=0&max_results=%d",
		keywordEscaped,
		count,
	)
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
			Title:            strings.ReplaceAll(queryResults.Items[i].Title, "\n", ""),
			Authors:          authors,
			Summary:          queryResults.Items[i].Description,
			Published:        queryResults.Items[i].Published,
			JournalReference: getOptionalField("journal_ref", arxivFields),
			Doi:              getOptionalField("doi", arxivFields),
			PrimaryCategory:  arxivFields["primary_category"][0].Value,
			Categories:       queryResults.Items[i].Categories,
			// Annoyingly arXiv doesn't appear to populate the Links field with the PDF link, but according to the
			// arXiv API specification we can construct the link.
			PdfUrl:   strings.Replace(queryResults.Items[i].Link, "abs", "pdf", 1),
			ArxivUrl: queryResults.Items[i].Link,
		}
		papers = append(papers, paper)
	}
	return papers
}
