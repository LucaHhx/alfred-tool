package dialog

import (
	"alfred-tool/dialog/field"
	"testing"
)

func TestDialog(t *testing.T) {
	open, err := NewDialog(
		WithTitle("完整功能测试对话框"),
		WithSize(800, 700),
		WithFields(
			// 文本输入框 - 带复制按钮
			field.NewTextField("username", "用户名",
				field.WithDefaultValue("admin"),
				field.WithCopy(true),
				field.WithNote("请输入您的用户名（必填）"),
				field.WithOrder(1),
			),

			// 多行文本编辑器
			field.NewTextEditorField("description", "描述",
				field.WithDefaultValue("请输入详细描述信息..."),
				field.WithNote("最多500字"),
				field.WithOrder(2),
			),

			// 复选框
			field.NewCheckBoxField("remember", "记住我",
				field.WithDefaultValue("true"),
				field.WithNote("勾选后下次自动登录"),
				field.WithOrder(3),
			),

			// 下拉选择框
			field.NewDropdownField("role", "用户角色",
				[]string{"管理员", "普通用户", "访客"},
				field.WithDefaultValue("普通用户"),
				field.WithNote("请选择您的角色"),
				field.WithOrder(4),
			),

			// 分段选择器
			field.NewSegmentedField("theme", "主题模式",
				[]string{"浅色", "深色", "自动"},
				field.WithDefaultValue("自动"),
				field.WithNote("选择应用的主题外观"),
				field.WithOrder(5),
			),

			// 文件选择器
			field.NewFileField("avatar", "头像文件",
				field.WithNote("支持 JPG、PNG 格式，大小不超过 2MB"),
				field.WithOrder(6),
				field.WithVisibleWhen("theme", "深色"),
			),

			// 文件夹选择器
			field.NewFolderField("output", "输出目录",
				field.WithDefaultValue("/Users"),
				field.WithNote("选择文件的保存位置"),
				field.WithOrder(7),
				field.WithVisibleWhen("theme", "自动"),
			),
		),
	).Open()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("对话框返回结果:", open)
}
