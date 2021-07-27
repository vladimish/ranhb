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
	keyboard, err := b.generateKeyboard(b.db.GetAllDistinctField, "form", "fucks")
	msg.Text = "Выберите форму обучения."
	msg.ReplyMarkup = keyboard
	message, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	log.Println("Message sent: ", message)
	return nil
}

type dataDrainer func(args ...string) ([]string, error)

func (b *Bot) generateKeyboard(dd dataDrainer, args ...string) (*tgbotapi.ReplyKeyboardMarkup, error) {
	var buttons [][]tgbotapi.KeyboardButton
	data := make([][]string, 1)
	var err error
	data[0], err = dd(args...)
	if err != nil {
		return nil, err
	}

	for len(data[0]) > 75 {
		for i := range data {
			data = append(data, data[i])
			data[i] = data[i][:len(data[i])/2]
			data[len(data)-1] = data[len(data)-1][len(data[i])/2:]
		}
	}

	for i := 0; i < len(data); i++ {
		var kbButtons []tgbotapi.KeyboardButton
		for k := range data[i] {
			kbButtons = append(kbButtons, tgbotapi.NewKeyboardButton(data[i][k]))
		}
		var row []tgbotapi.KeyboardButton
		row = tgbotapi.NewKeyboardButtonRow(kbButtons...)

		buttons = append(buttons, row)
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons...)

	return &keyboard, nil
}

func (b *Bot) handleAllCommand(msg tgbotapi.MessageConfig) error {
	groups, err := b.db.GetAllDistinctField("groups", "tt")
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
