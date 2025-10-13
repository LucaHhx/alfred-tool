package dialog

import (
	"alfred-tool/dialog/field"
	"encoding/json"
	"os/exec"
)

type Options func(*Dialog)

func WithTitle(title string) Options {
	return func(d *Dialog) {
		d.WindowTitle = title
	}
}

func WithSize(width int, height int) Options {
	return func(d *Dialog) {
		d.WindowWidth = width
		d.WindowHeight = height
	}
}

func WithOkLabel(label string) Options {
	return func(d *Dialog) {
		d.OkLabel = label
	}
}

func WithCancelLabel(label string) Options {
	return func(d *Dialog) {
		d.CancelLabel = label
	}
}

func WithAlwaysOnTop(alwaysOnTop bool) Options {
	return func(d *Dialog) {
		d.AlwaysOnTop = alwaysOnTop
	}
}

func WithFields(fields ...field.Field) Options {
	return func(d *Dialog) {
		for _, item := range fields {
			item.Order = len(d.Fields) + 1
			d.Fields[item.BindingKey] = item
		}
	}
}

func NewDialog(opts ...Options) *Dialog {
	d := &Dialog{
		WindowTitle:  "对话框",
		WindowWidth:  500,
		WindowHeight: 400,
		OkLabel:      "确定",
		CancelLabel:  "取消",
		AlwaysOnTop:  false,
		Fields:       make(map[string]field.Field),
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func (d *Dialog) AddFields(fields ...field.Field) {
	for _, item := range fields {
		item.Order = len(d.Fields) + 1
		d.Fields[item.BindingKey] = item
	}
}

func (d *Dialog) Open() (map[string]any, error) {
	dialogJSON, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}
	// 运行 dialog 可执行文件，直接传递 JSON 字符串作为参数
	cmd := exec.Command("/Users/luca/github/LucaHhx/alfred-tool/dialog/dialog", string(dialogJSON))
	cmd.Dir = "/Users/luca/github/LucaHhx/alfred-tool/dialog"
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	value := make(map[string]any)

	err = json.Unmarshal(output, &value)
	if err != nil {
		return value, err
	}
	return value, nil
}
