package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewmenu "coffee-like-helper-bot/view/menu"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (h *CallbackHandler) editMenuName(update *tgbotapi.Update, user *models.User, args []string) {
	menuID, _ := strconv.ParseUint(args[1], 10, 32)
	var menu models.Menu
	err := h.database.First(&menu, menuID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Меню с id ", menu.ID, " не найдено.")))
			return
		}
		panic(err)
	}
	h.stepHandler.AddText(user, h.StepUpdateMenuName, []any{uint(menuID)})
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Отправь новое название меню:"))
}

// args: models.Menu.ID
func (h *CallbackHandler) StepUpdateMenuName(update *tgbotapi.Update, user *models.User, args []any) {
	menuID := args[0].(uint)
	var menu models.Menu
	err := h.database.First(&menu, menuID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Меню с id ", menu.ID, " не найдено.")))
			return
		}
		panic(err)
	}
	menu.Name = update.Message.Text
	err = h.database.Model(&menu).UpdateColumn("name", menu.Name).Error
	if err != nil {
		panic(err)
	}
	h.bot.Send(viewmenu.Message(update.FromChat().ID, user, &menu, nil, 0, 0))
}
