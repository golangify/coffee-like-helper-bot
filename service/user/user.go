package userservice

import (
	"coffee-like-helper-bot/models"
	"errors"
	"slices"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var whitelist = []string{
	"golangify",
	"zhukitaa",
	"valyuhin",
	"nikolaeva136",
}

var whitelistString = strings.Join(whitelist, ", ")

func isWhitelisted(username string) bool {
	return slices.Contains(whitelist, username)
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

func (s *UserService) UpdateUserInfo(user *models.User, tgProfile *tgbotapi.User) error {
	var updated bool

	if user.FirstName != tgProfile.FirstName {
		user.FirstName = tgProfile.FirstName
		updated = true
	}
	if user.LastName != tgProfile.LastName {
		user.LastName = tgProfile.LastName
		updated = true
	}
	if user.UserName != tgProfile.UserName {
		user.UserName = tgProfile.UserName
		updated = true
	}

	if updated {
		err := s.db.Model(&user).UpdateColumns(map[string]any{"first_name": user.FirstName, "last_name": user.LastName, "user_name": user.UserName}).Error
		return err
	}

	return nil
}

func (s *UserService) NewUser(user *models.User) error {
	if err := s.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) BanUser(user *models.User) error {
	if isWhitelisted(user.UserName) {
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
	if !isWhitelisted(user.UserName) {
		return errors.New("к сожалению только эти администраторы могут назначать новых администраторов: " + whitelistString)
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
	if isWhitelisted(user.UserName) {
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
