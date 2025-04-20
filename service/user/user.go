package userservice

import (
	"coffee-like-helper-bot/models"
	"errors"
	"slices"
	"strings"

	"gorm.io/gorm"
)

var whitelist = []string{
	"golangify",
	"zhukitaa",
	"valyuhin",
	"nikolaeva136",
}

var (
	errUserNotFound                = errors.New("пользователь не найден")
	errUserAlreadyBanned           = errors.New("пользователь уже забанен")
	errUserAlreadyUnbanned         = errors.New("пользователь не в бане")
	errUserAlreadyBarista          = errors.New("пользователь уже является бариста")
	errUserAlreadyNotBarista       = errors.New("пользователь и так не является бариста")
	errUserAlreadyNotAdministrator = errors.New("пользователь уже не является администратором")
	errUserAlreadyAdministrator    = errors.New("пользователь уже является администратором")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) UserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UserByTgID(tgID int64) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "telegram_id = ?", tgID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) BanUser(user *models.User) error {
	if slices.Contains(whitelist, user.UserName) {
		return errors.New("бан этого пользователя невозможен")
	}
	if user.IsBanned {
		return errUserAlreadyBanned
	}
	user.IsBanned = true
	if err := s.db.Model(&user).Update("is_banned", user.IsBanned).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) UnbanUser(user *models.User) error {
	if !user.IsBanned {
		return errUserAlreadyUnbanned
	}
	user.IsBanned = true
	if err := s.db.Model(&user).Update("is_banned", user.IsBanned).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) MakeBarista(user *models.User) error {
	if user.IsBarista {
		return errUserAlreadyBarista
	}
	user.IsBarista = true
	if err := s.db.Model(&user).Update("is_barista", user.IsBarista).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) RemoveBarista(user *models.User) error {
	if !user.IsBarista {
		return errUserAlreadyNotBarista
	}
	user.IsBarista = false
	if err := s.db.Model(&user).Update("is_barista", user.IsBarista).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) MakeAdministrator(initiator *models.User, user *models.User) error {
	if !slices.Contains(whitelist, initiator.UserName) {
		return errors.New("к сожалению только эти администраторы могут назначать новых администраторов: " + strings.Join(whitelist, ", "))
	}
	if user.IsAdministrator {
		return errUserAlreadyAdministrator
	}
	user.IsAdministrator = true
	if err := s.db.Model(&user).Update("is_administrator", user.IsAdministrator).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) RemoveAdministrator(user *models.User) error {
	if slices.Contains(whitelist, user.UserName) {
		return errors.New("нельзя удалить этого пользователя из администраторов")
	}
	if !user.IsAdministrator {
		return errUserAlreadyNotAdministrator
	}
	user.IsAdministrator = false
	if err := s.db.Model(&user).Update("is_administrator", user.IsAdministrator).Error; err != nil {
		return err
	}
	return nil
}
