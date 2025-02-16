package viewnotification

import (
	"coffee-like-helper-bot/models"
	workernotificator "coffee-like-helper-bot/workers/notificator"
	"encoding/json"
	"fmt"
	"html"
	"slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Message(chatID int64, user *models.User, notification *workernotificator.Notification) (tgbotapi.Chattable, error) {
	msgText, err := Text(notification)
	if err != nil {
		return nil, err
	}
	keyboard := InlineKeyboardEdit(notification)
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = keyboard
	return &msg, nil
}

func Text(notification *workernotificator.Notification) (string, error) {
	var weekdays []int
	var weekdaysStr string
	err := json.Unmarshal(notification.WeekDays, &weekdays)
	if err != nil {
		return "", err
	}
	for _, weekday := range weekdays {
		weekdaysStr += workernotificator.WeekdaysRu[weekday] + ", "
	}
	if len(weekdaysStr) > 3 {
		weekdaysStr = weekdaysStr[:len(weekdaysStr)-2]
	}
	msgText := fmt.Sprint(
		"<b>", html.EscapeString(notification.Name), "</b>\n\n",
		"Категория: ", notification.UserCategory, "\n",
		"Дни недели: ", weekdaysStr, "\n",
		"Время: ", notification.HourAndMinute, "\n",
		"Текст: ", notification.Text,
	)
	return msgText, nil
}

func InlineKeyboardEdit(notification *workernotificator.Notification) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⚙️редактировать уведомление", fmt.Sprint("edit_notification ", notification.ID)),
	))
	return &keyboard
}

func InlineKeyboardEditList(notification *workernotificator.Notification) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✏️изменить название", fmt.Sprint("edit_notification_name ", notification.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("👪изменить категрию пользователей", fmt.Sprint("edit_notification_user_category ", notification.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🗓️изменить дни недели", fmt.Sprint("edit_notification_weekdays ", notification.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⏲️изменить время", fmt.Sprint("edit_notification_time ", notification.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📝изменить текст", fmt.Sprint("edit_notification_text ", notification.ID))),
		// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("сделать беззвуч", fmt.Sprint("edit_notification_description ", notification.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🗑удалить", fmt.Sprint("delete_notification ", notification.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🏠готово", fmt.Sprint("notification_home ", notification.ID))),
	)
	return &keyboard
}

func InlineKeyboardSetUserCategory(notification *workernotificator.Notification) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	if notification.UserCategory != "all" {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("all (вообще все пользователи)", fmt.Sprint("set_notification_user_category ", notification.ID, " all"))))
	}
	if notification.UserCategory != "admin" {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("admin (только администраторы)", fmt.Sprint("set_notification_user_category ", notification.ID, " admin"))))
	}
	if notification.UserCategory != "barista" {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("barista (только бариста)", fmt.Sprint("set_notification_user_category ", notification.ID, " barista"))))
	}
	return &keyboard
}

func InlineKeyboardNotificationWeekdays(notification *workernotificator.Notification) (*tgbotapi.InlineKeyboardMarkup, error) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	var weekdays []int
	if err := json.Unmarshal(notification.WeekDays, &weekdays); err != nil {
		return nil, err
	}
	for i := 1; i < 8; i++ {
		if slices.Contains(weekdays, i) {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✅"+workernotificator.WeekdaysRu[i], fmt.Sprint("remove_weekday_from_notification ", notification.ID, " ", i))),
			)
		} else {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("➕"+workernotificator.WeekdaysRu[i], fmt.Sprint("add_weekday_to_notification ", notification.ID, " ", i))),
			)
		}
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🏠готово", fmt.Sprint("notification_home ", notification.ID)),
	))
	return &keyboard, nil
}
