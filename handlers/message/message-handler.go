package messagehandler

import (
	"coffee-like-helper-bot/config"
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/service/search"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type MessageHandler struct {
	bot      *tgbotapi.BotAPI
	database *gorm.DB

	searchEngine *search.SearchEngine
}

func NewMessageHandler(config *config.Config, bot *tgbotapi.BotAPI, database *gorm.DB) *MessageHandler {
	h := &MessageHandler{
		bot:      bot,
		database: database,

		searchEngine: search.NewEngine(config, database),
	}
	return h
}

func (h *MessageHandler) Process(update *tgbotapi.Update, user *models.User) {
	if update.Message.Text == "Меню" {
		h.Menus(update, user)
	} else {
		h.Search(update, user)
	}
}
