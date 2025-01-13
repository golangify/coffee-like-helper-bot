package stephandler

import (
	"coffee-like-helper-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type stepFunction func(update *tgbotapi.Update, user *models.User, args []any)

type step struct {
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

func (h *StepHandler) AddStepHandler(userID uint, function stepFunction, args []any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.handlers[userID] = step{
		function: function,
		args:     args,
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
		step.function(update, user, step.args)
	}

	return ok
}
