package main

import (
	"coffee-like-helper-bot/config"
	updatehandler "coffee-like-helper-bot/handlers/update"
	"coffee-like-helper-bot/models"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := config.JsonLoadFromFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	database, err := gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	if err = database.AutoMigrate(&models.User{}, &models.Menu{}, &models.Product{}, &models.Search{}); err != nil {
		log.Fatalln(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramApiToken)
	if err != nil {
		log.Fatalln(err)
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
					log.Println(err)
					bot.Send(tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprint("Произошла ошибка: ", err, "\n\nЕсли это что-то серьезное пиши - @golangify")))
				}
			}()
			updHandler.Process(&update)
		}(update)
	}
}
