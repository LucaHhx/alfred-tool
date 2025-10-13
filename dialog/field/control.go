package field

// FieldOption 是用于配置 Field 的函数式选项
type FieldOption func(*Field)

// WithLabel 设置字段的显示标签
// 参数:
//   - label: 在对话框中显示的字段名称
//
// 示例:
//
//	WithLabel("用户名")
func WithLabel(label string) FieldOption {
	return func(f *Field) {
		f.Label = label
	}
}

// WithBindingKey 设置字段的绑定键名
// 参数:
//   - bindingKey: 用于在返回结果 JSON 中标识该字段值的键名
//
// 示例:
//
//	WithBindingKey("username")
func WithBindingKey(bindingKey string) FieldOption {
	return func(f *Field) {
		f.BindingKey = bindingKey
	}
}

// WithDefaultValue 设置字段的默认值
// 参数:
//   - defaultValue: 字段的默认值
//     对于文本框：任意字符串
//     对于复选框：使用 "true" 或 "false"
//     对于下拉框/分段器：应该是 options 中的一个值
//
// 示例:
//
//	WithDefaultValue("admin")           // 文本框
//	WithDefaultValue("true")            // 复选框
//	WithDefaultValue("管理员")           // 下拉框
func WithDefaultValue(defaultValue string) FieldOption {
	return func(f *Field) {
		f.DefaultValue = defaultValue
	}
}

// WithCopy 设置是否在字段旁显示复制按钮
// 参数:
//   - copy: true 显示复制按钮，false 不显示
//
// 注意: 仅对 Text 和 TextEditor 类型字段有效
//
// 示例:
//
//	WithCopy(true)  // 显示复制按钮，方便用户复制字段内容
func WithCopy(copy bool) FieldOption {
	return func(f *Field) {
		f.Copy = copy
	}
}

// WithFilePickerType 设置文件选择器的类型
// 参数:
//   - filePickerType: File（选择文件）或 Folder（选择文件夹）
//
// 注意: 仅对 FilePicker 类型字段有效
//
// 示例:
//
//	WithFilePickerType(field.File)      // 选择文件
//	WithFilePickerType(field.Folder)    // 选择文件夹
func WithFilePickerType(filePickerType FilePickerType) FieldOption {
	return func(f *Field) {
		f.FilePickerType = filePickerType
	}
}

// WithOptions 设置下拉框或分段器的选项列表
// 参数:
//   - options: 字符串数组，表示可供选择的选项
//
// 注意: 仅对 Dropdown 和 Segmented 类型字段有效
//
// 示例:
//
//	WithOptions([]string{"管理员", "用户", "访客"})
func WithOptions(options []string) FieldOption {
	return func(f *Field) {
		f.Options = options
	}
}

// WithNote 设置字段下方显示的备注文字
// 参数:
//   - note: 备注文字内容，将以红色小字显示在字段控件下方
//
// 示例:
//
//	WithNote("请输入6-20位字符")
//	WithNote("此操作不可撤销")
func WithNote(note string) FieldOption {
	return func(f *Field) {
		f.Note = note
	}
}

// WithOrder 设置字段的显示顺序
// 参数:
//   - order: 排序值，数字越小越靠前显示
//
// 示例:
//
//	WithOrder(1)  // 第一个显示
//	WithOrder(2)  // 第二个显示
func WithOrder(order int) FieldOption {
	return func(f *Field) {
		f.Order = order
	}
}

func NewField(fieldType FieldType, opts ...FieldOption) Field {
	f := Field{
		Type: fieldType,
	}
	for _, opt := range opts {
		opt(&f)
	}
	return f
}

// NewTextField 创建一个文本输入框字段
// 参数:
//   - key: 字段绑定的键名（bindingKey），用于在结果中标识该字段的值
//   - name: 字段显示的标签名称（label）
//   - opts: 可选的字段配置项，如 WithDefaultValue、WithCopy、WithNote、WithOrder 等
//
// 示例:
//
//	field.NewTextField("username", "用户名", WithDefaultValue("admin"), WithCopy(true), WithNote("请输入用户名"))
func NewTextField(key, name string, opts ...FieldOption) Field {
	if opts == nil {
		opts = make([]FieldOption, 0)
	}
	opts = append(opts, WithBindingKey(key), WithLabel(name))
	return NewField(Text, opts...)
}

// NewTextEditorField 创建一个多行文本编辑器字段
// 参数:
//   - key: 字段绑定的键名（bindingKey）
//   - name: 字段显示的标签名称（label）
//   - opts: 可选的字段配置项，如 WithDefaultValue、WithCopy、WithNote、WithOrder 等
//
// 示例:
//
//	field.NewTextEditorField("description", "描述", WithDefaultValue("请输入详细描述"), WithNote("最多500字"))
func NewTextEditorField(key, name string, opts ...FieldOption) Field {
	if opts == nil {
		opts = make([]FieldOption, 0)
	}
	opts = append(opts, WithBindingKey(key), WithLabel(name))
	return NewField(TextEditor, opts...)
}

