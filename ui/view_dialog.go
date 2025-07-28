package ui

import (
	"alfred-tool/ui/component"
	"alfred-tool/ui/icons"
	"alfred-tool/ui/xtheme"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"alfred-tool/database"
	"alfred-tool/models"
	"alfred-tool/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

// ShowViewDialog 显示统一的SSH连接对话框
// connectionName: 连接名称，如果数据库中存在则为修改模式，否则为新增模式
func ShowViewDialog(connectionName string) (err error) {
	var conn *models.SSHConnection
	var isUpdateMode bool

	conn, err = services.GetConnectionByName(connectionName)
	if err == nil && conn != nil {
		// 找到连接，进入修改模式
		isUpdateMode = true
	} else {
		// 未找到连接，进入新增模式
		isUpdateMode = false
		conn = nil
	}

	myApp := app.New()
	myApp.SetIcon(icons.Icon_Copy)
	myApp.Settings().SetTheme(&xtheme.XTheme{})
	// 根据模式设置窗口标题
	var windowTitle string
	if isUpdateMode {
		windowTitle = "修改 SSH 连接"
	} else {
		windowTitle = "添加 SSH 连接"
	}
	myWindow := myApp.NewWindow(windowTitle)
	myWindow.Resize(fyne.NewSize(502, 355))
	myWindow.SetFixedSize(true)
	myWindow.SetIcon(icons.Icon_Query)

	// 设置窗口居中并保持在前台
	myWindow.CenterOnScreen()
	myWindow.RequestFocus()

	// 设置窗口置顶 (macOS)
	SetWindowAlwaysOnTop(myWindow, "sshd")

	// 创建输入控件
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("输入连接名称...")

	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("输入服务器地址...")

	portEntry := widget.NewEntry()
	portEntry.SetText("22")
	portEntry.SetPlaceHolder("端口号...")

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("输入用户名...")

	localIPEntry := widget.NewEntry()
	localIPEntry.SetPlaceHolder("输入局域网IP地址（可选）...")

	passwordTypeSelect := widget.NewSelect([]string{"密码", "私钥"}, nil)
	passwordTypeSelect.SetSelected("私钥")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("输入密码...")
	passwordEntry.Hide()

	keyPathEntry := widget.NewEntry()
	keyPathEntry.SetPlaceHolder("选择私钥文件...")
	//keyPathEntry.Hide()

	browseButton := widget.NewButton("", func() {
		keyPathEntry.SetText(component.ShowMacNativeFileDialog())
	})
	browseButton.SetIcon(icons.Icon_Ellipsis)

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("输入描述信息...")
	descEntry.Resize(fyne.NewSize(0, 80))

	// 如果是修改模式，填充现有数据
	if isUpdateMode && conn != nil {
		nameEntry.SetText(conn.Name)
		addressEntry.SetText(conn.Address)
		portEntry.SetText(strconv.Itoa(conn.Port))
		usernameEntry.SetText(conn.Username)
		localIPEntry.SetText(conn.LocalIP)
		descEntry.SetText(conn.Description)

		if conn.PasswordType == models.PasswordTypePassword {
			passwordTypeSelect.SetSelected("密码")
			passwordEntry.SetText(conn.Password)
		} else {
			passwordTypeSelect.SetSelected("私钥")
			keyPathEntry.SetText(conn.KeyPath)
			passwordEntry.Hide()
			keyPathEntry.Show()
			browseButton.Show()
		}
	} else {
		//nameEntry.SetText(conn.Name)
		//addressEntry.SetText(conn.Address)
		portEntry.SetText("22")
		usernameEntry.SetText("root")
		//localIPEntry.SetText(conn.LocalIP)
		//descEntry.SetText(conn.Description)
		keyPathEntry.SetText("/Users/luca/.ssh/key")
	}

	// 密码类型选择处理
	passwordTypeSelect.OnChanged = func(value string) {
		if value == "密码" {
			passwordEntry.Show()
			keyPathEntry.Hide()
			browseButton.Hide()
		} else {
			passwordEntry.Hide()
			keyPathEntry.Show()
			browseButton.Show()
		}
	}

	// 创建复制按钮用于服务器地址
	addressCopyButton := widget.NewButton("", func() {
		err := clipboard.WriteAll(addressEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("复制失败: %v", err), myWindow)
			return
		}
		//dialog.ShowInformation("成功", "服务器地址已复制到剪贴板", myWindow)
	})
	addressCopyButton.SetIcon(icons.Icon_Copy)

	// 创建复制按钮用于局域网IP
	localIPCopyButton := widget.NewButton("", func() {
		if strings.TrimSpace(localIPEntry.Text) == "" {
			dialog.ShowError(fmt.Errorf("局域网IP为空"), myWindow)
			return
		}
		err := clipboard.WriteAll(localIPEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("复制失败: %v", err), myWindow)
			return
		}
		//dialog.ShowInformation("成功", "局域网IP已复制到剪贴板", myWindow)
	})
	localIPCopyButton.SetIcon(icons.Icon_Copy)

	// 创建带复制按钮的服务器地址容器
	addressContainer := container.NewBorder(nil, nil, nil, addressCopyButton, addressEntry)
	addressPortContainer := container.NewGridWithColumns(2, addressContainer, container.NewBorder(nil, nil, widget.NewLabel(":"), nil, portEntry))

	// 创建带复制按钮的局域网IP容器
	localIPContainer := container.NewBorder(nil, nil, nil, localIPCopyButton, localIPEntry)

	// 创建文件选择组合框
	fileContainer := container.NewBorder(nil, nil, nil, browseButton, keyPathEntry)

	// 创建密码输入区域
	passwordContainer := container.NewVBox(
		passwordEntry,
		fileContainer,
	)

	// 创建表单
	form := widget.NewForm(
		widget.NewFormItem("连接名称", nameEntry),
		widget.NewFormItem("服务器地址", addressPortContainer),
		widget.NewFormItem("用户名", usernameEntry),
		widget.NewFormItem("局域网IP", localIPContainer),
		widget.NewFormItem("密码类型", passwordTypeSelect),
		widget.NewFormItem("认证信息", passwordContainer),
		widget.NewFormItem("描述", descEntry),
	)

	// 设置表单提交和取消按钮
	form.OnSubmit = func() {
		var err error
		if isUpdateMode {
			err = updateConnection(conn.ID, nameEntry.Text, addressEntry.Text, portEntry.Text,
				usernameEntry.Text, localIPEntry.Text, passwordTypeSelect.Selected, passwordEntry.Text,
				keyPathEntry.Text, descEntry.Text)
		} else {
			err = saveConnection(nameEntry.Text, addressEntry.Text, portEntry.Text,
				usernameEntry.Text, localIPEntry.Text, passwordTypeSelect.Selected, passwordEntry.Text,
				keyPathEntry.Text, descEntry.Text)
		}

		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		var successMsg string
		if isUpdateMode {
			successMsg = "SSH 连接已更新"
		} else {
			successMsg = "SSH 连接已保存"
		}
		fmt.Println(nameEntry.Text, successMsg)
		dialog.ShowInformation("成功", successMsg, myWindow)
		myWindow.Close()
	}

	form.OnCancel = func() {
		var cancelMsg string
		if isUpdateMode {
			cancelMsg = "修改操作已取消"
		} else {
			cancelMsg = "添加操作已取消"
		}
		// 取消操作，关闭窗口
		fmt.Println(nameEntry.Text, cancelMsg)
		myWindow.Close()
	}

	// 设置按钮文本
	if isUpdateMode {
		form.SubmitText = "更新"
	} else {
		form.SubmitText = "保存"
	}
	form.CancelText = "取消"

	// 创建内容布局
	content := container.NewScroll(form)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()

	return nil
}

