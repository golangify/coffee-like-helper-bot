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
	if err := h.database.First(&notification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Уведомление с id ", notificationID, " не найдено.")))
			return
		}
		panic(err)
	}

	h.stepHandler.AddStepHandler(user.ID, h.stepEditNotificationTime, []any{&notification, update.CallbackQuery.Message.MessageID})
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Отправь время. Пример: <code>20:30</code> - будет приходить в пол девятого вечера.\n\n/cancel - отмена")
	msg.ParseMode = tgbotapi.ModeHTML
	h.bot.Send(msg)
}

// args: *workernotificator.Notification, tgbotapi.Message.MessageID
func (h *CallbackHandler) stepEditNotificationTime(update *tgbotapi.Update, user *models.User, args []any) {
	notification := args[0].(*workernotificator.Notification)
	notificationMessageID := args[1].(int)
	if update.Message == nil || update.Message.Text == "" || len(update.Message.Text) > 250 {
		panic("в сообщении отсутсвтует текст, или он слишком длинный")
	}

	notification.HourAndMinute = update.Message.Text
	if _, err := notification.TimeUntilNextNotification(); err != nil {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Ошибка подсчёта нового времени отправки - ", err.Error(), "\nИзменения не сохранены.")))
		return
	}
	if err := h.database.Model(&notification).Update("hour_and_minute", notification.HourAndMinute).Error; err != nil {
		panic(err)
	}

	msgText, err := viewnotification.Text(notification)
	if err != nil {
		panic(err)
	}
	keyboard := viewnotification.InlineKeyboardEdit(notification)
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, notificationMessageID, msgText, *keyboard)
	msg.ParseMode = tgbotapi.ModeHTML
	h.bot.Send(msg)
}
