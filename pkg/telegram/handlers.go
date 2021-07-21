package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	b.bot.Send(msg)
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда")
	switch message.Command() {
	case "start":
		b.handleStartCommand(msg)

	case "all":
		b.handleAllCommand(msg)

	default:
		b.handleUnknownCommand(msg)
	}

	return nil
}

func (b *Bot) handleStartCommand(msg tgbotapi.MessageConfig) error {
	msg.Text = "Hi"
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleAllCommand(msg tgbotapi.MessageConfig) error {
	groups, err := b.db.GetAllGroupsSlice()
	if err != nil {
		return err
	}
	msg.Text = strings.Join(groups, "\n")
	if len(msg.Text) > 4096 {
		msg.Text = msg.Text[:4096]
	}
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(msg tgbotapi.MessageConfig) error {
	_, err := b.bot.Send(msg)
	return err
}
