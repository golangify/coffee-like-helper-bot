package callbackhandler

import (
	"coffee-like-helper-bot/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"strconv"
)

func (h *CallbackHandler) deleteMenu(update *tgbotapi.Update, user *models.User, args []string) {
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

	checkPhrase := fmt.Sprint("удалить \"", menu.Name, "\" вместе с напитками")
	h.stepHandler.AddStepHandler(user.ID, h.stepDeleteMenu, []any{menu.ID, checkPhrase})

	msg := tgbotapi.NewMessage(update.FromChat().ID, "<b>ВНИМАНИЕ</b>\nЭто действие удалит меню и все напитки в нём!\n\nОтправь(можно нажать на текст ниже для копирования)\n\n<code>"+checkPhrase+"</code>\n\nдля потдверждения.\n\n/cancel для отмены")
	msg.ParseMode = tgbotapi.ModeHTML

	h.bot.Send(msg)
}

// args: models.Menu.ID, checkPhrase string
func (h *CallbackHandler) stepDeleteMenu(update *tgbotapi.Update, user *models.User, args []any) {
	menuID := args[0].(uint)
	checkPhrase := args[1].(string)
	var menu models.Menu
	err := h.database.First(&menu, menuID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Меню с id ", menuID, " не найдено.")))
			return
		}
		panic(err)
	}

	if update.Message == nil || update.Message.Text == "" {
		panic("это не сообщение")
		return
	}

	if update.Message.Text == "/cancel" {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Удаление успешно отменено."))
		return
	}

	if update.Message.Text != checkPhrase {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Неправильная подтверждающая фраза. Удаление отменено."))
		return
	}

	err = h.database.Where("menu_id = ?", menu.ID).Delete(&models.Product{}).Error
	if err != nil {
		panic(err)
	}

	err = h.database.Delete(&menu).Error
	if err != nil {
		panic(err)
	}

	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Меню было успешно удалено."))
}
