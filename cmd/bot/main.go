package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/configurator"
	"github.com/telf01/ranhb/pkg/db"
	"github.com/telf01/ranhb/pkg/telegram"
	"log"
)

func main() {
	// Init database connection.
	base, err := db.InitDBConncetion(configurator.Cfg.DbLogin, configurator.Cfg.DbPassword)
	if err != nil {
		log.Fatal(err)
	}
	botDB := db.NewDataBase(base)

	// Authorize bot.
	bot, err := Authorize()
	if err != nil {
		log.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, botDB)

	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func Authorize() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(configurator.Cfg.TgKey)
	if err != nil {
		return bot, err
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot, nil
}
