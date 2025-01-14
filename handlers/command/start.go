package commandhandler

import (
	"coffee-like-helper-bot/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Start(update *tgbotapi.Update, user *models.User, _ []string) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint(
		"Привет, <b>", user.FirstName, "</b>!\n\n",
		"В этом боте можно просматривать ТТК напитков.\n\n",
		"/menus - список всех меню\n",
		"/help - все команды бота(<b>обязательно посмотри</b>)",
	))
	msg.ParseMode = tgbotapi.ModeHTML
	h.bot.Send(msg)
}
