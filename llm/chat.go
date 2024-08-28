package llm

import "encoding/json"

type role string
type ContentType string
type ToolType string

const (
	RoleSystem    role = "system"
	RoleUser      role = "user"
	RoleAssistant role = "assistant"
	RoleFunction  role = "function"
	RoleTool      role = "tool"
)

func (r role) String() string {
	return string(r)
}

const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
	ContentTypeAudio ContentType = "audio"
	ContentTypeVideo ContentType = "video"
)

const (
	ToolTypeFunction ToolType = "function"
)

type ChatRequest struct {
	Model     string
	Messages  []*Message
	Functions []*Function
}

type ChatResponse struct {
	Message *Message
}

type ChatOptions interface {
}

type ChatOption interface {
	Apply(options ChatOptions)
}

type Message struct {
	Role          role
	Contents      []*MessageContent
	FunctionCalls []*FunctionCall
}

type MessageContent struct {
	ContentType
	Content []byte
}

type Function struct {
	Name        string
	Description string
	Parameter   json.RawMessage
}

type FunctionCall struct {
	Name      string
	Parameter json.RawMessage
}
