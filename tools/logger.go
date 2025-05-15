package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
)

type LogHandler struct {
	callbacks.SimpleHandler
}

var Logger = LogHandler{}

func (l LogHandler) HandleChainStart(_ context.Context, inputs map[string]any) {
	log.Println("Entering chain with input:", strings.TrimSpace(inputs["input"].(string)))
}

func (l LogHandler) HandleChainEnd(_ context.Context, outputs map[string]any) {
	log.Println("Exiting chain with outputs:", l.formatAsJson(outputs))
}

func (l LogHandler) HandleToolError(ctx context.Context, err error) {
	log.Printf("Failed while running tool: %s\n", err.Error())
}

func (l LogHandler) formatAsJson(data any) string {
	content, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("failed while marshalling data: %s", err)
	}
	return string(content)
}
