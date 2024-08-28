package openai

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"go.qingyu31.com/llmbridge/llm"
	"io"
)

type gptClientOptions struct {
	Endpoint string
	APIKey   string
	Azure    bool
}

func (g gptClientOptions) InitWithDefault() {
}

func NewGPTClient(opts ...llm.ClientOption[*gptClientOptions]) (llm.Client, error) {
	c := new(GPTClient)
	os := new(gptClientOptions)
	for _, o := range opts {
		o.Apply(os)
	}
	var err error
	if os.Azure {
		c.client, err = azopenai.NewClientWithKeyCredential(os.Endpoint, azcore.NewKeyCredential(os.APIKey), nil)
	} else {
		c.client, err = azopenai.NewClientForOpenAI(os.Endpoint, azcore.NewKeyCredential(os.APIKey), nil)
	}
	return c, err
}

func WithGPTEndpoint(endpoint string) llm.ClientOption[*gptClientOptions] {
	return llm.ClientOptionFunc[*gptClientOptions](func(o *gptClientOptions) {
		o.Endpoint = endpoint
	})
}

func WithAPIKey(key string) llm.ClientOption[*gptClientOptions] {
	return llm.ClientOptionFunc[*gptClientOptions](func(o *gptClientOptions) {
		o.APIKey = key
	})
}

func WithAzure() llm.ClientOption[*gptClientOptions] {
	return llm.ClientOptionFunc[*gptClientOptions](func(o *gptClientOptions) {
		o.Azure = true
	})
}

type GPTClient struct {
	client *azopenai.Client
}

func (c GPTClient) Complete(ctx context.Context, req *llm.CompleteRequest, opts ...llm.CompleteOption) (*llm.Result[llm.CompleteResponse], error) {
	var co azopenai.CompletionsOptions
	co.Prompt = []string{req.Prompt}
	co.DeploymentName = &req.Model
	cr, err := c.client.GetCompletions(ctx, co, nil)
	if err != nil {
		return nil, err
	}
	if len(cr.Completions.Choices) == 0 {
		return nil, nil
	}
	result := new(llm.Result[llm.CompleteResponse])
	if cr.Completions.Choices[0].FinishReason != nil {
		result.Reason = string(*cr.Completions.Choices[0].FinishReason)
	}
	if cr.Completions.Choices[0].Text != nil {
		result.Response.Text = *cr.Completions.Choices[0].Text
	}
	return result, nil
}

func (c GPTClient) CompleteStream(ctx context.Context, req *llm.CompleteRequest, opts ...llm.CompleteOption) (*llm.StreamResult[llm.CompleteResponse], error) {
	var co azopenai.CompletionsOptions
	co.Prompt = []string{req.Prompt}
	co.DeploymentName = &req.Model
	cr, er := c.client.GetCompletionsStream(ctx, co, nil)
	if er != nil {
		return nil, er
	}
	result := new(llm.StreamResult[llm.CompleteResponse])
	iter := llm.NewItemIterator[llm.CompleteResponse]()
	result.Iterator = iter
	for {
		resp, err := cr.CompletionsStream.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(resp.Choices) == 0 {
			continue
		}
		if resp.Choices[0].FinishReason != nil {
			result.Reason = string(*resp.Choices[0].FinishReason)
		}
		if resp.Choices[0].Text != nil {
			iter.Write(&llm.CompleteResponse{Text: *resp.Choices[0].Text})
		}
	}
	return result, nil
}

func (c GPTClient) Chat(ctx context.Context, req *llm.ChatRequest, opts ...llm.ChatOption) (*llm.Result[llm.ChatResponse], error) {
	co := c.transformChatRequest(req)
	cr, err := c.client.GetChatCompletions(ctx, *co, nil)
	if err != nil {
		return nil, err
	}
	result := new(llm.Result[llm.ChatResponse])
	result.Response = new(llm.ChatResponse)
	if len(cr.ChatCompletions.Choices) == 0 {
		return result, nil
	}
	if cr.ChatCompletions.Choices[0].FinishReason != nil {
		result.Reason = string(*cr.ChatCompletions.Choices[0].FinishReason)
	}
	if cr.ChatCompletions.Choices[0].Message != nil {
		result.Response.Message = c.transformChatMessage(cr.ChatCompletions.Choices[0].Message)
	}
	return result, nil
}

