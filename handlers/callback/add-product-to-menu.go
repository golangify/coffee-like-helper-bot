package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewproduct "coffee-like-helper-bot/view/menu/product"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (h *CallbackHandler) addProductToMenu(update *tgbotapi.Update, user *models.User, args []string) {
	menuID, _ := strconv.ParseUint(args[1], 10, 32)
	var menu models.Menu
	err := h.database.First(&menu, menuID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Меню с id ", menuID, " не найдено.")))
			return
		}
		panic(err)
	}
	h.stepHandler.AddText(user, h.stepAddProductToMenu, []any{uint(menuID)})
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Отправь новое название напитка:"))
}

// args: models.Menu.ID
func (h *CallbackHandler) stepAddProductToMenu(update *tgbotapi.Update, user *models.User, args []any) {
	if update.Message == nil || update.Message.Text == "" || len(update.Message.Text) > 250 {
		panic("в сообщении отсутсвтует текст названия меню, или он слишком длинный")
	}

	menuID := args[0].(uint)
	var menu models.Menu
	err := h.database.First(&menu, menuID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Меню с id ", menuID, " не найдено.")))
			return
		}
		panic(err)
	}

	product := models.Product{
		MenuID: menu.ID,
		Name:   update.Message.Text,
	}
	err = h.database.Create(&product).Error
	if err != nil {
		panic(err)
	}

	h.bot.Send(viewproduct.Message(update.FromChat().ID, user, &product))
}
