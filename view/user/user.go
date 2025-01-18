package viewuser

import (
	"coffee-like-helper-bot/models"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Message(chatID int64, user *models.User) tgbotapi.MessageConfig {

	msgText := fmt.Sprint(
		"ID: ", user.ID,
		"\n\n Первое сообщение: ", user.CreatedAt,
		"\n\n Обновление: ", user.UpdatedAt,
		"\n\n ID telegram: ", user.TelegramID,
		"\n\n Бариста: ", user.IsBarista,
		"\n\n Администратор: ", user.IsAdministrator,
		"\n\n Имя: ", user.FirstName,
		"\n\n Фамилия: ", user.LastName,
		"\n\n Имя пользователя: @", user.UserName,
	)

	msg := tgbotapi.NewMessage(chatID, msgText)
	return msg
}

func Text(user *models.User) string {
	text := /* fmt.Sprintf("[%d;%d] ", user.ID, user.TelegramID) + */ user.FirstName
	if user.LastName != "" {
		text += " " + user.LastName
	}
	if user.UserName != "" {
		text += fmt.Sprintf("(@%s)", user.UserName)
	}

	return text
}

func InlineKeyboardEdit(user *models.User) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("действия", fmt.Sprint("edit_user ", user.ID)),
	))
	return &keyboard
}

func InlineKeyboardEditList(user *models.User) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	if user.IsBarista {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("удалить из бариста", fmt.Sprint("remove_barista ", user.ID))),
		)
	} else {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("сделать бариста", fmt.Sprint("make_barista ", user.ID))),
		)
	}
	if user.IsAdministrator {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("удалить из администраторов", fmt.Sprint("remove_administrator ", user.ID))),
		)
	} else {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("сделать администратором", fmt.Sprint("make_administrator ", user.ID))),
		)
	}

	return &keyboard
}

func InlineKeyboardList(users []models.User, callbackData string, page, limit int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()

	if len(users) == 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("пусто", "nothing")))
		return &keyboard
	}

	totalPages := (len(users) + limit - 1) / limit
	startIndex := page * limit
	endIndex := (page + 1) * limit
	if endIndex > len(users) {
		endIndex = len(users)
	}

	users = users[startIndex:endIndex]

	for _, user := range users {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(Text(&user), fmt.Sprintf("user %d", user.ID))),
		)
	}

	if totalPages > 1 {
		navigationRow := tgbotapi.NewInlineKeyboardRow()
		if page > 0 {
			navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData("<<", fmt.Sprintf("page %s %d", callbackData, page-1)))
		}
		navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("(%d / %d)", page+1, totalPages), "nothing"))
		if page < totalPages-1 {
			navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData(">>", fmt.Sprintf("page %s %d", callbackData, page+1)))
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, navigationRow)
	}

	return &keyboard
}
