package commandhandler

import (
	"coffee-like-helper-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Panic(update *tgbotapi.Update, user *models.User, args []string) {
	panic(args[1])
}
