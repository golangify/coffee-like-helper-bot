package main

import (
	"coffee-like-helper-bot/config"
	updatehandler "coffee-like-helper-bot/handlers/update"
	"coffee-like-helper-bot/logger"
	"coffee-like-helper-bot/models"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log = logger.NewLoggerWithPrefix("MAIN")

func main() {
	config, err := config.JsonLoadFromFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	database, err := gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err = database.AutoMigrate(&models.User{}, &models.Menu{}, &models.Product{}, &models.Search{}); err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramApiToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("авторизован в @%s\n", bot.Self.UserName)

	updHandler := updatehandler.New(config, bot, database)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		go func(update tgbotapi.Update) {
			defer func() {
				if bot.Debug {
					return
				}
				if err := recover(); err != nil {
					log.Printf("%+v", err)
					bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Произошла ошибка: ", err, "\n\nЕсли ошибка критическая - /report")))
				}
			}()
			updHandler.Process(&update)
		}(update)
	}
}
