package main

import (
	"log"

	"tmwong.org/arxiv-researcher-go/index"
	papers "tmwong.org/arxiv-researcher-go/tools"
)

func main() {
	var err error
	index, err := index.GetIndex()
	if err != nil {
		log.Fatalln("Failed while creating index.")
	}
	papers := papers.FetchPapers("Language Models", 10)
	if len(papers) == 0 {
		log.Fatalln("Failed while getting papers: Got 0 papers")
	}
	err = index.AddPapers(papers)
	if err != nil {
		log.Fatalln("Failed while adding papers to index:", err)
	}
	log.Printf("Successfully added '%d' papers to index.\n", len(papers))
}
