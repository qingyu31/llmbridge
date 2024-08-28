package ollama

import (
	"context"
	"encoding/json"
	"github.com/ollama/ollama/api"
	"go.qingyu31.com/llmbridge/llm"
)

func (c Client) Chat(ctx context.Context, req *llm.ChatRequest, opts ...llm.ChatOption) (*llm.Result[llm.ChatResponse], error) {
	cr := transformChatRequest(req)
	cr.Stream = toPtr(false)
	result := new(llm.Result[llm.ChatResponse])
	er := c.client.Chat(ctx, cr, func(response api.ChatResponse) error {
		if response.Done {
			result.Reason = response.DoneReason
		}
		result.Response = buildChatResponse(response)
		return nil
	})
	return result, er

}

func (c Client) ChatStream(ctx context.Context, req *llm.ChatRequest, opts ...llm.ChatOption) (*llm.StreamResult[llm.ChatResponse], error) {
	cr := transformChatRequest(req)
	cr.Stream = toPtr(true)
	result := new(llm.StreamResult[llm.ChatResponse])
	iter := new(iterator[llm.ChatResponse])
	result.Iterator = iter
	er := c.client.Chat(ctx, cr, func(response api.ChatResponse) error {
		if response.Done {
			result.Reason = response.DoneReason
		}
		res := buildChatResponse(response)
		iter.Write(res)
		return nil
	})
	return result, er
}

func transformChatRequest(req *llm.ChatRequest) *api.ChatRequest {
	cr := new(api.ChatRequest)
	cr.Model = req.Model
	cr.Messages = make([]api.Message, 0, len(req.Messages))
	for _, m := range req.Messages {
		var msg api.Message
		msg.Role = m.Role.String()
		for _, c := range m.Contents {
			switch c.ContentType {
			case llm.ContentTypeText:
				msg.Content = string(c.Content)
			case llm.ContentTypeImage:
				msg.Images = append(msg.Images, c.Content)
			}
		}
		cr.Messages = append(cr.Messages, msg)
	}
	return cr
}

func buildChatResponse(response api.ChatResponse) *llm.ChatResponse {
	res := new(llm.ChatResponse)
	res.Message = new(llm.Message)
	switch response.Message.Role {
	case "system":
		res.Message.Role = llm.RoleSystem
	case "user":
		res.Message.Role = llm.RoleUser
	case "assistant":
		res.Message.Role = llm.RoleAssistant
	}
	if response.Message.Content != "" {
		c := new(llm.MessageContent)
		c.ContentType = llm.ContentTypeText
		c.Content = []byte(response.Message.Content)
		res.Message.Contents = append(res.Message.Contents, c)
	}
	for _, img := range response.Message.Images {
		c := new(llm.MessageContent)
		c.ContentType = llm.ContentTypeImage
		c.Content = img
		res.Message.Contents = append(res.Message.Contents, c)
	}
	res.Message.FunctionCalls = make([]*llm.FunctionCall, 0, len(response.Message.ToolCalls))
	for _, fc := range response.Message.ToolCalls {
		call := new(llm.FunctionCall)
		call.Name = fc.Function.Name
		b, _ := json.Marshal(fc.Function.Arguments)
		call.Parameter = b
		res.Message.FunctionCalls = append(res.Message.FunctionCalls, call)
	}
	return res
}
