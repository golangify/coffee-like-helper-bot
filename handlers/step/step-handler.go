package stephandler

import (
	"coffee-like-helper-bot/models"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	stepTypeAny = iota
	stepTypeText
	stepTypeImage
)

type stepFunction func(update *tgbotapi.Update, user *models.User, args []any)

type step struct {
	stepType int
	user     *models.User
	function stepFunction
	args     []any
}

type StepHandler struct {
	bot      *tgbotapi.BotAPI
	handlers map[uint]step
	mu       sync.Mutex
}

func New(bot *tgbotapi.BotAPI) *StepHandler {
	h := &StepHandler{
		bot:      bot,
		handlers: make(map[uint]step),
	}
	return h
}

func (h *StepHandler) addStepHandler(stepType int, user *models.User, function stepFunction, args []any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.handlers[user.ID] = step{
		stepType: stepTypeAny,
		user:     user,
		function: function,
		args:     args,
	}
	h.bot.Send(tgbotapi.NewMessage(user.TelegramID, "Чтобы отменить действие отправь /cancel"))
}

func (h *StepHandler) AddAny(user *models.User, function stepFunction, args []any) {
	h.addStepHandler(stepTypeAny, user, function, args)
}

func (h *StepHandler) AddText(user *models.User, function stepFunction, args []any) {
	h.addStepHandler(stepTypeText, user, function, args)
}

func (h *StepHandler) AddImage(user *models.User, function stepFunction, args []any) {
	h.addStepHandler(stepTypeText, user, function, args)
}

func typeMatch(step *step, update *tgbotapi.Update) bool {
	switch step.stepType {
	case stepTypeAny:
		return true
	case stepTypeText:
		return update.Message != nil && update.Message.Text != ""
	case stepTypeImage:
		return update.Message != nil &&
			(update.Message.Text == "/delete" ||
				(update.Message.Photo != nil && len(update.Message.Photo) > 0))
	default:
		return false
	}
}

func (h *StepHandler) Process(update *tgbotapi.Update, user *models.User) bool {
	h.mu.Lock()
	step, ok := h.handlers[user.ID]
	delete(h.handlers, user.ID)
	h.mu.Unlock()

	if ok {
		if update.Message != nil && update.Message.Text == "/cancel" {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Действие отменено."))
			return ok
		}

		if !typeMatch(&step, update) {
			h.bot.Send(tgbotapi.NewMessage(step.user.TelegramID, "Не принято, повтори. Отмена - /cancel"))
			return false
		}

		step.function(update, user, step.args)
	}

	return ok
}
