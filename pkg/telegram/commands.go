package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/users"
	"strings"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда")
	switch message.Command() {
	case "start":
		err := b.handleStartCommand(message.Chat.ID)
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

func (b *Bot) handleStartCommand(id int64) error {
	user, err := users.Get(id, b.db)
	if err != nil {
		return err
	}

	if !user.U.IsPrivacyAccepted {
		err := b.sendPrivacyNote(id)
		if err != nil {
			return err
		}

		return nil
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

	return b.sendGroupsKeyboard(id, user, 0)
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
