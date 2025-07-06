package callbackhandler

import (
	"coffee-like-helper-bot/models"
	viewmenu "coffee-like-helper-bot/view/menu"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (h *CallbackHandler) editMenuImage(update *tgbotapi.Update, user *models.User, args []string) {
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
	h.stepHandler.AddImage(user, h.StepUpdateMenuImage, []any{uint(menuID)})
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID,
		"Чтобы добавить изображение - просто отправь изображение.\n\n"+
			"Чтобы удалить изображение - /delete",
	))
}

// args: models.Menu.ID
func (h *CallbackHandler) StepUpdateMenuImage(update *tgbotapi.Update, user *models.User, args []any) {
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

	if update.Message != nil {
		if update.Message.Text == "/delete" {
			menu.ImageFileID = nil
			err = h.database.Model(&menu).UpdateColumn("image_file_id", menu.ImageFileID).Error
			if err != nil {
				panic(err)
			}
			_, err = h.bot.Send(viewmenu.Message(update.FromChat().ID, user, &menu, nil, 0, 0))
			if err != nil {
				panic(err)
			}
			return
		}
	}

	menu.ImageFileID = &update.Message.Photo[0].FileID
	err = h.database.Model(&menu).UpdateColumn("image_file_id", menu.ImageFileID).Error
	if err != nil {
		panic(err)
	}
	_, err = h.bot.Send(viewmenu.Message(update.FromChat().ID, user, &menu, nil, 0, 0))
	if err != nil {
		panic(err)
	}
}
