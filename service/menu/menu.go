package menuservice

import (
	"coffee-like-helper-bot/models"
	"errors"

	"gorm.io/gorm"
)

var (
	errMenuNotFound = errors.New("меню не найдено")
)

type MenuService struct {
	db *gorm.DB
}

func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{
		db: db,
	}
}

func (s *MenuService) MenuByID(id uint) (*models.Menu, error) {
	var menu models.Menu
	if err := s.db.First(&menu, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errMenuNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (s *MenuService) NewMenu(menu *models.Menu) error {
	if err := s.db.Create(menu).Error; err != nil {
		return err
	}
	return nil
}
