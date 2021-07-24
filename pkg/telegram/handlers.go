package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда")
	switch message.Command() {
	case "start":
		err := b.handleStartCommand(msg)
		if err != nil {
			return err
		}

	case "all":
		err := b.handleAllCommand(msg)
		if err != nil {
			return err
		}

	default:
		err := b.handleUnknownCommand(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handleStartCommand(msg tgbotapi.MessageConfig) error {
	keyboard, err := b.generateKeyboard()
	msg.Text = "Hi"
	msg.ReplyMarkup = keyboard
	message, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	log.Println("Message sent: ", message)
	return nil
}

func (b *Bot) generateKeyboard() (*tgbotapi.ReplyKeyboardMarkup, error) {
	var buttons [][]tgbotapi.KeyboardButton
	groups, err := b.db.GetAllGroupsSlice()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(groups); i++ {
		var row []tgbotapi.KeyboardButton
		if i+3 < len(groups) {
			row = tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(groups[i]),
				tgbotapi.NewKeyboardButton(groups[i+1]),
				tgbotapi.NewKeyboardButton(groups[i+2]),
				tgbotapi.NewKeyboardButton(groups[i+3]),
			)
			i += 3
		} else {
			row = tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(groups[i]),
			)
		}
		buttons = append(buttons, row)
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons...)

	return &keyboard, nil
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
