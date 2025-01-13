package commandhandler

import (
	"coffee-like-helper-bot/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Help(update *tgbotapi.Update, user *models.User, _ []string) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Сводка команд\n\n")
	for _, command := range h.commands {
		if command.isForStaff && !user.IsAdministrator {
			continue
		}
		msg.Text += fmt.Sprint(command.help(), "\n\n")
	}
	h.bot.Send(msg)
}
