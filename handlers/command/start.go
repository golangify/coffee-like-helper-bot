package commandhandler

import (
	"coffee-like-helper-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Start(update *tgbotapi.Update, user *models.User, _ []string) {
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "хай"))
}
