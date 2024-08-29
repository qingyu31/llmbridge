package llm

type role string

// ContentType represents the type of content in message.
type ContentType string

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

type ChatRequest struct {
	Model          string
	Messages       []*Message
	Functions      []*Function         // Functions is a list of definitions of functions that can be called in the chat.
	FunctionChoice *FunctionCallChoice // FunctionChoice is defined to specify how to choose the function to call.
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

// MessageContent represents the content in message.
type MessageContent struct {
	ContentType ContentType
	Content     []byte
}

// NewSystemMessage creates a new system message with the given text.
func NewSystemMessage(text string) *Message {
	return &Message{
		Role:     RoleSystem,
		Contents: []*MessageContent{{ContentType: ContentTypeText, Content: []byte(text)}},
	}
}

// NewUserMessage creates a new user message with the given contents.
func NewUserMessage(contents ...*MessageContent) *Message {
	return &Message{Role: RoleUser, Contents: contents}
}

// NewTextContent creates a new MessageContent with the given text.
func NewTextContent(text string) *MessageContent {
	return &MessageContent{ContentType: ContentTypeText, Content: []byte(text)}
}

// NewImageContent creates a new MessageContent with the given image.
func NewImageContent(image []byte) *MessageContent {
	return &MessageContent{ContentType: ContentTypeImage, Content: image}
}
