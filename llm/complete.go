package llm

type CompleteRequest struct {
	Model  string
	Prompt string
}

type CompleteResponse struct {
	Text string
}

type CompleteOptions struct {
}

type CompleteOption interface {
	Apply(options *CompleteOptions)
}
