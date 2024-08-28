package llm

import "context"

type Client interface {
	Complete(ctx context.Context, req *CompleteRequest, opts ...CompleteOption) (*Result[CompleteResponse], error)
	CompleteStream(ctx context.Context, req *CompleteRequest, opts ...CompleteOption) (*StreamResult[CompleteResponse], error)
	Chat(ctx context.Context, req *ChatRequest, opts ...ChatOption) (*Result[ChatResponse], error)
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

func NewClient[T ClientOptions, C Client](constructor ClientConstructor[T, C], opts ...ClientOption[T]) (Client, error) {
	return constructor(opts...)
}
