package tools

import (
	"context"
	"fmt"
)

const (
	paperDownloaderName        = "PaperDownloader"
	paperDownloaderDescription = `
Download a paper from a URL. The caller should ensure that the file name is a valid file name for the local file
system and ends with ".pdf".

JSON input format: { "fileName": "<file name>", "url": "<paper URL>" }

Sucess: Returns a success message.

Failure: Returns an error message.
`
)

type downloadPaperArgs struct {
	FileName string `json:"fileName"`
	URL      string `json:"url"`
}

func downloadPaper(_ context.Context, args downloadPaperArgs) (string, error) {
	err := DownloadPaper(args.FileName, args.URL)
	if err != nil {
		// Failing to download the paper is non-fatal. Perhaps the URL is bad.
		return fmt.Sprintf("failed while downloading paper: %s", err), nil
	}
	return fmt.Sprintf("Tool downloaded paper to '%s' successfully.", args.FileName), nil
}

var PaperDownloader = Tool[downloadPaperArgs]{
	name:                   paperDownloaderName,
	description:            paperDownloaderDescription,
	callback:               downloadPaper,
	introspectionCallbacks: Logger,
}
