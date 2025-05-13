package tools

import (
	"context"
	"encoding/json"
	"log"

	"github.com/tmc/langchaingo/callbacks"
)

// PaperDownloader implements the LangChainGo Tool interface to download a paper from a URL to a local file.

type PaperDownloader struct {
	CallbacksHandler callbacks.Handler
}

const (
	paperDownloaderName        = "PaperDownloader"
	paperDownloaderDescription = `
Download a paper from a URL. The caller should ensure that the file name is a valid file name for the local file
system and ends with ".pdf".

JSON input format: { "fileName": "<file name>", "ul": "<paper URL>" }

Sucess: Returns a success message.

Failure: Returns an error message.
`
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
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolError(ctx, err)
		}
		return makeToolErrorMessage(tool, "failed while unmarshalling arguments: %s", err), nil
	}
	err := DownloadPaper(args.FileName, args.URL)
	if err != nil {
		// Failing to download the paper is non-fatal. Perhaps the URL is bad.
		if tool.CallbacksHandler != nil {
			tool.CallbacksHandler.HandleToolError(ctx, err)
		}
		return makeToolErrorMessage(tool, "failed while downloading paper", err), nil
	}
	return "Paper downloaded successfully.", nil
}
