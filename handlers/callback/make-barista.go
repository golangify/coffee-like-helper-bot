package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewuser "coffee-like-helper-bot/view/user"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CallbackHandler) makeBarista(update *tgbotapi.Update, user *models.User, args []string) {
	targetUserID, _ := strconv.ParseUint(args[1], 10, 32)
	targetUser, err := h.userService.UserByID(uint(targetUserID))
	if err != nil {
		panic(err)
	}
	if err = h.userService.MakeBarista(targetUser); err != nil {
		panic(err)
	}

	h.bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
		*viewuser.InlineKeyboardEditList(targetUser),
	))

	notifMsg := tgbotapi.NewMessage(0, fmt.Sprint(
		viewuser.Text(user), " сделал(а) бариста ", viewuser.Text(targetUser),
	))
	notifMsg.ReplyMarkup = viewuser.InlineKeyboardEdit(targetUser)
	h.mailer.Administrator(&notifMsg)

	h.bot.Send(tgbotapi.NewMessage(targetUser.TelegramID, "Тебя назначили на роль бариста в боте!\n\nНажимай -> /help"))
}
