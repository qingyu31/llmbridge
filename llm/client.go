package llm

// Client is the client interface for LLM.
type Client[O1, O2 any] interface {
	ChatClient[O1]
	CompleteClient[O2]
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

type ClientConstructor[O ClientOptions, C any] func(opts ...ClientOption[O]) (C, error)

// NewClient creates a new client with the given provider and options.
func NewClient[T ClientOptions, C any](constructor ClientConstructor[T, C], opts ...ClientOption[T]) (C, error) {
	return constructor(opts...)
}
