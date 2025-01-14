package commandhandler

import (
	"coffee-like-helper-bot/config"
	"coffee-like-helper-bot/handlers/step"
	"coffee-like-helper-bot/models"

	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"regexp"
)

type commandFunction func(update *tgbotapi.Update, user *models.User, args []string)

type command struct {
	string           string
	description      string
	argsRegexp       *regexp.Regexp
	activatorRegexps []*regexp.Regexp
	isForStaff       bool
	function         commandFunction
}

func (c *command) help() string {
	text := fmt.Sprint(c.string, " - ", c.description)
	if c.function == nil {
		text += " (недоступна)"
	}
	return text
}

type CommandHandler struct {
	config   *config.Config
	bot      *tgbotapi.BotAPI
	database *gorm.DB

	stepHandler *stephandler.StepHandler

	commands []*command
}

func NewCommandHandler(cfg *config.Config, bot *tgbotapi.BotAPI, database *gorm.DB, stepHandler *stephandler.StepHandler) *CommandHandler {
	h := &CommandHandler{
		config:      cfg,
		bot:         bot,
		database:    database,
		stepHandler: stepHandler,
	}

	h.commands = []*command{
		{
			string:      "/start",
			description: "приветственное сообщение",
			argsRegexp:  regexp.MustCompile(`^\/start$`),
			function:    h.Start,
		},
		{
			string:      "/help",
			description: "помощь",
			argsRegexp:  regexp.MustCompile(`^\/help$`),
			function:    h.Help,
		},
		{
			string:      "/menus",
			description: "список меню",
			argsRegexp:  regexp.MustCompile(`^\/menus$`),
			function:    h.Menus,
		},
		{
			string:      "/menu_[id]",
			description: "получить меню по id",
			argsRegexp:  regexp.MustCompile(`^\/menu(?:_| )(\d+)$`),
			activatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`\/menu`),
				regexp.MustCompile(`\/menu_$`),
			},
			function: h.Menu,
		},
		{
			string:      "/products",
			description: "список всех позиций",
			argsRegexp:  regexp.MustCompile(`^\/products$`),
			function:    nil,
		},
		{
			string:      "/product_[id]",
			description: "получить напиок по id",
			argsRegexp:  regexp.MustCompile(`^\/product(?:_| )(\d+)$`),
			activatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`\/product`),
				regexp.MustCompile(`\/product_$`),
			},
			function: h.Product,
		},
		{
			string:      "/addmenu",
			description: "добавить новое меню",
			argsRegexp:  regexp.MustCompile(`^\/addmenu$`),
			function:    h.AddMenu,
			isForStaff:  true,
		},
		{
			string:      "/addnotification",
			description: "добавить уведомление",
			argsRegexp:  regexp.MustCompile(`^\/addnotification$`),
			function:    nil,
			isForStaff:  true,
		},
		{
			string:      "/debug",
			description: "режим debug",
			argsRegexp:  regexp.MustCompile(`^\/debug$`),
			function:    h.Debug,
			isForStaff:  true,
		},
		{
			string:      "/panic [text]",
			description: "вызвать ошибку",
			argsRegexp:  regexp.MustCompile(`^\/panic (.+)$`),
			activatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`\/panic$`),
			},
			function:   h.Panic,
			isForStaff: true,
		},
	}
	return h
}

func (h *CommandHandler) Process(update *tgbotapi.Update, user *models.User) {
	for _, cmd := range h.commands {
		if cmd.argsRegexp.MatchString(update.Message.Text) {
			if cmd.isForStaff && !user.IsAdministrator {
				h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Недостаточно прав для доступа."))
				return
			}
			if cmd.function == nil {
				h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Команда временно недоступна."))
				return
			}
			cmd.function(update, user, cmd.argsRegexp.FindAllStringSubmatch(update.Message.Text, -1)[0])
			return
		}

		for _, activatorRegexp := range cmd.activatorRegexps {
			if activatorRegexp.MatchString(update.Message.Text) {
				if cmd.isForStaff && !user.IsAdministrator {
					h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Недостаточно прав для доступа."))
					return
				}
				h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, cmd.help()))
				return
			}
		}
	}
	h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Неизвестная команда."))
}
