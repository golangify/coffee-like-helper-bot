package commandhandler

import (
	"coffee-like-helper-bot/models"
	viewnotification "coffee-like-helper-bot/view/notification"
	workernotificator "coffee-like-helper-bot/workers/notificator"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) addNotification(update *tgbotapi.Update, user *models.User, _ []string) {
	h.stepHandler.AddStepHandler(user.ID, h.stepAddNotificationName, nil)
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Отправь название для уведомления\n\n/cnacel - отмена"))
}

func (h *CommandHandler) stepAddNotificationName(update *tgbotapi.Update, user *models.User, _ []any) {
	if update.Message == nil || update.Message.Text == "" || len([]rune(update.Message.Text)) > maxMenuNameLength {
		panic("в сообщении отсутствует текст или он слишком длинный")
	}
	notification := workernotificator.Notification{
		Name:     update.Message.Text,
		WeekDays: []byte("[]"),
	}
	err := h.notificator.AddNotification(&notification)
	if err != nil {
	}
	msg, err := viewnotification.Message(update.FromChat().ID, user, &notification)
	if err != nil {
		panic(err)
	}
	h.bot.Send(msg)

}

// /addnotification "<name>" <user category> [<weekdays>] <text>
// func (h *CommandHandler) addNotification(update *tgbotapi.Update, user *models.User, args []string) {
// 	name := args[1]
// 	userCategory := args[2]
// 	weekdays := args[3]
// 	hoursAndMinute := args[4]
// 	text := args[5]

// 	notification := workernotificator.Notification{
// 		Name:          name,
// 		UserCategory:  userCategory,
// 		WeekDays:      []byte(weekdays),
// 		HourAndMinute: hoursAndMinute,
// 		Text:          text,
// 	}

// 	timeUntil, err := notification.TimeUntilNextNotification()
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = h.notificator.AddNotification(notification)
// 	if err != nil {
// 		panic(err)
// 	}

// 	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Уведомление успешно добавлено. Будет разослано через %.1fч",
// 		timeUntil.Hours(),
// 	)))
// }
