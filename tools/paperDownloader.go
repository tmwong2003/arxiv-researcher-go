package tools

import (
	"context"
	"fmt"
)

// Singleton [Tool] instance of a tool to download a paper from a URL to the local file system.
var PaperDownloader = Tool[downloadPaperArgs]{
	name:                   paperDownloaderName,
	description:            paperDownloaderDescription,
	Callback:               downloadPaper,
	introspectionCallbacks: Logger,
}

const (
	paperDownloaderName        = "PaperDownloader"
	paperDownloaderDescription = `
Download a paper from a URL to the local file system. The caller should ensure that the file name is a valid file name
for the local file system and ends with ".pdf".

JSON input format: { "fileName": "<file name>", "url": "<paper URL>" }

Sucess: Returns a success message.

Failure: Returns an error message.
`
)

// The arguments for the [PaperDownloader] tool. The structure and the [PaperDownloader] tool description must remain
// in sync with each other to ensure that agents call the tool with the correct JSON argument keys.
type downloadPaperArgs struct {
	FileName string `json:"fileName"`
	URL      string `json:"url"`
}

// Download a paper from a URL to the local file system. The caller should ensure that the file name is a valid file
// name for the local file system and ends with ".pdf".
//
// Returns a success message if the paper is downloaded successfully, otherwise returns an error message.
func downloadPaper(_ context.Context, args downloadPaperArgs) (string, error) {
	err := DownloadPaper(args.FileName, args.URL)
	if err != nil {
		return fmt.Sprintf("failed while downloading paper: %s", err), nil
	}
	return fmt.Sprintf("Tool downloaded paper to '%s' successfully.", args.FileName), nil
}
