package workermailer

import (
	"coffee-like-helper-bot/config"
	"coffee-like-helper-bot/logger"
	"coffee-like-helper-bot/models"
	viewuser "coffee-like-helper-bot/view/user"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var log = logger.NewLoggerWithPrefix("mailer")

type Mailer struct {
	bot      *tgbotapi.BotAPI
	database *gorm.DB
	config   *config.Config
}

func New(bot *tgbotapi.BotAPI, database *gorm.DB, config *config.Config) *Mailer {
	w := &Mailer{
		bot:      bot,
		database: database,
		config:   config,
	}

	return w
}

func (w *Mailer) mail(users []models.User, msg *tgbotapi.MessageConfig) error {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for _, user := range users {
		<-t.C
		msg.BaseChat.ChatID = user.TelegramID
		_, err := w.bot.Send(msg)
		if err != nil {
			if strings.Contains(err.Error(), "bot was blocked by the user") { // TODO: придумать лучше, чем искать вхождение строки ожидаемой ошибки в строку фактической ошибки
				w.bot.Send(tgbotapi.NewMessage(1218219057, fmt.Sprint(viewuser.Text(&user), err)))
			} else {
				log.Println(err)
			}
		}
	}
	return nil
}

func (w *Mailer) AllUsers(msg *tgbotapi.MessageConfig) error {
	var allUsers []models.User
	err := w.database.Find(&allUsers).Error
	if err != nil {
		return err
	}

	return w.mail(allUsers, msg)
}

func (w *Mailer) Barista(msg *tgbotapi.MessageConfig) error {
	var baristas []models.User
	err := w.database.Find(&baristas, "is_barista = ?", true).Error
	if err != nil {
		return err
	}

	return w.mail(baristas, msg)
}

func (w *Mailer) Administrator(msg *tgbotapi.MessageConfig) error {
	var admins []models.User
	err := w.database.Find(&admins, "is_administrator = ?", true).Error
	if err != nil {
		return err
	}

	return w.mail(admins, msg)
}
