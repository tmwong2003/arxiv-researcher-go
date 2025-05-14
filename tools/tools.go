package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
)

type Tool[T any] struct {
	tools.Tool
	name                   string
	description            string
	callback               func(ctx context.Context, args T) (string, error)
	introspectionCallbacks callbacks.Handler
}

func (tool Tool[T]) Name() string {
	return tool.name
}

func (tool Tool[T]) Description() string {
	return tool.description
}

func (tool Tool[T]) Call(ctx context.Context, input string) (string, error) {
	log.Printf("Calling tool '%s' with input '%s'.\n", tool.Name(), input)
	if tool.introspectionCallbacks != nil {
		tool.introspectionCallbacks.HandleToolStart(ctx, input)
	}
	var args T
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		if tool.introspectionCallbacks != nil {
			tool.introspectionCallbacks.HandleToolError(ctx, err)
		}
		return fmt.Sprintf("Tool '%s' failed while unmarshalling arguments: %s", tool.Name(), err), nil
	}
	log.Printf("Calling tool '%s' callback with args '%+v'.\n", tool.Name(), args)
	result, err := tool.callback(ctx, args)
	if err != nil {
		if tool.introspectionCallbacks != nil {
			tool.introspectionCallbacks.HandleToolError(ctx, err)
		}
		return fmt.Sprintf("Tool '%s' failed while running tool: %s", tool.Name(), err), nil
	}
	if tool.introspectionCallbacks != nil {
		tool.introspectionCallbacks.HandleToolEnd(ctx, result)
	}
	return result, nil
}
