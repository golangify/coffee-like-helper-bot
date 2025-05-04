package commandhandler

import (
	"coffee-like-helper-bot/config"
	stephandler "coffee-like-helper-bot/handlers/step"
	"coffee-like-helper-bot/models"
	menuservice "coffee-like-helper-bot/service/menu"
	workernotificator "coffee-like-helper-bot/workers/notificator"

	"fmt"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type commandFunction func(update *tgbotapi.Update, user *models.User, args []string)

type command struct {
	string           string
	description      string
	argsRegexp       *regexp.Regexp
	activatorRegexps []*regexp.Regexp
	isForStaff       bool
	function         commandFunction
	isHidden         bool
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

	menuService *menuservice.MenuService

	stepHandler *stephandler.StepHandler

	notificator *workernotificator.Notificator

	commands []*command
}

func NewCommandHandler(cfg *config.Config, bot *tgbotapi.BotAPI, database *gorm.DB, menuservice *menuservice.MenuService, stepHandler *stephandler.StepHandler, notificator *workernotificator.Notificator) *CommandHandler {
	h := &CommandHandler{
		config:   cfg,
		bot:      bot,
		database: database,

		menuService: menuservice,

		stepHandler: stepHandler,

		notificator: notificator,
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
			argsRegexp:  regexp.MustCompile(`^\/help(?:|(?:_| )(all))$`),
			function:    h.Help,
		},
		{
			string:      "/menus",
			description: "список меню",
			argsRegexp:  regexp.MustCompile(`^\/menus$`),
			function:    h.Menus,
		},
		{
			string:      "/report",
			description: "сообщить об ошибке",
			argsRegexp:  regexp.MustCompile(`^\/report$`),
			function:    h.report,
		},
		{
			string:      "/menu [id]",
			description: "получить меню по id",
			argsRegexp:  regexp.MustCompile(`^\/menu(?:_| )(\d+)$`),
			activatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`\/menu$`),
				regexp.MustCompile(`\/menu_$`),
			},
			function: h.Menu,
			isHidden: true,
		},
		{
			string:      "/addmenu",
			description: "добавить новое меню",
			argsRegexp:  regexp.MustCompile(`^\/addmenu$`),
			function:    h.AddMenu,
			isForStaff:  true,
		},
		{
			string:      "/faddnotification <название> [all | barista | admin] [<дни недели>] <час>:<минута> <содержимое>",
			description: "быстро добавить уведомление",
			argsRegexp:  regexp.MustCompile(`^\/faddnotification (.+) (all|barista|admin) (\[(?:[1-7](?:(?:,|, )[1-7]){0,6}){1,7}\]) ([0-2][0-9]:[0-5][0-9]) (.+)`),
			activatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`\/faddnotification`),
			},
			function:   nil,
			isForStaff: true,
			isHidden:   true,
		},
		{
			string:      "/addnotification",
			description: "добавить уведомление",
			argsRegexp:  regexp.MustCompile(`^\/addnotification$`),
			function:    h.addNotification,
			isForStaff:  true,
		},
		{
			string:      "/notifications",
			description: "список уведомлений",
			argsRegexp:  regexp.MustCompile(`^\/notifications$`),
			function:    h.notifications,
			isForStaff:  true,
		},
		{
			string:      "/notification [id]",
			description: "уведомление по id",
			argsRegexp:  regexp.MustCompile(`^\/notification(?:|(?: |_))(\d+)$`),
			activatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`^\/notification`),
			},
			function:   h.notification,
			isForStaff: true,
			isHidden:   true,
		},
		{
			string:      "/users [ all, admin, barista ](по умолчанию all)",
			description: "список пользователей",
			argsRegexp:  regexp.MustCompile(`^\/users(?:|(?: |_)(all|admin|barista))$`),
			activatorRegexps: []*regexp.Regexp{
				regexp.MustCompile(`^\/users.+`),
			},
			function:   h.Users,
			isForStaff: true,
		},
		{
			string:      "/debug",
			description: "режим debug",
			argsRegexp:  regexp.MustCompile(`^\/debug$`),
			function:    h.Debug,
			isForStaff:  true,
			isHidden:    true,
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
			isHidden:   true,
		},
		{
			string:      "/info",
			description: "тех. информация",
			argsRegexp:  regexp.MustCompile(`^\/info$`),
			function:    h.info,
			isForStaff:  true,
			isHidden:    true,
		},
		{
			string:      "/shutdown",
			description: "выключить бота",
			argsRegexp:  regexp.MustCompile(`^\/shutdown$`),
			function:    h.shutdown,
			isForStaff:  true,
			isHidden:    true,
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
