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
	Parameter   json.RawMessage // Parameter is the JSON representation of the function's parameter. It may be JSONSCHEMA in most cases.
}

// CallChoice returns the choice to call this function.
func (f *Function) CallChoice() FunctionCallChoice {
	return FunctionCallChoice(f.Name)
}

type FunctionCall struct {
	Name      string
	Parameter json.RawMessage
}

type FunctionCallChoice string

const (
	// FunctionCallChoiceNone means that the function call is not allowed.
	FunctionCallChoiceNone FunctionCallChoice = "none"
	// FunctionCallChoiceAuto means that the function call is automatically chosen.
	FunctionCallChoiceAuto FunctionCallChoice = "auto"
)

func (f FunctionCallChoice) String() string {
	return string(f)
}
