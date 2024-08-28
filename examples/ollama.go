package main

import (
	"context"
	"go.qingyu31.com/llmbridge/llm"
	"go.qingyu31.com/llmbridge/providers/ollama"
	"net/http"
	"net/url"
	"time"
)

func main() {
	hc := new(http.Client)
	hc.Timeout = 20 * time.Second
	u, _ := url.Parse("http://localhost:11434")
	client, err := llm.NewClient(ollama.New, ollama.WithHTTPClient(hc), ollama.WithBaseUrl(u))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	req := new(llm.CompleteRequest)
	req.Model = "llama3.1:8b"
	req.Prompt = "hello"
	result, err := client.Complete(ctx, req)
	if err != nil {
		panic(err)
	}
	println(result.Response.Text)
}
