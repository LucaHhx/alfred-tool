package services

import (
	"errors"
	"fmt"
	"strings"

	"alfred-tool/database"
	"alfred-tool/models"
)

func SearchConnections(query string) ([]models.SSHConnection, error) {
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("查询条件不能为空")
	}

	db := database.GetDB()
	var connections []models.SSHConnection

	query = strings.TrimSpace(query)
	searchPattern := "%" + query + "%"

	err := db.Where("name LIKE ? OR address LIKE ?", searchPattern, searchPattern).
		Find(&connections).Error

	if err != nil {
		return nil, fmt.Errorf("搜索失败: %v", err)
	}

	return connections, nil
}

func ListAllConnections() ([]models.SSHConnection, error) {
	db := database.GetDB()
	var connections []models.SSHConnection

	err := db.Order("usage_count DESC").Find(&connections).Error
	if err != nil {
		return nil, fmt.Errorf("获取连接列表失败: %v", err)
	}

	return connections, nil
}

func GetConnectionByName(name string) (*models.SSHConnection, error) {
	db := database.GetDB()
	var connection models.SSHConnection

	err := db.Where("name = ?", name).First(&connection).Error
	if err != nil {
		return nil, fmt.Errorf("未找到连接: %s", name)
	}

	return &connection, nil
}

func UpdateConnection(conn *models.SSHConnection) error {
	db := database.GetDB()
	err := db.Save(conn).Error
	if err != nil {
		return fmt.Errorf("更新连接失败: %v", err)
	}
	return nil
}

func DeleteConnection(name string) error {
	db := database.GetDB()
	var connection models.SSHConnection

	err := db.Where("name = ?", name).First(&connection).Error
	if err != nil {
		return fmt.Errorf("未找到连接: %s", name)
	}

	err = db.Delete(&connection).Error
	if err != nil {
		return fmt.Errorf("删除连接失败: %v", err)
	}

	return nil
}

func IncrementUsageCount(name string) error {
	db := database.GetDB()
	var connection models.SSHConnection

	err := db.Where("name = ?", name).First(&connection).Error
	if err != nil {
		return nil
	}

	connection.UsageCount++
	err = db.Save(&connection).Error
	if err != nil {
		return fmt.Errorf("更新使用次数失败: %v", err)
	}

	return nil
}
