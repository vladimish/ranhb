package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/configurator"
	"github.com/telf01/ranhb/pkg/users"
	"log"
	"strings"
	"time"
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
		err := b.handleStartMessage(message, user)
		if err != nil {
			return err
		}
	case "menu":
		err := b.handleMenuMessage(message, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handleMenuMessage(message *tgbotapi.Message, user *users.User) error {
	switch message.Text {
	case configurator.Cfg.Consts.Today:
		tts, err := b.db.GetSpecificTt(user.U.Group, time.Now().Day(), int(time.Now().Month()))
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%+v", tts))
		m, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
		log.Println(m)
	}
	return nil
}

func (b *Bot) handleStartMessage(message *tgbotapi.Message, user *users.User) error {
	switch message.Text {
	case configurator.Cfg.Consts.Right:
		groups, err := b.db.GetAllDistinctField("groups", "tt", "0", "10000")
		if err != nil {
			return err
		}
		if user.U.LastActionValue >= len(groups) {
			return b.sendGroupsKeyboard(message.Chat.ID, user, 0)
		} else {
			return b.sendGroupsKeyboard(message.Chat.ID, user, configurator.Cfg.PageSize)
		}

	case configurator.Cfg.Consts.Left:
		if user.U.LastActionValue <= 0 {
			return b.sendGroupsKeyboard(message.Chat.ID, user, 0)
		} else {
			return b.sendGroupsKeyboard(message.Chat.ID, user, -configurator.Cfg.PageSize)
		}
	default:
		groups, err := b.db.GetAllDistinctFieldWhere("groups", "tt", "0", "1", "groups", message.Text)
		if err != nil {
			return err
		}

		if len(groups) == 1 {
			user.U.LastActionValue = 0
			user.U.LastAction = "menu"
			user.U.Group = groups[0]
			err := user.Save()
			if err != nil {
				return err
			}

			msg := tgbotapi.NewMessage(message.Chat.ID, "Ваша группа сохранена.")
			m, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
			log.Println("Message sent: ", m)

			err = b.sendMenuKeyboard(message.Chat.ID, user)
			if err != nil {
				return err
			}
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Группа не найдена.")
			m, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
			log.Println("Message sent: ", m)
		}
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

	if user.U.Id == 0 || user.U.LastAction != "start" {
		if user.U.Id == 0 {
			user.Init(id, b.db)
		}
		user.U.LastAction = "start"
		user.U.LastActionValue = 0
		err := user.Save()
		if err != nil {
			return err
		}
	}

	return b.sendGroupsKeyboard(msg.ChatID, user, 0)
}

func (b *Bot) sendMenuKeyboard(chatId int64, user *users.User) error {
	msg := tgbotapi.NewMessage(chatId, "Меню")
	user.U.LastAction = "menu"
	err := user.Save()
	if err != nil {
		return err
	}

	keyboard := b.generateMenuKeyboard()
	msg.ReplyMarkup = keyboard

	message, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	log.Println("Message sent: ", message)

	return nil
}

func (b *Bot) sendGroupsKeyboard(chatId int64, user *users.User, pageOffset int) error {
	msg := tgbotapi.NewMessage(chatId, "Выберите форму обучения.")

	user.U.LastActionValue += pageOffset
	err := user.Save()
	if err != nil {
		return err
	}

	keyboard, err := b.generateGroupsKeyboard(b.db.GetAllDistinctField, user, "groups", "tt", fmt.Sprintf("%d", user.U.LastActionValue), fmt.Sprintf("%d", configurator.Cfg.PageSize))
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

func (b *Bot) generateMenuKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton
	row := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Today),
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Tomorrow),
	)
	buttons = append(buttons, row)
	keyboard := tgbotapi.NewReplyKeyboard(buttons...)

	return &keyboard
}

type dataDrainer func(args ...string) ([]string, error)

func (b *Bot) generateGroupsKeyboard(dd dataDrainer, u *users.User, args ...string) (*tgbotapi.ReplyKeyboardMarkup, error) {
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
	if len(groups) < configurator.Cfg.PageSize {
		row = tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Left),
		)
	} else if u.U.LastActionValue <= 0 {
		row = tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Right),
		)
	} else {
		row = tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Left),
			tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Right),
		)
	}
	buttons = append(buttons, row)

	keyboard := tgbotapi.NewReplyKeyboard(buttons...)

	return &keyboard, nil
}

func (b *Bot) handleAllCommand(msg tgbotapi.MessageConfig) error {
	groups, err := b.db.GetAllDistinctField("groups", "tt", "0", " 100")
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
