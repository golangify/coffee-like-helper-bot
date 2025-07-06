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
func (h *CallbackHandler) editNotificationTime(update *tgbotapi.Update, user *models.User, args []string) {
	notificationID, _ := strconv.ParseUint(args[1], 10, 32)
	var notification workernotificator.Notification
	if err := h.database.First(&notification, notificationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Уведомление с id ", notificationID, " не найдено.")))
			return
		}
		panic(err)
	}

	h.stepHandler.AddText(user, h.stepEditNotificationTime, []any{&notification, update.CallbackQuery.Message.MessageID})
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Отправь время. Пример: <code>20:30</code> - будет приходить в пол девятого вечера:")
	msg.ParseMode = tgbotapi.ModeHTML
	h.bot.Send(msg)
}

// args: *workernotificator.Notification, tgbotapi.Message.MessageID
func (h *CallbackHandler) stepEditNotificationTime(update *tgbotapi.Update, user *models.User, args []any) {
	notification := args[0].(*workernotificator.Notification)
	notificationMessageID := args[1].(int)

	notification.HourAndMinute = update.Message.Text
	if _, err := notification.TimeUntilNextNotification(h.notificator.TimeZone); err != nil {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Ошибка подсчёта нового времени отправки - ", err.Error(), "\nИзменения не сохранены.")))
		return
	}
	if err := h.database.Model(&notification).Update("hour_and_minute", notification.HourAndMinute).Error; err != nil {
		panic(err)
	}

	go h.notificator.NotificationProcess(notification, false)

	msgText, err := viewnotification.Text(notification)
	if err != nil {
		panic(err)
	}
	keyboard := viewnotification.InlineKeyboardEdit(notification)
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, notificationMessageID, msgText, *keyboard)
	msg.ParseMode = tgbotapi.ModeHTML
	h.bot.Send(msg)
}
