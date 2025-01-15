package viewbot

import (
	"coffee-like-helper-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartKeyboard(user *models.User) *tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Меню")),
		// tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Помощь")),
	)
	return &keyboard
}
