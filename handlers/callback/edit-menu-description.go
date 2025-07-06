package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewmenu "coffee-like-helper-bot/view/menu"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (h *CallbackHandler) editMenuDescription(update *tgbotapi.Update, user *models.User, args []string) {
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
	h.stepHandler.AddText(user, h.StepUpdateMenuDescription, []any{uint(menuID)})
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Отправь новое описание меню:"))
}

// args: models.Menu.ID
func (h *CallbackHandler) StepUpdateMenuDescription(update *tgbotapi.Update, user *models.User, args []any) {
	if update.Message == nil || update.Message.Text == "" || len(update.Message.Text) > 1000 {
		panic("в сообщении отсутсвтует текст описания меню, или он слишком длинный")
	}
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
	menu.Description = update.Message.Text
	err = h.database.Model(&menu).UpdateColumn("description", menu.Description).Error
	if err != nil {
		panic(err)
	}
	h.bot.Send(viewmenu.Message(update.FromChat().ID, user, &menu, nil, 0, 0))
}
