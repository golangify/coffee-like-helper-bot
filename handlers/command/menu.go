package commandhandler

import (
	"coffee-like-helper-bot/models"
	viewmenu "coffee-like-helper-bot/view/menu"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) Menu(update *tgbotapi.Update, user *models.User, args []string) {
	menuID, _ := strconv.ParseUint(args[1], 10, 32)

	menu, err := h.menuService.MenuByID(uint(menuID))
	if err != nil {
		panic(err)
	}

	var products []models.Product
	err = h.database.Find(&products, "menu_id = ?", menu.ID).Error
	if err != nil {
		panic(err)
	}

	_, err = h.bot.Send(viewmenu.Message(update.FromChat().ID, user, menu, products, 0, 5))
	if err != nil {
		panic(err)
	}
}
