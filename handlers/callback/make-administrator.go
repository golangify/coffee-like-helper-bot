package callbackhandler

import (
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/view/user"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strconv"
)

func (h *CallbackHandler) makeAdministrator(update *tgbotapi.Update, user *models.User, args []string) {
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

	euser.IsAdministrator = true
	err = h.database.Model(&euser).UpdateColumn("is_administrator", true).Error
	if err != nil {
		panic(err)
	}

	h.bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, 
		*viewuser.InlineKeyboardEditList(&euser),
	))

	notifMsg := tgbotapi.NewMessage(0, fmt.Sprint(
		viewuser.Text(user), " сделал(а) администратором бота ", viewuser.Text(&euser),
	))
	notifMsg.ReplyMarkup = viewuser.InlineKeyboardEdit(&euser)
	h.mailer.Administrator(&notifMsg)
	
	h.bot.Send(tgbotapi.NewMessage(euser.TelegramID, "Тебя назначили на роль администратора в боте!\n\nНажимай -> /help"))
}
