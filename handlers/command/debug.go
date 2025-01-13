package commandhandler

import (
	"coffee-like-helper-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Debug(update *tgbotapi.Update, user *models.User, _ []string) {
	h.bot.Debug = !h.bot.Debug
	if h.bot.Debug {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Debug активирован"))
	} else {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Debug деактивирован"))
	}
}
