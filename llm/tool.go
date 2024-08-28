package llm

import "encoding/json"

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

func (t ToolType) String() string {
	return string(t)
}

type Function struct {
	Name        string
	Description string
	Parameter   json.RawMessage
}

func (f *Function) CallChoice() FunctionCallChoice {
	return FunctionCallChoice(f.Name)
}

type FunctionCall struct {
	Name      string
	Parameter json.RawMessage
}

type FunctionCallChoice string

const (
	FunctionCallChoiceNone FunctionCallChoice = "none"
	FunctionCallChoiceAuto FunctionCallChoice = "auto"
)

func (f FunctionCallChoice) String() string {
	return string(f)
}
