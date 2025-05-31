/*
Populate a private knowledge database with papers from arXiv on a specifc topic of interest.

Usage:

	$ go run cmd/agent/main.go <topic phrase>

where <topic phrase> is a query phrase describing the topic.
The indexer takes the query phrase,
searches arXiv for relevant papers,
and saves metadata for the papers (including abstracts) in its database.
*/
package main

import (
	"log"
	"os"
	"strings"

	"tmwong.org/arxiv-researcher-go/tools"
)

func main() {
	var err error
	index, err := tools.GetIndex()
	if err != nil {
		log.Fatalln("Failed while creating index.")
	}
	var query string
	if len(os.Args) > 1 {
		query = strings.Join(os.Args[1:], " ")
	} else {
		query = "Language Models"
	}
	log.Println("Query: ", query)
	papers := tools.FetchPapers(query, 10)
	if len(papers) == 0 {
		log.Fatalln("Failed while getting papers: Got 0 papers")
	}
	err = index.AddPapers(papers)
	if err != nil {
		log.Fatalln("Failed while adding papers to index:", err)
	}
	log.Printf("Successfully added '%d' papers to index.\n", len(papers))
}
