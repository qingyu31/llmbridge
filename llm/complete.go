package llm

import "context"

type CompleteClient[T any] interface {
	// Complete completes the prompt with the LLM model.
	Complete(ctx context.Context, req *CompleteRequest, opts ...T) (*Result[CompleteResponse], error)
	// CompleteStream completes the prompt with the LLM model and returns a stream result.
	CompleteStream(ctx context.Context, req *CompleteRequest, opts ...T) (*StreamResult[CompleteResponse], error)
}

type CompleteRequest struct {
	Model  string
	Prompt string
}

type CompleteResponse struct {
	Text string
}
