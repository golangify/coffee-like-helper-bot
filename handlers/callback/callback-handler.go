package callbackhandler

import (
	"coffee-like-helper-bot/config"
	stephandler "coffee-like-helper-bot/handlers/step"
	"coffee-like-helper-bot/models"
	workermailer "coffee-like-helper-bot/workers/mailer"
	"fmt"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type callback struct {
	regexp              *regexp.Regexp
	function            func(update *tgbotapi.Update, user *models.User, args []string)
	isForStaff          bool
	customCallbackQuery bool
}

type CallbackHandler struct {
	config   *config.Config
	bot      *tgbotapi.BotAPI
	database *gorm.DB

	stepHandler *stephandler.StepHandler

	callbacks []*callback

	mailer *workermailer.Mailer
}

func NewCallbackHandler(cfg *config.Config, bot *tgbotapi.BotAPI, database *gorm.DB, stepHandler *stephandler.StepHandler, mailer *workermailer.Mailer) *CallbackHandler {
	h := &CallbackHandler{
		config:   cfg,
		bot:      bot,
		database: database,

		stepHandler: stepHandler,

		mailer: mailer,
	}

	h.callbacks = []*callback{
		{
			regexp:   regexp.MustCompile(`^nothing$`),
			function: func(_ *tgbotapi.Update, _ *models.User, _ []string) {},
		},
		{
			regexp:   regexp.MustCompile(`^menu (\d+)$`),
			function: h.menu,
		},
		{
			regexp:   regexp.MustCompile(`^page all_menus (\d+)$`),
			function: h.pageAllMenus,
		},
		{
			regexp:   regexp.MustCompile(`^product (\d+)$`),
			function: h.product,
		},
		{
			regexp:   regexp.MustCompile(`^page all_products (\d+)$`),
			function: h.pageAllMenus,
		},
		{
			regexp:   regexp.MustCompile(`^page product_from_menu (\d+) (\d+)$`),
			function: h.pageProductFromMenu,
		},
		{
			regexp:     regexp.MustCompile(`^add_product_to_menu (\d+)$`),
			function:   h.addProductToMenu,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_menu (\d+)$`),
			function:   h.editMenu,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_menu_name (\d+)$`),
			function:   h.editMenuName,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_menu_description (\d+)$`),
			function:   h.editMenuDescription,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_menu_image (\d+)$`),
			function:   h.editMenuImage,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^delete_menu (\d+)$`),
			function:   h.deleteMenu,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_product (\d+)$`),
			function:   h.editProduct,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_product_name (\d+)$`),
			function:   h.editProductName,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_product_description (\d+)$`),
			function:   h.editProductDescription,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_product_image (\d+)$`),
			function:   h.editProductImage,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^delete_product (\d+)$`),
			function:   h.deleteProduct,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_user (\d+)$`),
			function:   h.editUser,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^make_barista (\d+)$`),
			function:   h.makeBarista,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^remove_barista (\d+)$`),
			function:   h.removeBarista,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^make_administrator (\d+)$`),
			function:   h.makeAdministrator,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^remove_administrator (\d+)$`),
			function:   h.removeAdministrator,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^user (\d+)$`),
			function:   h.user,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^page users_(all|barista|admin|) (\d+)$`),
			function:   h.pageUsers,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_notification (\d+)$`),
			function:   h.editNotification,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_notification_name (\d+)$`),
			function:   h.editNotificationName,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^edit_notification_user_category (\d+)$`),
			function:   h.editNotificationUserCategory,
			isForStaff: true,
		},
		{
			regexp:     regexp.MustCompile(`^set_notification_user_category (\d+) (all|admin|barista)$`),
			function:   h.setNotificationUserCategory,
			isForStaff: true,
		},
	}

	return h
}

func (h *CallbackHandler) Process(update *tgbotapi.Update, user *models.User) {
	for _, callback := range h.callbacks {
		if callback.isForStaff && !user.IsAdministrator {
			continue
		}
		if callback.regexp.MatchString(update.CallbackQuery.Data) {
			if callback.function == nil {
				h.bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Действие недоступно."))
				return
			}
			if !callback.customCallbackQuery {
				go h.bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "..."))
			}
			callback.function(update, user, callback.regexp.FindAllStringSubmatch(update.CallbackQuery.Data, -1)[0])
			return
		}
	}
	h.bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "произошла ошибка"))
	panic(fmt.Sprint("действие ", update.CallbackQuery.Data, " недоступно"))
}
