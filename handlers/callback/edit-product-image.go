package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewproduct "coffee-like-helper-bot/view/menu/product"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (h *CallbackHandler) editProductImage(update *tgbotapi.Update, user *models.User, args []string) {
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
	h.stepHandler.AddImage(user, h.stepUpdateProductImage, []any{product.ID})
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID,
		"Чтобы добавить изображение - просто отправь изображение.\n\n"+
			"Чтобы удалить изображение - /delete",
	))
}

// args: models.Product.ID
func (h *CallbackHandler) stepUpdateProductImage(update *tgbotapi.Update, user *models.User, args []any) {
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

	if update.Message != nil {
		if update.Message.Text == "/delete" {
			product.ImageFileID = nil
			err = h.database.Model(&product).UpdateColumn("image_file_id", product.ImageFileID).Error
			if err != nil {
				panic(err)
			}
			_, err = h.bot.Send(viewproduct.Message(update.FromChat().ID, user, &product))
			if err != nil {
				panic(err)
			}
			return
		}
	}

	product.ImageFileID = &update.Message.Photo[0].FileID
	err = h.database.Model(&product).UpdateColumn("image_file_id", product.ImageFileID).Error
	if err != nil {
		panic(err)
	}
	_, err = h.bot.Send(viewproduct.Message(update.FromChat().ID, user, &product))
	if err != nil {
		panic(err)
	}
}
