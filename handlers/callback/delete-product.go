package callbackhandler

import (
	"coffee-like-helper-bot/models"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (h *CallbackHandler) deleteProduct(update *tgbotapi.Update, user *models.User, args []string) {
	productID, _ := strconv.ParseUint(args[1], 10, 32)
	var product models.Product
	err := h.database.First(&product, productID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Напиток с id ", productID, " не найден.")))
			return
		}
		panic(err)
	}

	checkPhrase := fmt.Sprint("удалить \"", product.Name, "\"")
	h.stepHandler.AddStepHandler(user.ID, h.stepDeleteProduct, []any{product.ID, checkPhrase})

	msg := tgbotapi.NewMessage(update.FromChat().ID,
		"Отправь(можно нажать на текст ниже для копирования)\n\n<code>"+checkPhrase+"</code>\n\nдля потдверждения удаления напитка.\n\n/cancel для отмены",
	)
	msg.ParseMode = tgbotapi.ModeHTML

	h.bot.Send(msg)
}

// args: models.Product.ID, checkPhrase string
func (h *CallbackHandler) stepDeleteProduct(update *tgbotapi.Update, user *models.User, args []any) {
	productID := args[0].(uint)
	checkPhrase := args[1].(string)

	var product models.Product
	err := h.database.First(&product, productID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Напиток с id ", productID, " не найден.")))
			return
		}
		panic(err)
	}

	if update.Message == nil || update.Message.Text == "" {
		panic("это не сообщение")
	}

	if update.Message.Text != checkPhrase {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Неправильная подтверждающая фраза. Удаление отменено."))
		return
	}

	err = h.database.Delete(&product).Error
	if err != nil {
		panic(err)
	}

	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Напиток был успешно удален."))
}
