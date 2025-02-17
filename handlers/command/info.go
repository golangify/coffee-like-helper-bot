package commandhandler

import (
	"coffee-like-helper-bot/models"
	workernotificator "coffee-like-helper-bot/workers/notificator"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) info(update *tgbotapi.Update, user *models.User, _ []string) {
	var (
		totalUsersCount,
		totalProductsCount,
		totalMenusCount,
		totalQueriesCount,
		totalNotificationsCount int64

		numActiveNotifications int
	)
	if err := h.database.Model(&models.User{}).Count(&totalUsersCount).Error; err != nil {
		panic(err)
	}
	if err := h.database.Model(&models.Product{}).Count(&totalProductsCount).Error; err != nil {
		panic(err)
	}
	if err := h.database.Model(&models.Menu{}).Count(&totalMenusCount).Error; err != nil {
		panic(err)
	}
	if err := h.database.Model(&models.Search{}).Count(&totalQueriesCount).Error; err != nil {
		panic(err)
	}
	if err := h.database.Model(&workernotificator.Notification{}).Count(&totalNotificationsCount).Error; err != nil {
		panic(err)
	}
	numActiveNotifications = h.notificator.NumRunning()

	msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint(
		"Всего пользователей: ", totalUsersCount, "\n",
		"Всего напитков: ", totalProductsCount, "\n",
		"Всего меню: ", totalMenusCount, "\n",
		"Всего поисковых запросов: ", totalQueriesCount, "\n",
		"Всего уведомлений: ", totalNotificationsCount, "\n",
		"Уведомлений сейчас запущено: ", numActiveNotifications, "\n",
	))
	h.bot.Send(msg)
}
