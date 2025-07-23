package ui

import (
	"errors"
	"fmt"
	"alfred-tool/ui/component"
	"alfred-tool/ui/icons"
	"alfred-tool/ui/xtheme"
	"strings"

	"alfred-tool/database"
	"alfred-tool/models"
	"alfred-tool/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowRsyncDialog 显示rsync配置对话框
func ShowRsyncDialog(configName string) error {
	var config *models.RsyncConfig
	var isUpdateMode bool

	if configName != "" {
		var err error
		config, err = getRsyncConfigByName(configName)
		if err == nil && config != nil {
			isUpdateMode = true
		}
	}

	if !isUpdateMode {
		config = &models.RsyncConfig{}
	}

	myApp := app.New()
	myApp.Settings().SetTheme(&xtheme.XTheme{})

	var windowTitle string
	if isUpdateMode {
		windowTitle = "修改 Rsync 配置"
	} else {
		windowTitle = "添加 Rsync 配置"
	}
	myWindow := myApp.NewWindow(windowTitle)
	myWindow.Resize(fyne.NewSize(650, 650))
	myWindow.SetFixedSize(true)

	myWindow.CenterOnScreen()
	myWindow.RequestFocus()
	SetWindowAlwaysOnTop(myWindow, "sshd")

	// 创建输入控件
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("输入配置名称...")

	// SSH连接选择
	sshConnections, err := services.ListAllConnections()
	if err != nil {
		return fmt.Errorf("获取SSH连接失败: %v", err)
	}

	var sshNames []string
	for _, conn := range sshConnections {
		sshNames = append(sshNames, conn.Name)
	}

	if len(sshNames) == 0 {
		return errors.New("请先添加SSH连接")
	}

	sshSelect := widget.NewSelect(sshNames, nil)
	sshSelect.SetSelected(sshNames[0])

	// 传输方向选择
	directionSelect := widget.NewSelect([]string{"上传 (本地→服务器)", "下载 (服务器→本地)"}, nil)
	directionSelect.SetSelected("上传 (本地→服务器)")

	// 路径输入
	localPathEntry := widget.NewEntry()
	localPathEntry.SetPlaceHolder("输入本地路径...")

	localBrowseButton := widget.NewButton("", func() {
		path := component.ShowMacNativeFolderDialog()
		if path != "" {
			localPathEntry.SetText(path)
		}
	})
	localBrowseButton.SetIcon(icons.Icon_Ellipsis)

	remotePathEntry := widget.NewEntry()
	remotePathEntry.SetPlaceHolder("输入远程路径...")

	// 排除规则
	excludeEntry := widget.NewMultiLineEntry()
	excludeEntry.SetPlaceHolder("输入排除规则，每行一个:\n*.log\n*.tmp\n.DS_Store")
	excludeEntry.Resize(fyne.NewSize(0, 80))

	// 常用rsync选项复选框
	verboseCheck := widget.NewCheck("详细输出 (-v)", nil)
	recursiveCheck := widget.NewCheck("递归 (-r)", nil)
	archiveCheck := widget.NewCheck("归档模式 (-a)", nil)
	compressCheck := widget.NewCheck("压缩 (-z)", nil)
	timesCheck := widget.NewCheck("保持时间戳 (-t)", nil)
	progressCheck := widget.NewCheck("显示进度 (--progress)", nil)
	deleteCheck := widget.NewCheck("删除多余文件 (--delete)", nil)
	checksumCheck := widget.NewCheck("使用校验和 (-c)", nil)
	linksCheck := widget.NewCheck("复制符号链接 (-l)", nil)
	permsCheck := widget.NewCheck("保持权限 (-p)", nil)
	ownerCheck := widget.NewCheck("保持所有者 (-o)", nil)
	groupCheck := widget.NewCheck("保持组 (-g)", nil)

	// 设置默认值
	verboseCheck.SetChecked(true)
	archiveCheck.SetChecked(true)
	compressCheck.SetChecked(true)
	progressCheck.SetChecked(true)

	// 创建选项容器
	optionsContainer1 := container.NewGridWithColumns(3,
		verboseCheck, recursiveCheck, archiveCheck)
	optionsContainer2 := container.NewGridWithColumns(3,
		compressCheck, timesCheck, progressCheck)
	optionsContainer3 := container.NewGridWithColumns(3,
		deleteCheck, checksumCheck, linksCheck)
	optionsContainer4 := container.NewGridWithColumns(3,
		permsCheck, ownerCheck, groupCheck)

	// 额外选项
	optionsEntry := widget.NewEntry()
	optionsEntry.SetPlaceHolder("输入额外rsync选项，如: --backup")

	// 描述
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("输入描述信息...")
	descEntry.Resize(fyne.NewSize(0, 60))

	// 命令预览
	previewEntry := widget.NewMultiLineEntry()
	previewEntry.SetPlaceHolder("rsync命令预览将在这里显示...")
	previewEntry.Resize(fyne.NewSize(0, 100))
	previewEntry.Disable() // 设置为只读

	// 更新预览的函数
	updatePreview := func() {
		// 检查必要字段是否为空
		if sshSelect.Selected == "" || localPathEntry.Text == "" || remotePathEntry.Text == "" {
			previewEntry.SetText("请先填写SSH连接、本地路径和远程路径")
			return
		}

		// 获取SSH连接信息
		sshConn, err := services.GetConnectionByName(sshSelect.Selected)
		if err != nil {
			previewEntry.SetText("获取SSH连接失败: " + err.Error())
			return
		}

		// 创建临时配置用于生成预览
		var dir models.RsyncDirection
		if directionSelect.Selected == "上传 (本地→服务器)" {
			dir = models.RsyncDirectionUpload
		} else {
			dir = models.RsyncDirectionDownload
		}

		tempConfig := &models.RsyncConfig{
			Direction:    dir,
			LocalPath:    localPathEntry.Text,
			RemotePath:   remotePathEntry.Text,
			ExcludeRules: excludeEntry.Text,
			Options:      optionsEntry.Text,
			Verbose:      verboseCheck.Checked,
			Recursive:    recursiveCheck.Checked,
			Archive:      archiveCheck.Checked,
			Compress:     compressCheck.Checked,
			Times:        timesCheck.Checked,
			Progress:     progressCheck.Checked,
			Delete:       deleteCheck.Checked,
			Checksum:     checksumCheck.Checked,
			Links:        linksCheck.Checked,
			Perms:        permsCheck.Checked,
			Owner:        ownerCheck.Checked,
			Group:        groupCheck.Checked,
		}

		// 生成rsync命令
		cmdArgs := tempConfig.BuildRsyncCommand(sshConn)
		previewEntry.SetText(strings.Join(cmdArgs, " "))
	}

	// 如果是修改模式，填充现有数据
	if isUpdateMode && config != nil {
		nameEntry.SetText(config.Name)
		sshSelect.SetSelected(config.SSHName)

		if config.Direction == models.RsyncDirectionUpload {
			directionSelect.SetSelected("上传 (本地→服务器)")
		} else {
			directionSelect.SetSelected("下载 (服务器→本地)")
		}

		localPathEntry.SetText(config.LocalPath)
		remotePathEntry.SetText(config.RemotePath)
		excludeEntry.SetText(config.ExcludeRules)
		optionsEntry.SetText(config.Options)
		descEntry.SetText(config.Description)

		// 设置复选框状态
		verboseCheck.SetChecked(config.Verbose)
		recursiveCheck.SetChecked(config.Recursive)
		archiveCheck.SetChecked(config.Archive)
		compressCheck.SetChecked(config.Compress)
		timesCheck.SetChecked(config.Times)
		progressCheck.SetChecked(config.Progress)
		deleteCheck.SetChecked(config.Delete)
		checksumCheck.SetChecked(config.Checksum)
		linksCheck.SetChecked(config.Links)
		permsCheck.SetChecked(config.Perms)
		ownerCheck.SetChecked(config.Owner)
		groupCheck.SetChecked(config.Group)
	}

	// 为所有相关控件添加事件监听器，实时更新预览
	sshSelect.OnChanged = func(string) { updatePreview() }
	directionSelect.OnChanged = func(string) { updatePreview() }
	localPathEntry.OnChanged = func(string) { updatePreview() }
	remotePathEntry.OnChanged = func(string) { updatePreview() }
	excludeEntry.OnChanged = func(string) { updatePreview() }
	optionsEntry.OnChanged = func(string) { updatePreview() }

	verboseCheck.OnChanged = func(bool) { updatePreview() }
	recursiveCheck.OnChanged = func(bool) { updatePreview() }
	archiveCheck.OnChanged = func(bool) { updatePreview() }
	compressCheck.OnChanged = func(bool) { updatePreview() }
	timesCheck.OnChanged = func(bool) { updatePreview() }
	progressCheck.OnChanged = func(bool) { updatePreview() }
	deleteCheck.OnChanged = func(bool) { updatePreview() }
	checksumCheck.OnChanged = func(bool) { updatePreview() }
	linksCheck.OnChanged = func(bool) { updatePreview() }
	permsCheck.OnChanged = func(bool) { updatePreview() }
	ownerCheck.OnChanged = func(bool) { updatePreview() }
	groupCheck.OnChanged = func(bool) { updatePreview() }

	// 初始预览更新
	updatePreview()

	// 创建本地路径容器
	localPathContainer := container.NewBorder(nil, nil, nil, localBrowseButton, localPathEntry)

	// 创建表单
	form := widget.NewForm(
		widget.NewFormItem("配置名称", nameEntry),
		widget.NewFormItem("SSH连接", sshSelect),
		widget.NewFormItem("传输方向", directionSelect),
		widget.NewFormItem("本地路径", localPathContainer),
		widget.NewFormItem("远程路径", remotePathEntry),
		widget.NewFormItem("排除规则", excludeEntry),
		widget.NewFormItem("常用选项", container.NewVBox(
			optionsContainer1,
			optionsContainer2,
			optionsContainer3,
			optionsContainer4,
		)),
		widget.NewFormItem("额外选项", optionsEntry),
		widget.NewFormItem("描述", descEntry),
		widget.NewFormItem("命令预览", previewEntry),
	)

	// 设置表单提交和取消按钮
	form.OnSubmit = func() {
		var err error
		if isUpdateMode {
			err = updateRsyncConfig(config.ID, nameEntry.Text, sshSelect.Selected, directionSelect.Selected,
				localPathEntry.Text, remotePathEntry.Text, excludeEntry.Text, optionsEntry.Text, descEntry.Text,
				verboseCheck.Checked, recursiveCheck.Checked, archiveCheck.Checked, compressCheck.Checked,
				timesCheck.Checked, progressCheck.Checked, deleteCheck.Checked, checksumCheck.Checked,
				linksCheck.Checked, permsCheck.Checked, ownerCheck.Checked, groupCheck.Checked)
		} else {
			err = saveRsyncConfig(nameEntry.Text, sshSelect.Selected, directionSelect.Selected,
				localPathEntry.Text, remotePathEntry.Text, excludeEntry.Text, optionsEntry.Text, descEntry.Text,
				verboseCheck.Checked, recursiveCheck.Checked, archiveCheck.Checked, compressCheck.Checked,
				timesCheck.Checked, progressCheck.Checked, deleteCheck.Checked, checksumCheck.Checked,
				linksCheck.Checked, permsCheck.Checked, ownerCheck.Checked, groupCheck.Checked)
		}

		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		var successMsg string
		if isUpdateMode {
			successMsg = "Rsync 配置已更新"
		} else {
			successMsg = "Rsync 配置已保存"
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

// ShowAddRsyncDialog 显示添加rsync配置对话框
func ShowAddRsyncDialog() error {
	return ShowRsyncDialog("")
}

// ShowUpdateRsyncDialog 显示修改rsync配置对话框
func ShowUpdateRsyncDialog(name string) error {
	return ShowRsyncDialog(name)
}

func saveRsyncConfig(name, sshName, direction, localPath, remotePath, excludeRules, options, description string,
	verbose, recursive, archive, compress, times, progress, delete, checksum, links, perms, owner, group bool) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(sshName) == "" ||
		strings.TrimSpace(localPath) == "" || strings.TrimSpace(remotePath) == "" {
		return errors.New("配置名称、SSH连接、本地路径和远程路径不能为空")
	}

	// 转换传输方向
	var dir models.RsyncDirection
	if direction == "上传 (本地→服务器)" {
		dir = models.RsyncDirectionUpload
	} else {
		dir = models.RsyncDirectionDownload
	}

	config := models.RsyncConfig{
		Name:         strings.TrimSpace(name),
		SSHName:      strings.TrimSpace(sshName),
		Direction:    dir,
		LocalPath:    strings.TrimSpace(localPath),
		RemotePath:   strings.TrimSpace(remotePath),
		ExcludeRules: strings.TrimSpace(excludeRules),
		Options:      strings.TrimSpace(options),
		Description:  strings.TrimSpace(description),
		Verbose:      verbose,
		Recursive:    recursive,
		Archive:      archive,
		Compress:     compress,
		Times:        times,
		Progress:     progress,
		Delete:       delete,
		Checksum:     checksum,
		Links:        links,
		Perms:        perms,
		Owner:        owner,
		Group:        group,
	}

	db := database.GetDB()
	return db.Create(&config).Error
}

func updateRsyncConfig(id uint, name, sshName, direction, localPath, remotePath, excludeRules, options, description string,
	verbose, recursive, archive, compress, times, progress, delete, checksum, links, perms, owner, group bool) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(sshName) == "" ||
		strings.TrimSpace(localPath) == "" || strings.TrimSpace(remotePath) == "" {
		return errors.New("配置名称、SSH连接、本地路径和远程路径不能为空")
	}

	// 转换传输方向
	var dir models.RsyncDirection
	if direction == "上传 (本地→服务器)" {
		dir = models.RsyncDirectionUpload
	} else {
		dir = models.RsyncDirectionDownload
	}

	config := &models.RsyncConfig{
		Name:         strings.TrimSpace(name),
		SSHName:      strings.TrimSpace(sshName),
		Direction:    dir,
		LocalPath:    strings.TrimSpace(localPath),
		RemotePath:   strings.TrimSpace(remotePath),
		ExcludeRules: strings.TrimSpace(excludeRules),
		Options:      strings.TrimSpace(options),
		Description:  strings.TrimSpace(description),
		Verbose:      verbose,
		Recursive:    recursive,
		Archive:      archive,
		Compress:     compress,
		Times:        times,
		Progress:     progress,
		Delete:       delete,
		Checksum:     checksum,
		Links:        links,
		Perms:        perms,
		Owner:        owner,
		Group:        group,
	}
	config.ID = id

	db := database.GetDB()
	return db.Save(config).Error
}

func getRsyncConfigByName(name string) (*models.RsyncConfig, error) {
	var config models.RsyncConfig
	db := database.GetDB()
	err := db.Where("name = ?", name).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}
