package llm

type Option[T any] interface {
	Apply(t T)
}

type OptionFunc[T any] func(t T)

func (f OptionFunc[T]) Apply(t T) {
	f(t)
}
