package field

type FieldType string

const (
	Text       FieldType = "text"       //文本
	TextEditor FieldType = "texteditor" //文本编辑器
	CheckBox   FieldType = "checkbox"   //多选框
	Dropdown   FieldType = "dropdown"   //下拉框
	Segmented  FieldType = "segmented"  // 分段器
	FilePicker FieldType = "filepicker" // 文件选择器
)

type FilePickerType string

const (
	Folder FilePickerType = "folder"
	File   FilePickerType = "file"
)

type Field struct {
	Type           FieldType      `json:"type"`
	Label          string         `json:"label"`
	BindingKey     string         `json:"bindingKey"`
	DefaultValue   string         `json:"defaultValue"`
	Copy           bool           `json:"copy,omitempty"`
	FilePickerType FilePickerType `json:"filePickerType,omitempty"`
	Options        []string       `json:"options,omitempty"`
	Note           string         `json:"note,omitempty"` // 字段备注，显示在控件下方的红色小字
	Order          int            `json:"order"`          // 字段显示顺序，数字越小越靠前
}
