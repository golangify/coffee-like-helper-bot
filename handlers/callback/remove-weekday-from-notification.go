package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewnotification "coffee-like-helper-bot/view/notification"
	workernotificator "coffee-like-helper-bot/workers/notificator"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

// args: _, workernotificator.Notification.ID, workernotificator.Weekdays
func (h *CallbackHandler) removeWeekdayFromNotification(update *tgbotapi.Update, _ *models.User, args []string) {
	var err error
	notificationID, _ := strconv.ParseUint(args[1], 10, 32)
	weekday, _ := strconv.Atoi(args[2])
	var notification workernotificator.Notification
	if err := h.database.First(&notification, notificationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Уведомление с id ", notificationID, " не найдено.")))
			return
		}
		panic(err)
	}

	var weekdays []int
	if err := json.Unmarshal(notification.WeekDays, &weekdays); err != nil {
		panic(err)
	}

	if !slices.Contains(weekdays, weekday) {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint(`В уведомлении "`, notification.Name, `" уже нет `, workernotificator.WeekdaysRu[weekday], ". (где склонение? — его муторно делать...)")))
		return
	}
	weekdays = slices.DeleteFunc(weekdays, func(wd int) bool { return wd == weekday })
	if notification.WeekDays, err = json.Marshal(&weekdays); err != nil {
		panic(err)
	}
	if err = h.database.Model(&notification).Update("week_days", notification.WeekDays).Error; err != nil {
		panic(err)
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, "", tgbotapi.NewInlineKeyboardMarkup())
	msg.Text, err = viewnotification.Text(&notification)
	if err != nil {
		panic(err)
	}
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup, err = viewnotification.InlineKeyboardNotificationWeekdays(&notification)
	if err != nil {
		panic(err)
	}
	h.bot.Send(msg)
}
