package commandhandler

import (
	"coffee-like-helper-bot/models"
	viewuser "coffee-like-helper-bot/view/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Users(update *tgbotapi.Update, user *models.User, args []string) {
	userCategory := args[1]
	var users []models.User
	var err error
	switch userCategory {
	case "admin":
		err = h.database.Find(&users, "is_administrator = ?", true).Error
	case "barista":
		err = h.database.Find(&users, "is_barista = ?", true).Error
	default:
		userCategory = "all"
		err = h.database.Find(&users).Error
	}
	if err != nil {
		panic(err)
	}

	msg := tgbotapi.NewMessage(update.FromChat().ID, "Пользователи категории "+userCategory)
	msg.ReplyMarkup = viewuser.InlineKeyboardList(users, "users_"+userCategory, 0, 5)
	h.bot.Send(msg)
}
