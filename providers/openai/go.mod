module go.qingyu31.com/llmbridge/providers/openai

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai v0.6.1
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.13.0
	go.qingyu31.com/llmbridge v0.0.0
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	go.qingyu31.com/gtl v0.1.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace go.qingyu31.com/llmbridge v0.0.0 => ../../
