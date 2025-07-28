package services

import (
	"alfred-tool/database"
	"alfred-tool/models"
	"errors"
	"strings"
)

type ServiceService struct{}

func NewServiceService() *ServiceService {
	return &ServiceService{}
}

func (s *ServiceService) CreateService(service *models.Service) error {
	if service.Name == "" {
		return errors.New("服务名称不能为空")
	}

	var existingService models.Service
	if err := database.GetDB().Where("name = ?", service.Name).First(&existingService).Error; err == nil {
		return errors.New("已存在同名服务")
	}

	return database.GetDB().Create(service).Error
}

func (s *ServiceService) GetServiceByID(id uint) (*models.Service, error) {
	var service models.Service
	if err := database.GetDB().Preload("SSHConnection").First(&service, id).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *ServiceService) GetAllServices() ([]models.Service, error) {
	var services []models.Service
	if err := database.GetDB().Preload("SSHConnection").Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (s *ServiceService) GetServicesBySSHConnection(sshConnectionID uint) ([]models.Service, error) {
	var services []models.Service
	if err := database.GetDB().Preload("SSHConnection").Where("ssh_connection_id = ?", sshConnectionID).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (s *ServiceService) SearchServices(keyword string) ([]models.Service, error) {
	var services []models.Service
	searchPattern := "%" + strings.ToLower(keyword) + "%"

	if err := database.GetDB().Preload("SSHConnection").Where(
		"LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(details) LIKE ?",
		searchPattern, searchPattern, searchPattern,
	).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (s *ServiceService) UpdateService(service *models.Service) error {
	if service.ID == 0 {
		return errors.New("服务ID不能为空")
	}
	if service.Name == "" {
		return errors.New("服务名称不能为空")
	}

	var existingService models.Service
	if err := database.GetDB().Where("name = ? AND id != ?", service.Name, service.ID).First(&existingService).Error; err == nil {
		return errors.New("已存在同名服务")
	}

	return database.GetDB().Save(service).Error
}

func (s *ServiceService) DeleteService(id uint) error {
	return database.GetDB().Delete(&models.Service{}, id).Error
}
