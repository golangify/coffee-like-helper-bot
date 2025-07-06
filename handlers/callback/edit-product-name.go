package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewproduct "coffee-like-helper-bot/view/menu/product"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (h *CallbackHandler) editProductName(update *tgbotapi.Update, user *models.User, args []string) {
	productID, _ := strconv.ParseUint(args[1], 10, 32)
	var product models.Product
	err := h.database.First(&product, productID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Напиток с id ", product.ID, " не найден.")))
			return
		}
		panic(err)
	}
	h.stepHandler.AddText(user, h.stepUpdateProductName, []any{uint(productID)})
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Отправь новое название напитка:"))
}

// args: models.Product.ID
func (h *CallbackHandler) stepUpdateProductName(update *tgbotapi.Update, user *models.User, args []any) {
	productID := args[0].(uint)
	var product models.Product
	err := h.database.First(&product, productID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Напиток с id ", product.ID, " не найден.")))
			return
		}
		panic(err)
	}
	product.Name = update.Message.Text
	err = h.database.Model(&product).UpdateColumn("name", product.Name).Error
	if err != nil {
		panic(err)
	}
	h.bot.Send(viewproduct.Message(update.FromChat().ID, user, &product))
}