// NewCheckBoxField 创建一个复选框字段
// 参数:
//   - key: 字段绑定的键名（bindingKey）
//   - name: 字段显示的标签名称（label）
//   - opts: 可选的字段配置项，如 WithDefaultValue（"true"/"false"）、WithNote、WithOrder 等
//
// 示例:
//
//	field.NewCheckBoxField("remember", "记住我", WithDefaultValue("true"), WithNote("勾选后下次自动登录"))
func NewCheckBoxField(key, name string, opts ...FieldOption) Field {
	if opts == nil {
		opts = make([]FieldOption, 0)
	}
	opts = append(opts, WithBindingKey(key), WithLabel(name))
	return NewField(CheckBox, opts...)
}

// NewDropdownField 创建一个下拉选择框字段
// 参数:
//   - key: 字段绑定的键名（bindingKey）
//   - name: 字段显示的标签名称（label）
//   - options: 下拉框的选项列表
//   - opts: 可选的字段配置项，如 WithDefaultValue、WithNote、WithOrder 等
//
// 示例:
//
//	field.NewDropdownField("role", "角色", []string{"管理员", "用户", "访客"}, WithDefaultValue("用户"), WithNote("请选择您的角色"))
func NewDropdownField(key, name string, options []string, opts ...FieldOption) Field {
	if opts == nil {
		opts = make([]FieldOption, 0)
	}
	opts = append(opts, WithBindingKey(key), WithLabel(name), WithOptions(options))
	return NewField(Dropdown, opts...)
}

// NewSegmentedField 创建一个分段选择器字段（单选按钮组）
// 参数:
//   - key: 字段绑定的键名（bindingKey）
//   - name: 字段显示的标签名称（label）
//   - options: 分段选项列表
//   - opts: 可选的字段配置项，如 WithDefaultValue、WithNote、WithOrder 等
//
// 示例:
//
//	field.NewSegmentedField("theme", "主题", []string{"浅色", "深色", "自动"}, WithDefaultValue("自动"), WithNote("选择应用主题"))
func NewSegmentedField(key, name string, options []string, opts ...FieldOption) Field {
	if opts == nil {
		opts = make([]FieldOption, 0)
	}
	opts = append(opts, WithBindingKey(key), WithLabel(name), WithOptions(options))
	return NewField(Segmented, opts...)
}

// NewFilePickerField 创建一个文件/文件夹选择器字段
// 参数:
//   - key: 字段绑定的键名（bindingKey）
//   - name: 字段显示的标签名称（label）
//   - filePickerType: 选择器类型，File（选择文件）或 Folder（选择文件夹）
//   - opts: 可选的字段配置项，如 WithDefaultValue、WithNote、WithOrder 等
//
// 示例:
//
//	field.NewFilePickerField("config", "配置文件", field.File, WithNote("请选择配置文件"))
func NewFilePickerField(key, name string, filePickerType FilePickerType, opts ...FieldOption) Field {
	if opts == nil {
		opts = make([]FieldOption, 0)
	}
	opts = append(opts, WithBindingKey(key), WithLabel(name), WithFilePickerType(filePickerType))
	return NewField(FilePicker, opts...)
}

// NewFileField 创建一个文件选择器字段（NewFilePickerField 的便捷方法）
// 参数:
//   - key: 字段绑定的键名（bindingKey）
//   - name: 字段显示的标签名称（label）
//   - opts: 可选的字段配置项，如 WithDefaultValue、WithNote、WithOrder 等
//
// 示例:
//
//	field.NewFileField("avatar", "头像", WithNote("支持 JPG/PNG 格式"))
func NewFileField(key, name string, opts ...FieldOption) Field {
	if opts == nil {
		opts = make([]FieldOption, 0)
	}
	return NewFilePickerField(key, name, File, opts...)
}

// NewFolderField 创建一个文件夹选择器字段（NewFilePickerField 的便捷方法）
// 参数:
//   - key: 字段绑定的键名（bindingKey）
//   - name: 字段显示的标签名称（label）
//   - opts: 可选的字段配置项，如 WithDefaultValue、WithNote、WithOrder 等
//
// 示例:
//
//	field.NewFolderField("output", "输出目录", WithNote("选择文件的保存位置"))
func NewFolderField(key, name string, opts ...FieldOption) Field {
	if opts == nil {
		opts = make([]FieldOption, 0)
	}
	return NewFilePickerField(key, name, Folder, opts...)
}
