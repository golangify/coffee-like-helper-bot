package commandhandler

import (
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/view/bot"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Start(update *tgbotapi.Update, user *models.User, _ []string) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint(
		"Привет, <b>", user.FirstName, "</b>!\n\n",
		"В этом боте можно просматривать ТТК напитков.\n\n",
		"/help - возможности бота",
	))
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = viewbot.StartKeyboard(user)
	h.bot.Send(msg)
}
