package main

import (
	"context"
	"go.qingyu31.com/llmbridge/impl/openai"
	"go.qingyu31.com/llmbridge/llm"
)

func main() {
	client, err := llm.NewClient(openai.NewGPTClient, openai.WithGPTEndpoint("https://api.openai.com"), openai.WithAPIKey("YOUR_API_KEY"))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	req := new(llm.CompleteRequest)
	req.Model = "gpt-4o"
	req.Prompt = "hello"
	result, err := client.Complete(ctx, req)
	if err != nil {
		panic(err)
	}
	println(result.Response.Text)
}