func (c GPTClient) ChatStream(ctx context.Context, req *llm.ChatRequest, opts ...llm.ChatOption) (*llm.StreamResult[llm.ChatResponse], error) {
	co := c.transformChatRequest(req)
	cs, er := c.client.GetChatCompletionsStream(ctx, *co, nil)
	if er != nil {
		return nil, er
	}
	result := new(llm.StreamResult[llm.ChatResponse])
	iter := llm.NewItemIterator[llm.ChatResponse]()
	result.Iterator = iter
	for {
		resp, err := cs.ChatCompletionsStream.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(resp.Choices) == 0 {
			continue
		}
		if resp.Choices[0].FinishReason != nil {
			result.Reason = string(*resp.Choices[0].FinishReason)
		}
		if resp.Choices[0].Delta != nil {
			cr := new(llm.ChatResponse)
			cr.Message = c.transformChatMessage(resp.Choices[0].Delta)
			iter.Write(cr)
		}
	}
	return result, nil
}

func (c GPTClient) transformChatRequest(req *llm.ChatRequest) *azopenai.ChatCompletionsOptions {
	co := new(azopenai.ChatCompletionsOptions)
	co.DeploymentName = &req.Model
	co.Functions = make([]azopenai.FunctionDefinition, 0, len(req.Functions))
	for _, f := range req.Functions {
		co.Functions = append(co.Functions, *c.transformFunctionDefinition(f))
	}
	if req.FunctionChoice != nil {
		choice := *req.FunctionChoice
		co.FunctionCall = new(azopenai.ChatCompletionsOptionsFunctionCall)
		if choice == llm.FunctionCallChoiceAuto || choice == llm.FunctionCallChoiceNone {
			co.FunctionCall.IsFunction = false
		} else {
			co.FunctionCall.IsFunction = true
			co.FunctionCall.Value = toPtr(choice.String())
		}
	}
	co.Messages = make([]azopenai.ChatRequestMessageClassification, 0, len(req.Messages))
	for _, m := range req.Messages {
		text := ""
		contents := make([]azopenai.ChatCompletionRequestMessageContentPartClassification, 0, len(m.Contents))
		for _, content := range m.Contents {
			switch content.ContentType {
			case llm.ContentTypeText:
				text = string(content.Content)
				part := new(azopenai.ChatCompletionRequestMessageContentPartText)
				part.Text = &text
				contents = append(contents, part)
			case llm.ContentTypeImage:
				part := new(azopenai.ChatCompletionRequestMessageContentPartImage)
				part.ImageURL = &azopenai.ChatCompletionRequestMessageContentPartImageURL{
					URL: toPtr(string(content.Content)),
				}
				contents = append(contents, part)
			}
		}
		switch m.Role {
		case llm.RoleUser:
			um := new(azopenai.ChatRequestUserMessage)
			um.Content = azopenai.NewChatRequestUserMessageContent(contents)
			co.Messages = append(co.Messages, um)
		case llm.RoleSystem:
			sm := new(azopenai.ChatRequestSystemMessage)
			sm.Content = &text
			co.Messages = append(co.Messages, sm)
		case llm.RoleAssistant:
			am := new(azopenai.ChatRequestAssistantMessage)
			am.Content = &text
			co.Messages = append(co.Messages, am)
		}
	}
	return co
}

func (c GPTClient) transformChatMessage(rm *azopenai.ChatResponseMessage) *llm.Message {
	lm := new(llm.Message)
	switch *rm.Role {
	case azopenai.ChatRoleSystem:
		lm.Role = llm.RoleSystem
	case azopenai.ChatRoleUser:
		lm.Role = llm.RoleUser
	case azopenai.ChatRoleAssistant:
		lm.Role = llm.RoleAssistant
	}
	if rm.Content != nil {
		lc := new(llm.MessageContent)
		lc.ContentType = llm.ContentTypeText
		lc.Content = []byte(*rm.Content)
	}
	if rm.FunctionCall != nil {
		lfc := new(llm.FunctionCall)
		lfc.Name = *rm.FunctionCall.Name
		if rm.FunctionCall.Arguments == nil {
			lfc.Parameter = json.RawMessage(*rm.FunctionCall.Arguments)
		}
		lm.FunctionCalls = append(lm.FunctionCalls, lfc)
	}
	return lm
}

func (c GPTClient) transformFunctionDefinition(f *llm.Function) *azopenai.FunctionDefinition {
	fd := new(azopenai.FunctionDefinition)
	fd.Name = toPtr(f.Name)
	fd.Description = toPtr(f.Description)
	fd.Parameters = f.Parameter
	return fd
}
