package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/configurator"
	"github.com/telf01/ranhb/pkg/users"
	"log"
	"strings"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	user, err := users.Get(message.Chat.ID, b.db)
	if err != nil {
		return err
	}
	if user.U == nil {
		user.Init(message.Chat.ID, b.db)
		user.U.LastActionValue = -configurator.Cfg.PageSize
	}

	switch user.U.LastAction {
	case "start":
		{
			switch message.Text {
			case "➡️️":
				groups, err := b.db.GetAllDistinctField("groups", "tt", "0", "10000")
				if err != nil {
					return err
				}
				if user.U.LastActionValue >= len(groups) {
					return b.sendKeyboard(message.Chat.ID, user, 0)
				} else {
					return b.sendKeyboard(message.Chat.ID, user, configurator.Cfg.PageSize)
				}

			case "⬅️":
				if user.U.LastActionValue <= 0 {
					return b.sendKeyboard(message.Chat.ID, user, 0)
				} else {
					return b.sendKeyboard(message.Chat.ID, user, -configurator.Cfg.PageSize)
				}
			default:
				log.Println("fuck")
			}
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда")
	switch message.Command() {
	case "start":
		err := b.handleStartCommand(msg, message.Chat.ID)
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

func (b *Bot) handleStartCommand(msg tgbotapi.MessageConfig, id int64) error {
	user, err := users.Get(id, b.db)
	if err != nil {
		return err
	}
	if user.U.LastAction != "start" {
		user.U.LastActionValue = -configurator.Cfg.PageSize
		user.U.LastAction = "start"
		err := user.Save()
		if err != nil {
			return err
		}
	}
	if user.U.Id == 0 {
		user.Init(id, b.db)
		user.U.LastAction = "start"
		user.U.LastActionValue = -configurator.Cfg.PageSize
		err := user.Save()
		if err != nil {
			return err
		}
	}

	return b.sendKeyboard(msg.ChatID, user, configurator.Cfg.PageSize)
}

func (b *Bot) sendKeyboard(chatId int64, user *users.User, pageOffset int) error {
	msg := tgbotapi.NewMessage(chatId, "Выберите форму обучения.")

	user.U.LastActionValue += pageOffset
	err := user.Save()
	if err != nil {
		return err
	}

	keyboard, err := b.generateKeyboard(b.db.GetAllDistinctField, user, "groups", "tt", fmt.Sprintf("%d", user.U.LastActionValue), fmt.Sprintf("%d", configurator.Cfg.PageSize))
	if err != nil {
		return err
	}

	msg.ReplyMarkup = keyboard
	if user.U.LastActionValue == configurator.Cfg.PageSize {
		msg.Text = "Выберите форму обучения."
	}
	message, err := b.bot.Send(msg)

	if err != nil {
		return err
	}
	log.Println("Message sent: ", message)

	return nil
}

type dataDrainer func(args ...string) ([]string, error)

func (b *Bot) generateKeyboard(dd dataDrainer, u *users.User, args ...string) (*tgbotapi.ReplyKeyboardMarkup, error) {
	var buttons [][]tgbotapi.KeyboardButton
	groups, err := dd(args...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(groups); i++ {
		var row []tgbotapi.KeyboardButton
		if i+1 < len(groups) {
			row = tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(groups[i]),
				tgbotapi.NewKeyboardButton(groups[i+1]),
			)
			i += 1
		} else {
			row = tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(groups[i]),
			)
		}
		buttons = append(buttons, row)
	}

	var row []tgbotapi.KeyboardButton
	if len(groups) < 30 {
		row = tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️"),
		)
	} else if u.U.LastActionValue <= 0 {
		row = tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➡️️"),
		)
	} else {
		row = tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️"),
			tgbotapi.NewKeyboardButton("➡️️"),
		)
	}
	buttons = append(buttons, row)

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
