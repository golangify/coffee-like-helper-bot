package callbackhandler

import (
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/view/user"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strconv"
)

func (h *CallbackHandler) editUser(update *tgbotapi.Update, user *models.User, args []string) {
	userID, _ := strconv.ParseUint(args[1], 10, 32)
	var euser models.User
	err := h.database.First(&euser, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Пользователь с id ", userID, " не найден.")))
			return
		}
		panic(err)
	}

	keyboard := viewuser.InlineKeyboardEditList(&euser)
	msg := tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, *keyboard)
	h.bot.Send(msg)
}
