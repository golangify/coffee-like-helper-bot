package workernotificator

import (
	"coffee-like-helper-bot/config"
	"coffee-like-helper-bot/logger"
	workermailer "coffee-like-helper-bot/workers/mailer"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

var log = logger.NewLoggerWithPrefix("notificator")

type Notificator struct {
	bot      *tgbotapi.BotAPI
	database *gorm.DB
	config   *config.Config

	mailer *workermailer.Mailer

	notifications []*Notification
}

func (w *Notificator) notificationProcess(notification *Notification) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
			msg := tgbotapi.NewMessage(0, fmt.Sprint("worker/notificator: ошибка в процессе уведомления \"", notification.Name, "\"", err))
			w.mailer.Administrator(&msg)
		}
	}()

	for {
		sleepTime, err := notification.TimeUntilNextNotification()
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprint("Уведомление \"", notification.Name, "\" будет разослано через ", sleepTime))
		time.Sleep(sleepTime)

		if !w.NotificationExist(notification.ID) {
			break
		}

		msg := tgbotapi.NewMessage(0, notification.Text)
		switch notification.UserCategory {
		case "all":
			w.mailer.AllUsers(&msg)
		case "barista":
			w.mailer.Barista(&msg)
		case "admin":
			w.mailer.Administrator(&msg)
		default:
			panic(fmt.Sprint("неизвестная категория пользователей для рассылки: ", notification.UserCategory))
		}
	}
}

func (w *Notificator) run() {
	defer func() {
		if w.bot.Debug {
			return
		}
		err := recover()
		if err == nil {
			return
		}
		log.Println(err)
		msg := tgbotapi.NewMessage(0, fmt.Sprint("worker/notificator упал с ошибкой: ", err))
		w.mailer.Administrator(&msg)
	}()

	for _, ntfctn := range w.notifications {
		go w.notificationProcess(ntfctn)
	}
}

func NewNotificator(bot *tgbotapi.BotAPI, database *gorm.DB, config *config.Config, mailer *workermailer.Mailer) (*Notificator, error) {
	if err := database.AutoMigrate(&Notification{}); err != nil {
		return nil, err
	}

	w := &Notificator{
		bot:      bot,
		database: database,
		config:   config,
		mailer:   mailer,
	}

	var notifications []*Notification
	if err := w.database.Find(&notifications).Error; err != nil {
		return nil, err
	}

	w.notifications = notifications

	w.run()
	log.Println("запущен")

	return w, nil
}

func (w *Notificator) AddNotification(notification Notification) error {
	_, err := notification.TimeUntilNextNotification()
	if err != nil {
		return err
	}
	err = w.database.Create(&notification).Error
	if err != nil {
		return err
	}

	w.notifications = append(w.notifications, &notification)

	go w.notificationProcess(&notification)

	return nil
}

func (w *Notificator) DeleteNotification(id uint) error {
	var notificationToDelete *Notification
	w.notifications = slices.DeleteFunc(w.notifications, func(n *Notification) bool {
		if n.ID == id {
			notificationToDelete = n
			return true
		}
		return false
	})
	if notificationToDelete == nil {
		return fmt.Errorf("notification with id %d not found", id)
	}

	err := w.database.Delete(notificationToDelete).Error
	if err != nil {
		return err
	}

	return nil
}

func (n *Notificator) NotificationExist(id uint) bool {
	for _, n := range n.notifications {
		if n.ID == id {
			return true
		}
	}

	return false
}
