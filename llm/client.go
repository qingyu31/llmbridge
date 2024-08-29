package llm

import "context"

// Client is the client interface for LLM.
type Client interface {
	// Complete completes the prompt with the LLM model.
	Complete(ctx context.Context, req *CompleteRequest, opts ...CompleteOption) (*Result[CompleteResponse], error)
	// CompleteStream completes the prompt with the LLM model and returns a stream result.
	CompleteStream(ctx context.Context, req *CompleteRequest, opts ...CompleteOption) (*StreamResult[CompleteResponse], error)
	// Chat sends a chat request to the LLM model.
	Chat(ctx context.Context, req *ChatRequest, opts ...ChatOption) (*Result[ChatResponse], error)
	// ChatStream sends a chat request to the LLM model and returns a stream result.
	ChatStream(ctx context.Context, req *ChatRequest, opts ...ChatOption) (*StreamResult[ChatResponse], error)
}

type ClientOptions interface {
	InitWithDefault()
}

type ClientOption[T ClientOptions] interface {
	Apply(T)
}

type ClientOptionFunc[T ClientOptions] func(T)

func (f ClientOptionFunc[T]) Apply(t T) {
	f(t)
}

type ClientConstructor[T ClientOptions, C Client] func(opts ...ClientOption[T]) (C, error)

// NewClient creates a new client with the given provider and options.
func NewClient[T ClientOptions, C Client](constructor ClientConstructor[T, C], opts ...ClientOption[T]) (Client, error) {
	return constructor(opts...)
}
