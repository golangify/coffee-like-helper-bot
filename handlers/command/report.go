package commandhandler

import (
	"coffee-like-helper-bot/models"
	viewuser "coffee-like-helper-bot/view/user"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *CommandHandler) report(update *tgbotapi.Update, user *models.User, _ []string) {
	h.stepHandler.AddStepHandler(user.ID, h.StepReport, nil)
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Опиши ошибку:\nпосле каких действий она возникает и скриншот ошибки. Заранее спасибо!\n\n/cancel - отмена"))
}

func (h *CommandHandler) StepReport(update *tgbotapi.Update, user *models.User, _ []any) {
	var admins []models.User
	err := h.database.Find(&admins, "is_administrator = ?", true).Error
	if err != nil {
		panic(err)
	}
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Отправляю отчёт - это может занят немного времени..."))
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for _, admin := range admins {
		h.bot.Send(tgbotapi.NewMessage(admin.TelegramID, fmt.Sprint("Пользователь ", viewuser.Text(user), " отправил отчёт об ошибке:")))
		msg := tgbotapi.NewCopyMessage(admin.TelegramID, update.FromChat().ID, update.Message.MessageID)
		h.bot.Send(msg)
		<-t.C
	}
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Твоё сообщение отправлено. Спасибо!"))
}
