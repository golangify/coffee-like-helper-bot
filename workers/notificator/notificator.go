package workernotificator

import (
	"coffee-like-helper-bot/config"
	"coffee-like-helper-bot/logger"
	workermailer "coffee-like-helper-bot/workers/mailer"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var log = logger.NewLoggerWithPrefix("notificator")

type Notificator struct {
	bot      *tgbotapi.BotAPI
	database *gorm.DB
	config   *config.Config

	mailer *workermailer.Mailer
}

func (w *Notificator) notificationProcess(notification Notification) {
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

		if notification, err := w.NotificationByID(notification.ID); err != nil || notification != nil {
			if err != nil {
				panic(err)
			}
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

func (w *Notificator) run(notifications []Notification) {
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

	for _, ntfctn := range notifications {
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

	var notifications []Notification
	if err := w.database.Find(&notifications).Error; err != nil {
		return nil, err
	}

	w.run(notifications)
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

	go w.notificationProcess(notification)

	return nil
}

func (w *Notificator) DeleteNotification(id uint) error {

	err := w.database.Model(&Notification{}).Delete(id).Error
	if err != nil {
		return err
	}

	return nil
}

func (n *Notificator) NotificationByID(id uint) (*Notification, error) {
	var notification Notification
	err := n.database.First(&notification, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &notification, nil
}
