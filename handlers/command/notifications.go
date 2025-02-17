package commandhandler

import (
	"coffee-like-helper-bot/models"
	viewnotification "coffee-like-helper-bot/view/notification"
	workernotificator "coffee-like-helper-bot/workers/notificator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) notifications(update *tgbotapi.Update, _ *models.User, _ []string) {
	var notifications []workernotificator.Notification
	err := h.database.Find(&notifications).Error
	if err != nil {
		panic(err)
	}

	msg := tgbotapi.NewMessage(update.FromChat().ID, "Уведомления")
	msg.ReplyMarkup = viewnotification.InlineKeyboardList(notifications, "all_notifications", 0, 100)

	h.bot.Send(msg)
}
