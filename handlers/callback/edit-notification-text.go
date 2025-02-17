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
func (h *CallbackHandler) editNotificationText(update *tgbotapi.Update, user *models.User, args []string) {
	notificationID, _ := strconv.ParseUint(args[1], 10, 32)
	var notification workernotificator.Notification
	if err := h.database.First(&notification, notificationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Уведомление с id ", notificationID, " не найдено.")))
			return
		}
		panic(err)
	}
	h.stepHandler.AddStepHandler(user.ID, h.stepEditNotificationText, []any{update, &notification})
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Отправь новый текст для уведомления:\n\n/cancel - отмена"))
}

// args: tgbotapi.Update, *workernotificator.Notification
func (h *CallbackHandler) stepEditNotificationText(update *tgbotapi.Update, user *models.User, args []any) {
	if update.Message == nil || update.Message.Text == "" || len(update.Message.Text) > 250 {
		panic("в сообщении отсутсвтует текст для уведомления, или он слишком длинный")
	}
	sourceUpdate := args[0].(*tgbotapi.Update)
	notification := args[1].(*workernotificator.Notification)
	notification.Text = update.Message.Text
	if err := h.database.Model(&notification).Update("text", notification.Text).Error; err != nil {
		panic(err)
	}

	text, err := viewnotification.Text(notification)
	if err != nil {
		panic(err)
	}
	keyboard := viewnotification.InlineKeyboardEdit(notification)
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, sourceUpdate.CallbackQuery.Message.MessageID, text, *keyboard)
	msg.ParseMode = tgbotapi.ModeHTML
	h.bot.Send(msg)
}
