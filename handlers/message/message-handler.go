package messagehandler

import (
	"coffee-like-helper-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type MessageHandler struct {
	bot      *tgbotapi.BotAPI
	database *gorm.DB
}

func NewMessageHandler(bot *tgbotapi.BotAPI, database *gorm.DB) *MessageHandler {
	h := &MessageHandler{
		bot:      bot,
		database: database,
	}

	return h
}

func (h *MessageHandler) Process(update *tgbotapi.Update, user *models.User) {
	h.Search(update, user)
}
