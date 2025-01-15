package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewuser "coffee-like-helper-bot/view/user"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (h *CallbackHandler) user(update *tgbotapi.Update, _ *models.User, args []string) {
	userID, _ := strconv.ParseUint(args[1], 10, 32)

	var targetUser models.User
	err := h.database.First(&targetUser, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Пользователь с id ", userID, " не найден.")))
			return
		}
		panic(err)
	}

	msg := viewuser.Message(update.FromChat().ID, &targetUser)
	msg.ReplyMarkup = viewuser.InlineKeyboardEdit(&targetUser)
	h.bot.Send(msg)
}
