package messagehandler

import (
	"coffee-like-helper-bot/models"
	viewproduct "coffee-like-helper-bot/view/menu/product"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func prepareQueryString(queryString string) string {
	tokens := map[string]string{
		"е":  "ё",
		"й":  "и",
		"э":  "ё",
		"сс": "с",
		"а":  "о",
	}

	queryRuneSlice := []rune(strings.ToLower(queryString))
	if len(queryRuneSlice) > 100 {
		queryRuneSlice = queryRuneSlice[:100]
	}
	queryString = string(queryRuneSlice)

	for s, d := range tokens {
		queryString = strings.ReplaceAll(queryString, s, d)
	}

	return queryString
}

func (h *MessageHandler) Search(update *tgbotapi.Update, user *models.User) {
	if len([]rune(update.Message.Text)) < 3 {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Запрос должен быть больше 3 символов."))
		return
	}

	h.database.Create(&models.Search{
		Text:   update.Message.Text,
		UserID: user.ID,
	})

	results, err := h.searchEngine.SearchProducts(update.Message.Text)
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Запрос не дал результатов"))
		return
	}

	msg := tgbotapi.NewMessage(update.FromChat().ID, "Результат запроса:")
	msg.ReplyMarkup = viewproduct.InlineKeyboardList(results, "search_product", 0, 100)

	h.bot.Send(msg)
}
