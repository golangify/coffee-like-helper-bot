package commandhandler

import (
	"coffee-like-helper-bot/models"
	viewbot "coffee-like-helper-bot/view/bot"
	"fmt"
	"html"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Help(update *tgbotapi.Update, user *models.User, _ []string) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Сводка команд\n\n")
	for _, command := range h.commands {
		if command.isForStaff && !user.IsAdministrator || command.isHidden {
			continue
		}
		msg.Text += fmt.Sprint(html.EscapeString(command.help()), "\n\n")
	}
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = viewbot.StartKeyboard(user)
	h.bot.Send(msg)
}
