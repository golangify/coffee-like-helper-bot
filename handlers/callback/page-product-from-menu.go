package callbackhandler

import (
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/view/menu/product"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func (h *CallbackHandler) pageProductFromMenu(update *tgbotapi.Update, user *models.User, args []string) {
	menuID, _ := strconv.ParseUint(args[1], 10, 32)
	page, _ := strconv.Atoi(args[2])

	var products []models.Product
	err := h.database.Find(&products, "menu_id = ?", menuID).Error
	if err != nil {
		panic(err)
	}

	_, err = h.bot.Send(tgbotapi.NewEditMessageReplyMarkup(
		update.FromChat().ID,
		update.CallbackQuery.Message.MessageID,
		*viewproduct.InlineKeyboardList(products, fmt.Sprint("product_from_menu ", menuID), page, 5),
	))
	if err != nil {
		panic(err)
	}
}
