package viewmenu

import (
	"coffee-like-helper-bot/models"
	"coffee-like-helper-bot/view/menu/product"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Message(chatID int64, user *models.User, menu *models.Menu, products []models.Product, page, limit int) tgbotapi.Chattable {
	var chattable tgbotapi.Chattable

	msgText := "<b>" + menu.Name + "</b>"
	if menu.Description != "" {
		msgText += "\n\n" + menu.Description
	}

	var keyboard tgbotapi.InlineKeyboardMarkup
	if user.IsAdministrator {
		keyboard = *InlineKeyboardEdit(menu)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, viewproduct.InlineKeyboardList(products, fmt.Sprint("product_from_menu ", menu.ID), page, limit).InlineKeyboard...)

	if menu.ImageFileID == nil {
		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = keyboard
		chattable = msg
	} else {
		msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(*menu.ImageFileID))
		msg.Caption = msgText
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = keyboard
		chattable = msg
	}

	return chattable
}

func InlineKeyboardEdit(menu *models.Menu) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⚙️редактировать", fmt.Sprint("edit_menu ", menu.ID)),
	))
	return &keyboard
}

func InlineKeyboardEditList(menu *models.Menu) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("➕добавить напиток", fmt.Sprint("add_product_to_menu ", menu.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✏️изменить название", fmt.Sprint("edit_menu_name ", menu.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📝изменить описание", fmt.Sprint("edit_menu_description ", menu.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🖼изменить изображение", fmt.Sprint("edit_menu_image ", menu.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🗑удалить", fmt.Sprint("delete_menu ", menu.ID))),
	)

	return &keyboard
}

func InlineKeyboardList(menus []models.Menu, callbackData string, page, limit int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()

	if len(menus) == 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("пусто", "nothing")))
		return &keyboard
	}

	totalPages := (len(menus) + limit - 1) / limit
	startIndex := page * limit
	endIndex := (page + 1) * limit
	if endIndex > len(menus) {
		endIndex = len(menus)
	}

	menus = menus[startIndex:endIndex]

	for _, menu := range menus {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(menu.Name, fmt.Sprintf("menu %d", menu.ID))),
		)
	}

	if totalPages > 1 {
		navigationRow := tgbotapi.NewInlineKeyboardRow()
		if page > 0 {
			navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData("<<", fmt.Sprintf("page %s %d", callbackData, page-1)))
		}
		navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("(%d / %d)", page+1, totalPages), "nothing"))
		if page < totalPages-1 {
			navigationRow = append(navigationRow, tgbotapi.NewInlineKeyboardButtonData(">>", fmt.Sprintf("page %s %d", callbackData, page+1)))
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, navigationRow)
	}

	return &keyboard
}
