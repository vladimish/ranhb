package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/configurator"
	"github.com/telf01/ranhb/pkg/db"
	"log"
	"net/http"
)

type Bot struct {
	db  *db.DataBase
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI, db *db.DataBase) *Bot {
	return &Bot{bot: bot, db: db}
}

func (b *Bot) Start() error {
	err := b.initWebhook()
	if err != nil {
		return err
	}
	err = b.checkWebhookStatus()
	if err != nil {
		return err
	}

	updates := b.bot.ListenForWebhook("/" + b.bot.Token)

	go func() {
		err := http.ListenAndServeTLS(":"+configurator.Cfg.Port, "public.pem", "private.key", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	b.handleUpdates(updates)
	return nil
}

func (b *Bot) checkWebhookStatus() error {
	info, err := b.bot.GetWebhookInfo()
	if err != nil {
		return err
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	return nil
}

func (b *Bot) initWebhook() error {
	response, err := b.bot.RemoveWebhook()
	if err != nil {
		return err
	}
	log.Println(response)

	address := "https://" + configurator.Cfg.Url + ":" + configurator.Cfg.Port + "/"
	response, err = b.bot.SetWebhook(tgbotapi.NewWebhookWithCert(address+b.bot.Token, "public.pem"))
	if err != nil {
		return err
	}
	log.Println(response)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		log.Printf("update recieved: %+v\n", update)

		if update.Message != nil {
			if update.Message.IsCommand() {
				err := b.handleCommand(update.Message)
				if err != nil {
					log.Fatal(err)
				}
				continue
			}

			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}
