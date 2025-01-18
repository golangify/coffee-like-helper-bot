package commandhandler

import (
	"coffee-like-helper-bot/models"
	workernotificator "coffee-like-helper-bot/workers/notificator"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// /addnotification "<name>" <user category> [<weekdays>] <text>
func (h *CommandHandler) addNotification(update *tgbotapi.Update, user *models.User, args []string) {
	name := args[1]
	userCategory := args[2]
	weekdays := args[3]
	hoursAndMinute := args[4]
	text := args[5]

	notification := workernotificator.Notification{
		Name:          name,
		UserCategory:  userCategory,
		WeekDays:      []byte(weekdays),
		HourAndMinute: hoursAndMinute,
		Text:          text,
	}

	timeUntil, err := notification.TimeUntilNextNotification()
	if err != nil {
		panic(err)
	}

	err = h.notificator.AddNotification(notification)
	if err != nil {
		panic(err)
	}

	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Уведомление успешно добавлено. Будет разослано через %.1fч",
		timeUntil.Hours(),
	)))
}
