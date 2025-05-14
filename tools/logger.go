package tools

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/callbacks"
)

type LogHandler struct {
	callbacks.SimpleHandler
}

func (h LogHandler) HandleToolError(ctx context.Context, err error) {
	log.Printf("failed while running tool: %s\n", err.Error())
}

var Logger = LogHandler{}
