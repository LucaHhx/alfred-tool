package ssh

import (
	"alfred-tool/database"
	"alfred-tool/dialog"
	"alfred-tool/dialog/field"
	"alfred-tool/models"
	"alfred-tool/services"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ShowAddDialogV2 显示添加SSH连接对话框（使用dialog包）
func ShowAddDialogV2() error {
	d := dialog.NewDialog(
		dialog.WithTitle("添加 SSH 连接"),
		dialog.WithSize(600, 500),
		dialog.WithOkLabel("保存"),
		dialog.WithCancelLabel("取消"),
		dialog.WithAlwaysOnTop(true),
		dialog.WithFields(
			field.NewTextField("name", "连接名称"),
			field.NewTextField("address", "服务器地址", field.WithCopy(true)),
			field.NewTextField("port", "端口", field.WithDefaultValue("22")),
			field.NewTextField("username", "用户名", field.WithDefaultValue("root")),
			field.NewTextField("localIP", "局域网IP", field.WithCopy(true), field.WithNote("可选")),
			field.NewSegmentedField("passwordType", "密码类型", []string{"私钥", "密码"}, field.WithDefaultValue("私钥")),
			field.NewFileField("keyPath", "私钥文件", field.WithDefaultValue("/Users/luca/.ssh/key"), field.WithVisibleWhen("passwordType", "私钥")),
			field.NewTextField("password", "密码", field.WithVisibleWhen("passwordType", "密码")),
			field.NewTextEditorField("description", "描述", field.WithNote("可选")),
		),
	)

	result, err := d.Open()
	if err != nil {
		return fmt.Errorf("打开对话框失败: %v", err)
	}

	// 提取表单数据
	name := getStringValue(result, "name")
	address := getStringValue(result, "address")
	port := getStringValue(result, "port")
	username := getStringValue(result, "username")
	localIP := getStringValue(result, "localIP")
	passwordType := getStringValue(result, "passwordType")
	keyPath := getStringValue(result, "keyPath")
	password := getStringValue(result, "password")
	description := getStringValue(result, "description")

	// 保存连接
	err = saveConnection(name, address, port, username, localIP, passwordType, password, keyPath, description)
	if err != nil {
		return fmt.Errorf("保存连接失败: %v", err)
	}

	fmt.Printf("SSH 连接 '%s' 已成功添加\n", name)
	return nil
}

// ShowUpdateDialogV2 显示更新SSH连接对话框（使用dialog包）
func ShowUpdateDialogV2(connectionName string) error {
	// 获取现有连接信息
	conn, err := services.GetConnectionByName(connectionName)
	if err != nil {
		return fmt.Errorf("未找到连接 '%s': %v", connectionName, err)
	}

	// 根据密码类型确定默认值
	var passwordTypeDefault string
	if conn.PasswordType == models.PasswordTypePassword {
		passwordTypeDefault = "密码"
	} else {
		passwordTypeDefault = "私钥"
	}

	d := dialog.NewDialog(
		dialog.WithTitle("修改 SSH 连接"),
		dialog.WithSize(600, 500),
		dialog.WithOkLabel("更新"),
		dialog.WithCancelLabel("取消"),
		dialog.WithAlwaysOnTop(true),
		dialog.WithFields(
			field.NewTextField("name", "连接名称", field.WithDefaultValue(conn.Name)),
			field.NewTextField("address", "服务器地址", field.WithDefaultValue(conn.Address), field.WithCopy(true)),
			field.NewTextField("port", "端口", field.WithDefaultValue(strconv.Itoa(conn.Port))),
			field.NewTextField("username", "用户名", field.WithDefaultValue(conn.Username)),
			field.NewTextField("localIP", "局域网IP", field.WithDefaultValue(conn.LocalIP), field.WithCopy(true), field.WithNote("可选")),
			field.NewSegmentedField("passwordType", "密码类型", []string{"私钥", "密码"}, field.WithDefaultValue(passwordTypeDefault)),
			field.NewFileField("keyPath", "私钥文件", field.WithDefaultValue(conn.KeyPath), field.WithVisibleWhen("passwordType", "私钥")),
			field.NewTextField("password", "密码", field.WithDefaultValue(conn.Password), field.WithVisibleWhen("passwordType", "密码")),
			field.NewTextEditorField("description", "描述", field.WithDefaultValue(conn.Description), field.WithNote("可选")),
		),
	)

	result, err := d.Open()
	if err != nil {
		return fmt.Errorf("打开对话框失败: %v", err)
	}

	// 提取表单数据
	name := getStringValue(result, "name")
	address := getStringValue(result, "address")
	port := getStringValue(result, "port")
	username := getStringValue(result, "username")
	localIP := getStringValue(result, "localIP")
	passwordType := getStringValue(result, "passwordType")
	keyPath := getStringValue(result, "keyPath")
	password := getStringValue(result, "password")
	description := getStringValue(result, "description")

	// 更新连接
	err = updateConnection(conn.ID, name, address, port, username, localIP, passwordType, password, keyPath, description)
	if err != nil {
		return fmt.Errorf("更新连接失败: %v", err)
	}

	fmt.Printf("SSH 连接 '%s' 已成功更新\n", name)
	return nil
}

// getStringValue 从结果中安全获取字符串值
func getStringValue(result map[string]any, key string) string {
	if val, ok := result[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// saveConnection 保存SSH连接（复用原有逻辑）
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

// updateConnection 更新SSH连接（复用原有逻辑）
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
