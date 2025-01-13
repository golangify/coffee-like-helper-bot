package commandhandler

import (
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/view/menu"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Menus(update *tgbotapi.Update, user *models.User, _ []string) {
	var menus []models.Menu
	err := h.database.Find(&menus).Error
	if err != nil {
		panic(err)
	}

	msg := tgbotapi.NewMessage(update.FromChat().ID, "Все меню")
	msg.ReplyMarkup = viewmenu.InlineKeyboardList(menus, "all_menus", 0, 5)

	h.bot.Send(msg)
}
