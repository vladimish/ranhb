package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vladimish/ranhb/internal/configurator"
	"github.com/vladimish/ranhb/internal/db"
	"github.com/vladimish/ranhb/internal/telegram"
	"github.com/vladimish/yookassa-go-sdk"
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

	kassa := yookassa.NewKassa(configurator.Cfg.ShopID, configurator.Cfg.ShopToken)
	result, err := kassa.Ping()
	if !result || err != nil {
		log.Fatal("Can't ping kassa", err)
	}

	telegramBot := telegram.NewBot(bot, botDB, kassa)

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
