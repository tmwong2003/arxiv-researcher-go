package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
)

type LogHandler struct {
	callbacks.SimpleHandler
}

func (h LogHandler) HandleToolError(ctx context.Context, err error) {
	log.Printf("failed while running tool: %s\n", err.Error())
}

func makeToolErrorMessage(tool tools.Tool, message string, err error) string {
	return fmt.Sprintf("Tool '%s' %s: %s", tool.Name(), message, err.Error())
}
