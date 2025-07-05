package updatehandler

import (
	"coffee-like-helper-bot/config"
	callbackhandler "coffee-like-helper-bot/handlers/callback"
	commandhandler "coffee-like-helper-bot/handlers/command"
	messagehandler "coffee-like-helper-bot/handlers/message"
	stephandler "coffee-like-helper-bot/handlers/step"
	"coffee-like-helper-bot/models"
	menuservice "coffee-like-helper-bot/service/menu"
	userservice "coffee-like-helper-bot/service/user"
	viewuser "coffee-like-helper-bot/view/user"
	workermailer "coffee-like-helper-bot/workers/mailer"
	workernotificator "coffee-like-helper-bot/workers/notificator"

	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type UpdateHandler struct {
	config   *config.Config
	bot      *tgbotapi.BotAPI
	database *gorm.DB

	userService *userservice.UserService
	menuService *menuservice.MenuService

	commandHandler  *commandhandler.CommandHandler
	callbackHandler *callbackhandler.CallbackHandler
	messageHandler  *messagehandler.MessageHandler

	stepHandler *stephandler.StepHandler

	mailer *workermailer.Mailer
}

func New(cfg *config.Config, bot *tgbotapi.BotAPI, database *gorm.DB) *UpdateHandler {
	userService := userservice.NewUserService(database)
	menuService := menuservice.NewMenuService(database)

	stepHandler := stephandler.New(bot)
	mailer := workermailer.New(bot, database, cfg)
	notificator, err := workernotificator.NewNotificator(bot, database, cfg, mailer)
	if err != nil {
		panic(err)
	}

	h := &UpdateHandler{
		config:   cfg,
		bot:      bot,
		database: database,

		userService: userService,
		menuService: menuService,

		commandHandler:  commandhandler.NewCommandHandler(cfg, bot, database, menuService, stepHandler, notificator),
		callbackHandler: callbackhandler.NewCallbackHandler(cfg, bot, database, userService, stepHandler, mailer, notificator),
		messageHandler:  messagehandler.NewMessageHandler(bot, database),

		stepHandler: stepHandler,

		mailer: mailer,
	}

	return h
}

func (h *UpdateHandler) Process(update *tgbotapi.Update) {
	sentFrom := update.SentFrom()
	if sentFrom == nil {
		return
	}

	var user models.User
	err := h.database.First(&user, "telegram_id = ?", sentFrom.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user = models.User{
				TelegramID: sentFrom.ID,
			}
			if err = h.userService.NewUser(&user); err != nil {
				panic(err)
			}
			err = h.userService.UpdateUserInfo(&user, sentFrom)
			if err != nil {
				panic(err)
			}
			h.сallHandlers(update, sentFrom, &user)
			return
		}
		panic(err)
	}
	err = h.userService.UpdateUserInfo(&user, sentFrom)
	if err != nil {
		panic(err)
	}
	// TODO: проверка на флуд
	// ...
	h.сallHandlers(update, sentFrom, &user)
}

func isFlood() {}

func (h *UpdateHandler) сallHandlers(update *tgbotapi.Update, sentFrom *tgbotapi.User, user *models.User) {
	if user.IsBanned {
		return
	}
	if !user.IsBarista && !user.IsAdministrator {
		msg := tgbotapi.NewMessage(0, fmt.Sprint(viewuser.Text(user), " запрашивает доступ к боту"))
		msg.ReplyMarkup = viewuser.InlineKeyboardEdit(user)
		go h.mailer.Administrator(&msg)
		h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Запрос на получение доступа отправлен. Ожидай."))
		return
	}

	if h.stepHandler.Process(update, user) {
		return
	}
	if update.Message != nil {
		if update.Message.IsCommand() {
			h.commandHandler.Process(update, user)
			return
		} else if update.Message.Text != "" {
			h.messageHandler.Process(update, user)
			return
		} else {
			h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Действие не поддерживается."))
			return
		}
	} else if update.CallbackQuery != nil {
		h.callbackHandler.Process(update, user)
		return
	}
	h.bot.Send(tgbotapi.NewMessage(sentFrom.ID, "Неподдерживаемое действие."))
}
