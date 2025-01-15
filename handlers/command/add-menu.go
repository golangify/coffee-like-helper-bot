package commandhandler

import (
	"coffee-like-helper-bot/models"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maxMenuNameLength = 250

func (h *CommandHandler) AddMenu(update *tgbotapi.Update, user *models.User, _ []string) {
	h.stepHandler.AddStepHandler(user.ID, h.StepAddMenuName, nil)
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Отправь название меню:"))
}

func (h *CommandHandler) StepAddMenuName(update *tgbotapi.Update, user *models.User, _ []any) {
	if update.Message == nil || update.Message.Text == "" || len([]rune(update.Message.Text)) > maxMenuNameLength {
		panic("в сообщении отсутствует текст или он слишком длинный")
	}

	newMenu := models.Menu{
		Name: update.Message.Text,
	}

	err := h.database.Create(&newMenu).Error
	if err != nil {
		panic(err)
	}

	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Меню успешно добавлено - /menu_", newMenu.ID)))
}
