package services

import (
	"alfred-tool/database"
	"alfred-tool/models"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GetAllRsyncConfigs 获取所有rsync配置
func GetAllRsyncConfigs() ([]models.RsyncConfig, error) {
	var configs []models.RsyncConfig
	db := database.GetDB()
	err := db.Find(&configs).Error
	return configs, err
}

// GetRsyncConfigByName 根据名称获取rsync配置
func GetRsyncConfigByName(name string) (*models.RsyncConfig, error) {
	var config models.RsyncConfig
	db := database.GetDB()
	err := db.Where("name = ?", name).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// CreateRsyncConfig 创建rsync配置
func CreateRsyncConfig(config *models.RsyncConfig) error {
	db := database.GetDB()
	return db.Create(config).Error
}

// UpdateRsyncConfig 更新rsync配置
func UpdateRsyncConfig(config *models.RsyncConfig) error {
	db := database.GetDB()
	return db.Save(config).Error
}

// DeleteRsyncConfig 删除rsync配置
func DeleteRsyncConfig(name string) error {
	db := database.GetDB()
	return db.Where("name = ?", name).Delete(&models.RsyncConfig{}).Error
}

// SearchRsyncConfigs 搜索rsync配置
func SearchRsyncConfigs(query string) ([]models.RsyncConfig, error) {
	var configs []models.RsyncConfig
	db := database.GetDB()

	query = strings.TrimSpace(query)
	if query == "" {
		return GetAllRsyncConfigs()
	}

	err := db.Where("name LIKE ? OR ssh_name LIKE ? OR local_path LIKE ? OR remote_path LIKE ? OR description LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&configs).Error
	return configs, err
}

// ExecuteRsyncConfig 执行rsync配置
func ExecuteRsyncConfig(configName string) error {
	// 获取rsync配置
	config, err := GetRsyncConfigByName(configName)
	if err != nil {
		return fmt.Errorf("获取rsync配置失败: %v", err)
	}

	// 获取SSH连接信息
	sshConn, err := GetConnectionByName(config.SSHName)
	if err != nil {
		return fmt.Errorf("获取SSH连接失败: %v", err)
	}

	// 构建rsync命令
	cmdArgs := config.BuildRsyncCommand(sshConn)

	fmt.Printf("执行命令: %s\n", strings.Join(cmdArgs, " "))

	// 执行命令
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println("开始执行rsync...")
	err = cmd.Run()
	if err == nil {
		fmt.Println("rsync执行完成")
	}
	if err != nil {
		return fmt.Errorf("rsync执行失败: %v", err)
	}

	// 更新使用次数
	db := database.GetDB()
	db.Model(config).Update("usage_count", config.UsageCount+1)
	db.Model(sshConn).Update("usage_count", sshConn.UsageCount+1)
	return nil
}

// DryRunRsyncConfig 预览rsync命令（不执行）
func DryRunRsyncConfig(configName string) (string, error) {
	// 获取rsync配置
	config, err := GetRsyncConfigByName(configName)
	if err != nil {
		return "", fmt.Errorf("获取rsync配置失败: %v", err)
	}

	// 获取SSH连接信息
	sshConn, err := GetConnectionByName(config.SSHName)
	if err != nil {
		return "", fmt.Errorf("获取SSH连接失败: %v", err)
	}

	// 构建rsync命令
	cmdArgs := config.BuildRsyncCommand(sshConn)

	// 添加 --dry-run 参数进行预览
	cmdArgs = append(cmdArgs[:1], append([]string{"--dry-run"}, cmdArgs[1:]...)...)

	return strings.Join(cmdArgs, " "), nil
}

// ValidateRsyncConfig 验证rsync配置
func ValidateRsyncConfig(config *models.RsyncConfig) error {
	// 检查SSH连接是否存在
	_, err := GetConnectionByName(config.SSHName)
	if err != nil {
		return fmt.Errorf("SSH连接 '%s' 不存在", config.SSHName)
	}

	// 检查本地路径
	if config.Direction == models.RsyncDirectionUpload {
		if _, err := os.Stat(config.LocalPath); os.IsNotExist(err) {
			return fmt.Errorf("本地路径 '%s' 不存在", config.LocalPath)
		}
	}

	// 检查配置名称唯一性
	existingConfig, err := GetRsyncConfigByName(config.Name)
	if err == nil && existingConfig != nil && existingConfig.ID != config.ID {
		return fmt.Errorf("配置名称 '%s' 已存在", config.Name)
	}

	return nil
}
