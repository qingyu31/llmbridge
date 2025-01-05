package ollama

import (
	"context"
	"github.com/ollama/ollama/api"
	"go.qingyu31.com/llmbridge/llm"
	"strings"
)

type CompleteOption llm.Option[*api.GenerateRequest]

func (c Client) Complete(ctx context.Context, req *llm.CompleteRequest, opts ...CompleteOption) (*llm.Result[llm.CompleteResponse], error) {
	gr := transformCompleteRequest(req)
	for _, opt := range opts {
		opt.Apply(gr)
	}
	result := new(llm.Result[llm.CompleteResponse])
	sb := new(strings.Builder)
	er := c.client.Generate(ctx, gr, func(response api.GenerateResponse) error {
		if response.Done {
			result.Reason = response.DoneReason
		}
		sb.WriteString(response.Response)
		return nil
	})
	result.Response = new(llm.CompleteResponse)
	result.Response.Text = sb.String()
	return result, er
}

func (c Client) CompleteStream(ctx context.Context, req *llm.CompleteRequest, opts ...CompleteOption) (*llm.StreamResult[llm.CompleteResponse], error) {
	gr := transformCompleteRequest(req)
	for _, opt := range opts {
		opt.Apply(gr)
	}
	result := new(llm.StreamResult[llm.CompleteResponse])
	iter := llm.NewItemIterator[llm.CompleteResponse]()
	result.Iterator = iter
	er := c.client.Generate(ctx, gr, func(response api.GenerateResponse) error {
		res := new(llm.CompleteResponse)
		if response.Done {
			result.Reason = response.DoneReason
		}
		res.Text = response.Response
		iter.Write(res)
		return nil
	})
	return result, er
}

func transformCompleteRequest(req *llm.CompleteRequest) *api.GenerateRequest {
	gr := new(api.GenerateRequest)
	gr.Prompt = req.Prompt
	gr.Model = req.Model
	return gr
}

func WithCompleteOptions(key string, value any) CompleteOption {
	return llm.OptionFunc[*api.GenerateRequest](func(req *api.GenerateRequest) {
		if req.Options == nil {
			req.Options = make(map[string]any)
		}
		req.Options[key] = value
	})
}

func WithCompleteFormat(format string) CompleteOption {
	return llm.OptionFunc[*api.GenerateRequest](func(req *api.GenerateRequest) {
		req.Format = format
	})
}
