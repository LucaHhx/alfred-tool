package dialog

import "alfred-tool/dialog/field"

type Dialog struct {
	WindowTitle  string                 `json:"windowTitle"`
	WindowWidth  int                    `json:"windowWidth"`
	WindowHeight int                    `json:"windowHeight"`
	OkLabel      string                 `json:"okLabel"`
	CancelLabel  string                 `json:"cancelLabel"`
	AlwaysOnTop  bool                   `json:"alwaysOnTop"`
	Fields       map[string]field.Field `json:"fields"`
}
