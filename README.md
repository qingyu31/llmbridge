# LLMBridge

LLMBridge is a lightweight SDK designed to simplify and streamline the process of interacting with various large language model (LLM) APIs. By abstracting the complexities of different LLM interfaces, LLMBridge allows developers to switch between models effortlessly and focus on building applications without worrying about the underlying API differences.

## Features
- Unified API: A consistent interface for interacting with multiple LLMs, regardless of their provider.
- Easy Model Switching: Seamlessly switch between different LLM providers with minimal code changes.
- Extensible: Easily add support for new LLMs as they become available.
- Elegant and User-Friendly: Designed with simplicity and developer experience in mind.

## Installation
```bash
go get -u go.qingyu31.com/llmbridge
```

## Usage
Here's a basic example of how to use LLMBridge:
```go
package main

import (
	"context"
	"go.qingyu31.com/llmbridge/llm"
	"go.qingyu31.com/llmbridge/providers/openai"
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

```
Switching to a different LLM provider is as simple as changing the `NewClient` function argument. For example, to use the llama3.1 model with ollama:
```go
client, err := llm.NewClient(ollama.New, ollama.WithHTTPClient(hc), ollama.WithBaseUrl(u))
```

## Supported LLM Providers
- OpenAI
- Azure OpenAI
- Ollama
- [More to come]

## Contributing
We welcome contributions! Please feel free to submit issues or pull requests.

## License
LLMBridge is licensed under the MIT License. See the LICENSE file for more information.