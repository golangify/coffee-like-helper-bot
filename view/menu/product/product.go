package viewproduct

import (
	"coffee-like-helper-bot/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Message(chatID int64, user *models.User, product *models.Product) tgbotapi.Chattable {
	var chattable tgbotapi.Chattable

	msgText := "<b>" + product.Name + "</b>"
	if product.Description != "" {
		msgText += "\n\n" + product.Description
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	if user.IsAdministrator {
		keyboard = *InlineKeyboardEdit(product)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("<< к меню", fmt.Sprint("menu ", product.MenuID)),
	))

	if product.ImageFileID == nil {
		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = keyboard
		chattable = msg
	} else {
		msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(*product.ImageFileID))
		msg.Caption = msgText
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = keyboard
		chattable = msg
	}

	return chattable
}

func InlineKeyboardEdit(product *models.Product) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⚙️редактировать", fmt.Sprint("edit_product ", product.ID)),
	))
	return &keyboard
}

func InlineKeyboardEditList(product *models.Product) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✏️изменить название", fmt.Sprint("edit_product_name ", product.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📝изменить описание", fmt.Sprint("edit_product_description ", product.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🖼изменить изображение", fmt.Sprint("edit_product_image ", product.ID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🗑удалить", fmt.Sprint("delete_product ", product.ID))),
	)

	return &keyboard
}

func InlineKeyboardList(products []models.Product, callbackData string, page, limit int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()

	if len(products) == 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("пусто", "nothing")))
		return &keyboard
	}

	totalPages := (len(products) + limit - 1) / limit
	startIndex := page * limit
	endIndex := (page + 1) * limit
	if endIndex > len(products) {
		endIndex = len(products)
	}

	products = products[startIndex:endIndex]

	for _, product := range products {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(product.Name, fmt.Sprintf("product %d", product.ID))),
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
