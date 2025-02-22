package workernotificator

import (
	"coffee-like-helper-bot/config"
	"coffee-like-helper-bot/logger"
	workermailer "coffee-like-helper-bot/workers/mailer"
	"context"
	"fmt"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var log = logger.NewLoggerWithPrefix("notificator")

type Notificator struct {
	bot      *tgbotapi.BotAPI
	database *gorm.DB
	config   *config.Config
	mailer   *workermailer.Mailer
	ctx      context.Context

	activeNotificationsCancel map[uint]context.CancelFunc
	mu                        sync.Mutex
}

func NewNotificator(bot *tgbotapi.BotAPI, database *gorm.DB, config *config.Config, mailer *workermailer.Mailer) (*Notificator, error) {
	if err := database.AutoMigrate(&Notification{}); err != nil {
		return nil, err
	}

	w := &Notificator{
		bot:                       bot,
		database:                  database,
		config:                    config,
		mailer:                    mailer,
		ctx:                       context.Background(),
		activeNotificationsCancel: make(map[uint]context.CancelFunc),
	}

	var notifications []Notification
	if err := w.database.Find(&notifications).Error; err != nil {
		return nil, err
	}

	w.run(notifications)
	log.Println("запущен")

	return w, nil
}

func (w *Notificator) CancelNotification(id uint) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if cancel, ok := w.activeNotificationsCancel[id]; ok {
		cancel()
		delete(w.activeNotificationsCancel, id)
		log.Println("notification", id, "cancelled")
	}
}

func (w *Notificator) NotificationProcess(notification *Notification, withWarning bool) {
	w.CancelNotification(notification.ID)
	ctx, cancel := context.WithCancel(w.ctx)
	w.mu.Lock()
	w.activeNotificationsCancel[notification.ID] = cancel
	w.mu.Unlock()

	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				msg := tgbotapi.NewMessage(0, fmt.Sprintf("worker/notificator: ошибка в процессе уведомления \"%s\": %v\n/notification_%d", notification.Name, err, notification.ID))
				w.mailer.Administrator(&msg)
			}
		}()

		for {
			sleepTime, err := notification.TimeUntilNextNotification()
			if err != nil {
				log.Println(err)
				if !withWarning {
					return
				}
				msg := tgbotapi.NewMessage(0, fmt.Sprintf("Ошибка подсчёта времени для уведомления: \"%s\": %v.\n\nВозможно стоит проверить настройки - /notification_%d", notification.Name, err, notification.ID))
				w.mailer.Administrator(&msg)
				return
			}

			log.Println(fmt.Sprintf("Уведомление \"%s\" будет разослано через %v", notification.Name, sleepTime))
			select {
			case <-ctx.Done():
				return
			case <-time.After(sleepTime):
			}

			notification, err = w.NotificationByID(notification.ID)
			if err != nil || notification == nil {
				if err != nil && err != gorm.ErrRecordNotFound {
					panic(err)
				}
				return
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
				panic(fmt.Sprintf("неизвестная категория пользователей для рассылки: %s", notification.UserCategory))
			}
		}
	}(ctx)
}

func (w *Notificator) run(notifications []Notification) {
	defer func() {
		if w.bot.Debug {
			return
		}
		if err := recover(); err != nil {
			log.Println(err)
			msg := tgbotapi.NewMessage(0, fmt.Sprintf("worker/notificator упал с ошибкой: %v", err))
			w.mailer.Administrator(&msg)
		}
	}()

	for _, ntfctn := range notifications {
		ntfctn := ntfctn
		go w.NotificationProcess(&ntfctn, true)
	}
}

func (w *Notificator) AddNotification(notification *Notification) error {
	return w.database.Create(&notification).Error
}

func (w *Notificator) DeleteNotification(id uint) error {
	if err := w.database.Where("id = ?", id).Delete(&Notification{}).Error; err != nil {
		return err
	}
	w.CancelNotification(id)
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

func (w *Notificator) NumRunning() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.activeNotificationsCancel)
}
