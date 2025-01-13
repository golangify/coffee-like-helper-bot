package commandhandler

import (
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/view/menu"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strconv"
)

func (h *CommandHandler) Menu(update *tgbotapi.Update, user *models.User, args []string) {
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

	var products []models.Product
	err = h.database.Find(&products, "menu_id = ?", menu.ID).Error
	if err != nil {
		panic(err)
	}

	_, err = h.bot.Send(viewmenu.Message(update.FromChat().ID, user, &menu, products, 0, 5))
	if err != nil {
		panic(err)
	}
}
