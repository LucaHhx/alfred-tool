package ui

import (
	"alfred-tool/ui/xtheme"
	"errors"
	"fmt"
	"strings"

	"alfred-tool/models"
	"alfred-tool/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowServiceDialog 显示服务配置对话框
func ShowServiceDialog(serviceID uint) error {
	var service *models.Service
	var isUpdateMode bool

	serviceService := services.NewServiceService()

	if serviceID > 0 {
		var err error
		service, err = serviceService.GetServiceByID(serviceID)
		if err == nil && service != nil {
			isUpdateMode = true
		}
	}

	if !isUpdateMode {
		service = &models.Service{}
	}

	myApp := app.New()
	myApp.Settings().SetTheme(&xtheme.XTheme{})

	var windowTitle string
	if isUpdateMode {
		windowTitle = "修改服务配置"
	} else {
		windowTitle = "添加服务配置"
	}
	myWindow := myApp.NewWindow(windowTitle)
	myWindow.Resize(fyne.NewSize(650, 500))
	myWindow.SetFixedSize(false)

	myWindow.CenterOnScreen()
	myWindow.RequestFocus()
	SetWindowAlwaysOnTop(myWindow, "sshd")

	// 创建输入控件
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("输入服务名称...")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("输入服务描述...")
	descEntry.Resize(fyne.NewSize(0, 60))

	detailsEntry := widget.NewMultiLineEntry()
	detailsEntry.SetPlaceHolder("输入服务详细信息，支持多行文本...")
	detailsEntry.Resize(fyne.NewSize(0, 200))
	detailsEntry.Wrapping = fyne.TextWrapWord

	// SSH连接选择
	sshConnections, err := services.ListAllConnections()
	if err != nil {
		return fmt.Errorf("获取SSH连接失败: %v", err)
	}

	var sshNames []string
	sshNames = append(sshNames, "不关联SSH连接")
	for _, conn := range sshConnections {
		sshNames = append(sshNames, fmt.Sprintf("%s (%s@%s)", conn.Name, conn.Username, conn.Address))
	}

	sshSelect := widget.NewSelect(sshNames, nil)
	sshSelect.SetSelected("不关联SSH连接")

	// 如果是修改模式，填充现有数据
	if isUpdateMode && service != nil {
		nameEntry.SetText(service.Name)
		descEntry.SetText(service.Description)
		detailsEntry.SetText(service.Details)

		// 设置SSH连接选择
		if service.SSHConnectionID > 0 {
			for _, conn := range sshConnections {
				if conn.ID == service.SSHConnectionID {
					sshSelect.SetSelected(fmt.Sprintf("%s (%s@%s)", conn.Name, conn.Username, conn.Address))
					break
				}
			}
		}
	}

	// 创建表单
	form := widget.NewForm(
		widget.NewFormItem("服务名称", nameEntry),
		widget.NewFormItem("关联SSH连接", sshSelect),
		widget.NewFormItem("服务描述", descEntry),
		widget.NewFormItem("服务详情", detailsEntry),
	)

	// 设置表单提交和取消按钮
	form.OnSubmit = func() {
		var err error

		// 获取SSH连接ID
		var sshConnectionID uint = 0
		if sshSelect.Selected != "不关联SSH连接" {
			for _, conn := range sshConnections {
				selectedText := fmt.Sprintf("%s (%s@%s)", conn.Name, conn.Username, conn.Address)
				if sshSelect.Selected == selectedText {
					sshConnectionID = conn.ID
					break
				}
			}
		}

		if isUpdateMode {
			err = updateService(service.ID, nameEntry.Text, descEntry.Text, detailsEntry.Text, sshConnectionID)
		} else {
			err = saveService(nameEntry.Text, descEntry.Text, detailsEntry.Text, sshConnectionID)
		}

		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		var successMsg string
		if isUpdateMode {
			successMsg = "服务配置已更新"
		} else {
			successMsg = "服务配置已保存"
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

// ShowAddServiceDialog 显示添加服务配置对话框
func ShowAddServiceDialog() error {
	return ShowServiceDialog(0)
}

// ShowUpdateServiceDialog 显示修改服务配置对话框
func ShowUpdateServiceDialog(serviceID uint) error {
	return ShowServiceDialog(serviceID)
}

func saveService(name, description, details string, sshConnectionID uint) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("服务名称不能为空")
	}

	service := models.Service{
		Name:            strings.TrimSpace(name),
		Description:     strings.TrimSpace(description),
		Details:         strings.TrimSpace(details),
		SSHConnectionID: sshConnectionID,
	}

	serviceService := services.NewServiceService()
	return serviceService.CreateService(&service)
}

func updateService(id uint, name, description, details string, sshConnectionID uint) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("服务名称不能为空")
	}

	service := &models.Service{
		Name:            strings.TrimSpace(name),
		Description:     strings.TrimSpace(description),
		Details:         strings.TrimSpace(details),
		SSHConnectionID: sshConnectionID,
	}
	service.ID = id

	serviceService := services.NewServiceService()
	return serviceService.UpdateService(service)
}
