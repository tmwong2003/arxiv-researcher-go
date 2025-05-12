package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/callbacks"
)

// PaperDownloader implements the LangChainGo Tool interface to download a paper from a URL to a local file.

type PaperDownloader struct {
	CallbacksHandler callbacks.Handler
}

const (
	paperDownloaderName        = "PaperDownloader"
	paperDownloaderDescription = "Download a paper from a URL."
)

func (tool PaperDownloader) Name() string {
	return paperDownloaderName
}
func (tool PaperDownloader) Description() string {
	return paperDownloaderDescription
}
func (tool PaperDownloader) Call(ctx context.Context, input string) (string, error) {
	log.Printf("Calling tool '%s' with input '%s'.\n", tool.Name(), input)
	if tool.CallbacksHandler != nil {
		tool.CallbacksHandler.HandleToolStart(ctx, input)
	}
	var args struct {
		FileName string `json:"fileName"`
		URL      string `json:"url"`
	}
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		// Failing to unmarshall is _not_ a fatal error. We have observed the agent iterate through different input
		// JSON formats until it discovers the "right" arguments to pass.
		errMessage := fmt.Sprintf("failed while unmarshalling arguments: %s", err)
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		return errMessage, nil
	}
	err := DownloadPaper(args.FileName, args.URL)
	if err != nil {
		// Failing to download the paper is non-fatal. Perhaps the URL is bad.
		errMessage := fmt.Sprintf("failed while downloading paper: %s", err)
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolEnd(ctx, errMessage)
		}
		return errMessage, nil
	}
	return "Paper downloaded successfully.", nil
}
