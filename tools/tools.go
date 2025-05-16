package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/callbacks"
	lcgtools "github.com/tmc/langchaingo/tools"
)

// A generic type to use for implementing tools for chatbot agents. This type implements the [lcgtools.Tool] interface
// and provides a way to define a tool with a name, description, and callback function. We provide this type instead of
// using the raw [lcgtools.Tool] interface to make it simpler to declare type-safe input argument structures for each
// tool, and to consolidate boilerplate introspection calls and argument unmarshalling logic in one place. The input
// argument structure type T should be a structure of the form
//
//	type MyToolArgs struct {
//		Arg1 string `json:"arg1"`
//		Arg2 int    `json:"arg2"`
//	}
//
// Tools should return a string result back to the calling agent that describes the result of their invocations. In
// particular, in the event of an internal error, tools should return a natural language string that describes the
// error and return a nil error, as opposed to returning Go error values, e.g.,
//
//	return fmt.Sprintf("Tool '%s' failed while running tool: %s", tool.Name(), err), nil
//
// In this way, the agent can attempt to recover from the error by rephrasing the question or trying a different tool
// Returning a non-nil error will likely cause the agent to panic and thus prevent it from trying to recover
// gracefully.
type Tool[T any] struct {
	lcgtools.Tool
	// A name for the tool. This name should be unique across all tools made available to LLM-based agents.
	name string
	// A natural language tool description that describes the purpose of the tool LLM-based agents and end users.
	// An agents will parse the description to help them decide when to call the tool. The description should
	// incorporate a pseudocode JSON schema for the input arguments to the tool so that agents will format arguments
	// appropriately. For example, if the tool expects to receive a string argument called "query" and an integer
	// argument called "n", the description should include a JSON schema like
	//
	//	JSON input format: { "fileName": "<file name>", "url": "<paper URL>" }
	description string
	// A callback function that implements the tool logic. The function should take a context and an input argument
	// structure of type T, and return a string result and an error. As discussed above, tools should avoid returning
	// non-nil error values unless the error is not recoverable by the agent.
	Callback func(ctx context.Context, args T) (string, error)
	// An optional set of introspection callback handlers that implement the [callbacks.Handler] interface. The tool
	// will call the HandleToolStart, HandleToolEnd, and HandleToolError methods in the set where appropriate.
	introspectionCallbacks callbacks.Handler
}

// Get the name of a tool.
//
// Implements the [lcgtools.Tool.Name] API call.
func (tool Tool[T]) Name() string {
	return tool.name
}

// Get the description of a tool.
//
// Implements the [lcgtools.Tool.Name] API call.
func (tool Tool[T]) Description() string {
	return tool.description
}

// Unmarshal the raw input from a chatbot agent into the input argument structure for a tool, and call the tool
// callback.
//
// Implements the [lcgtools.Tool.Call] API call.
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
	result, err := tool.Callback(ctx, args)
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
