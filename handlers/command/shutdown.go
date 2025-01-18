package commandhandler

import (
	"coffee-like-helper-bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) shutdown(update *tgbotapi.Update, user *models.User, _ []string) {
	checkPhrase := "выключить " + h.bot.Self.FirstName
	h.stepHandler.AddStepHandler(user.ID, h.stepShutdown, []any{checkPhrase})
	msg := tgbotapi.NewMessage(update.FromChat().ID, "Отправь(можно нажать на текст ниже для копирования)\n\n<code>"+checkPhrase+"</code>\n\nчтобы подтвердить выключение бота.\n\nотмена - /cancel")
	msg.ParseMode = tgbotapi.ModeHTML
	h.bot.Send(msg)
}

// args: checkPhrase string
func (h *CommandHandler) stepShutdown(update *tgbotapi.Update, user *models.User, args []any) {
	checkPhrase := args[0].(string)
	if update.Message == nil || update.Message.Text == "" {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "сообщение не содержит текст."))
		return
	}
	if checkPhrase != update.Message.Text {
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Неправильная подтверждающая фраза. Выключение отменено."))
		return
	}
	h.bot.StopReceivingUpdates()
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Бот выключен."))
}