// ShowAddDialog 显示添加SSH连接对话框
func ShowAddDialog() {
	ShowViewDialog("")
}

// ShowUpdateDialog 显示修改SSH连接对话框
func ShowUpdateDialog(name string) error {
	return ShowViewDialog(name)
}

func saveConnection(name, address, port, username, localIP, passwordType, password, keyPath, description string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(address) == "" ||
		strings.TrimSpace(username) == "" {
		return errors.New("名称、地址和用户名不能为空")
	}

	// 如果端口为空，使用默认端口22
	if strings.TrimSpace(port) == "" {
		port = "22"
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return errors.New("端口号无效")
	}

	// 将中文类型转换为英文存储
	var dbPasswordType models.PasswordType
	if passwordType == "密码" {
		dbPasswordType = models.PasswordTypePassword
	} else {
		dbPasswordType = models.PasswordTypeKeyPath
	}

	conn := models.SSHConnection{
		Name:         strings.TrimSpace(name),
		Address:      strings.TrimSpace(address),
		Port:         portNum,
		Username:     strings.TrimSpace(username),
		PasswordType: dbPasswordType,
		LocalIP:      strings.TrimSpace(localIP),
		Description:  strings.TrimSpace(description),
	}

	if passwordType == "密码" {
		conn.Password = password
	} else {
		conn.KeyPath = strings.TrimSpace(keyPath)
	}

	db := database.GetDB()
	return db.Create(&conn).Error
}

func updateConnection(id uint, name, address, port, username, localIP, passwordType, password, keyPath, description string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(address) == "" ||
		strings.TrimSpace(username) == "" {
		return errors.New("名称、地址和用户名不能为空")
	}

	if strings.TrimSpace(port) == "" {
		port = "22"
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return errors.New("端口号无效")
	}

	var dbPasswordType models.PasswordType
	if passwordType == "密码" {
		dbPasswordType = models.PasswordTypePassword
	} else {
		dbPasswordType = models.PasswordTypeKeyPath
	}

	conn := &models.SSHConnection{
		Name:         strings.TrimSpace(name),
		Address:      strings.TrimSpace(address),
		Port:         portNum,
		Username:     strings.TrimSpace(username),
		PasswordType: dbPasswordType,
		LocalIP:      strings.TrimSpace(localIP),
		Description:  strings.TrimSpace(description),
	}
	conn.ID = id

	if passwordType == "密码" {
		conn.Password = password
		conn.KeyPath = "" // 清除密钥路径
	} else {
		conn.KeyPath = strings.TrimSpace(keyPath)
		conn.Password = "" // 清除密码
	}

	return services.UpdateConnection(conn)
}
