package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/db"
	"github.com/telf01/ranhb/pkg/telegram"
	"log"
	"os"
)

func main(){
	// Init database connection.
	base, err := db.InitDBConncetion(os.Getenv("TGLOGIN"), os.Getenv("TGPASS"))
	if err != nil{
		log.Fatal(err)
	}
	botDB := db.NewDataBase(base)

	// Authorize bot.
	bot, err := Authorize()
	if err != nil{
		log.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, botDB)

	if err := telegramBot.Start(); err!=nil{
		log.Fatal(err)
	}
}

func Authorize() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGKEY"))
	if err != nil {
		return bot, err
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot, nil
}