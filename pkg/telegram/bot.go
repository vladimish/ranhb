package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/db"
	"log"
)

type Bot struct {
	db *db.DataBase
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI, db *db.DataBase) *Bot {
	return &Bot{bot: bot, db: db}
}

func (b *Bot) Start() error {
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand(){
			err := b.handleCommand(update.Message)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}
