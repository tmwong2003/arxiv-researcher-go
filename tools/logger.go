package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
)

// A simple introspection handler that logs the start and end of a chain, as well as any errors that occur while
// running a tool. Implements the [callbacks.Handler] interface.
type LogHandler struct {
	callbacks.SimpleHandler
}

// Singleton [LogHandler] instance used by a chatbot agent and all of its tools.
var Logger = LogHandler{}

func (l LogHandler) HandleChainStart(_ context.Context, inputs map[string]any) {
	var scratchpad_length int
	// The [agents.OneShotZeroAgent] agent uses the key `agent_scratchpad`` to store the scratchpad (i.e., the agent
	// memory).
	if scratchpad, ok := inputs["agent_scratchpad"].(string); ok {
		scratchpad_length = len(scratchpad)
	} else {
		scratchpad_length = 0
	}
	log.Printf(
		"Entering chain with input '%s' and scratchpad of '%d' characters.\n",
		strings.TrimSpace(inputs["input"].(string)),
		scratchpad_length,
	)
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
