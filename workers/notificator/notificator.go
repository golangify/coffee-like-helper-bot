package workernotificator

import (
	"coffee-like-helper-bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type Notification struct {
	UserID uint
	Type string // "all" || "barista" || "administrator"
	WeekDays string // []int
	Text string
}

type Notificator struct {
	bot *tgbotapi.BotAPI
	database *gorm.DB
	config *config.Config
}

func (w *Notificator)run() {
	//
}

func NewNotificator(bot *tgbotapi.BotAPI, database *gorm.DB, config *config.Config) (*Notificator, error) {
	if err := database.AutoMigrate(&Notification{}); err != nil {
		return nil, err
	}

	w := &Notificator{
		bot: bot,
		database: database,
		config: config,
	}

	go w.run()

	return w, nil
}
