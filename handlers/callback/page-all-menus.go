package callbackhandler

import (
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/view/menu"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

// args: _, page int
func (h *CallbackHandler) pageAllMenus(update *tgbotapi.Update, user *models.User, args []string) {
	page, _ := strconv.Atoi(args[1])

	var menus []models.Menu
	err := h.database.Find(&menus).Error
	if err != nil {
		panic(err)
	}

	msg := tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, *viewmenu.InlineKeyboardList(menus, "all_menus", page, 5))

	h.bot.Send(msg)
}
