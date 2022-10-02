package main

import (
	"birthdayBot/internal/configs"
	"birthdayBot/internal/db"
	"birthdayBot/internal/telegram"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	conf, err := configs.LoadConfiguration("config.json")
	if err != nil {
		fmt.Printf("error config, %s", err)
		panic(1)
	}
	dbCon, err := db.ConnectToDb(conf)
	if err != nil {
		fmt.Printf("error ConnectingToDb, %s", err)
		panic(1)
	}

	botApi, err := tgbotapi.NewBotAPI(conf.TelegramBotToken)
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	bot := telegram.NewBot(botApi)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	if err := bot.Start(dbCon, &conf); err != nil {
		log.Fatal(err)
	}

}
