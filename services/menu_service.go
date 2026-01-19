package services

import (
	"coffee_shop/dto"
	"coffee_shop/models"
	"coffee_shop/repositories"
	"errors"
)

type MenuService interface {
	// Define menu-related business logic methods here
	GetAllMenus() ([]dto.MenuResponse, error)
	CreateMenu(request models.Menu) (dto.MenuResponse, error)
	UpdateMenu(id uint, request models.Menu) (dto.MenuResponse, error)
	DeleteMenu(id uint) error
	GetMenuByID(id uint) (dto.MenuResponse, error)
}

type MenuServiceImpl struct {
	MenuRepo repositories.MenuRepository
}

// CreateMenu implements MenuService.
func (m *MenuServiceImpl) CreateMenu(request models.Menu) (dto.MenuResponse, error) {
	// Check if Menu is Exist
	findMenu, err := m.MenuRepo.FindByName(request.MenuName)
	if err == nil && findMenu != nil {
		return dto.MenuResponse{}, errors.New("menu name already exists: " + findMenu.MenuName)
	}

	requestMenu := models.Menu{
		MenuName:    request.MenuName,
		Price:       request.Price,
		Description: request.Description,
		ImageURL:    request.ImageURL,
		CategoryID:  request.CategoryID,
	}

	err = m.MenuRepo.CreateMenu(&requestMenu)
	if err != nil {
		return dto.MenuResponse{}, errors.New("failed to create menu: " + err.Error())
	}
	return dto.ToMenuResponse(&requestMenu), nil
}

// DeleteMenu implements MenuService.
func (m *MenuServiceImpl) DeleteMenu(id uint) error {
	err := m.MenuRepo.DeleteMenu(id)
	if err != nil {
		return errors.New("failed to delete menu: " + err.Error())
	}
	return nil
}

// GetAllMenus implements MenuService.
func (m *MenuServiceImpl) GetAllMenus() ([]dto.MenuResponse, error) {
	AllMenus, err := m.MenuRepo.GetAllMenus()
	if err != nil {
		return []dto.MenuResponse{}, err
	}

	var MenuResponses []dto.MenuResponse
	for _, menu := range AllMenus {
		MenuResponses = append(MenuResponses, dto.ToMenuResponse(&menu))
	}
	return MenuResponses, nil
}

// GetMenuByID implements MenuService.
func (m *MenuServiceImpl) GetMenuByID(id uint) (dto.MenuResponse, error) {
	menu, err := m.MenuRepo.GetMenuByID(id)
	if err != nil {
		return dto.MenuResponse{}, errors.New("failed to get menu by id: " + err.Error())
	}
	return dto.ToMenuResponse(menu), nil
}

// UpdateMenu implements MenuService.
func (m *MenuServiceImpl) UpdateMenu(id uint, request models.Menu) (dto.MenuResponse, error) {
	requestMenu := models.Menu{
		MenuName:    request.MenuName,
		Price:       request.Price,
		Description: request.Description,
		ImageURL:    request.ImageURL,
		CategoryID:  request.CategoryID,
	}

	err := m.MenuRepo.UpdateMenu(id, &requestMenu)
	if err != nil {
		return dto.MenuResponse{}, errors.New("failed to update menu: " + err.Error())
	}
	return dto.ToMenuResponse(&requestMenu), nil
}

func NewMenuService(menuRepo repositories.MenuRepository) MenuService {
	return &MenuServiceImpl{
		MenuRepo: menuRepo,
	}
}
