package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewuser "coffee-like-helper-bot/view/user"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CallbackHandler) pageUsers(update *tgbotapi.Update, user *models.User, args []string) {
	userCategory := args[1]
	page, err := strconv.Atoi(args[2])
	if err != nil {
		panic(err)
	}
	var users []models.User
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

	msg := tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, *viewuser.InlineKeyboardList(users, "users_"+userCategory, page, 5))
	h.bot.Send(msg)
}
