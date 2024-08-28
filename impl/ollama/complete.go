package ollama

import (
	"context"
	"github.com/ollama/ollama/api"
	"go.qingyu31.com/llmbridge/llm"
	"strings"
)

func (c Client) Complete(ctx context.Context, req *llm.CompleteRequest, opts ...llm.CompleteOption) (*llm.Result[llm.CompleteResponse], error) {
	gr := transformCompleteRequest(req)
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

func (c Client) CompleteStream(ctx context.Context, req *llm.CompleteRequest, opts ...llm.CompleteOption) (*llm.StreamResult[llm.CompleteResponse], error) {
	gr := transformCompleteRequest(req)
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
