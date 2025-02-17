package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewnotification "coffee-like-helper-bot/view/notification"
	workernotificator "coffee-like-helper-bot/workers/notificator"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

// args: _, workernotificator.Notification.ID
func (h *CallbackHandler) setNotificationUserCategory(update *tgbotapi.Update, user *models.User, args []string) {
	notificationID, _ := strconv.ParseUint(args[1], 10, 32)
	userCategory := args[2]
	var notification workernotificator.Notification
	if err := h.database.First(&notification, notificationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Уведомление с id ", notificationID, " не найдено.")))
			return
		}
		panic(err)
	}
	notification.UserCategory = userCategory
	if err := h.database.Model(notification).Update("user_category", notification.UserCategory).Error; err != nil {
		panic(err)
	}

	text, err := viewnotification.Text(&notification)
	if err != nil {
		panic(err)
	}
	keyboard := viewnotification.InlineKeyboardEdit(&notification)
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, text, *keyboard)
	msg.ParseMode = tgbotapi.ModeHTML
	h.bot.Send(msg)
}
