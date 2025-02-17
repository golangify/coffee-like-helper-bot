package callbackhandler

import (
	"coffee-like-helper-bot/models"
	workernotificator "coffee-like-helper-bot/workers/notificator"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

// args: _, workernotificator.Notification.ID
func (h *CallbackHandler) deleteNotification(update *tgbotapi.Update, user *models.User, args []string) {
	notificationID, _ := strconv.ParseUint(args[1], 10, 32)
	var notification workernotificator.Notification
	if err := h.database.First(&notification, notificationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Уведомление с id ", notificationID, " не найдено.")))
			return
		}
		panic(err)
	}

	if err := h.database.Delete(&notification).Error; err != nil {
		panic(err)
	}
	h.bot.Send(tgbotapi.NewDeleteMessage(update.FromChat().ID, update.CallbackQuery.Message.MessageID))
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Уведомление успешно удалено."))
}
