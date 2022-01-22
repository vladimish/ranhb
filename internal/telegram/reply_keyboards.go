package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vladimish/ranhb/internal/configurator"
	"github.com/vladimish/ranhb/internal/users"
	"log"
)

func (b *Bot) buildMenuKeyboard(isPremium bool) *tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton
	row1 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Today),
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Tomorrow),
	)
	row2 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.ThisWeek),
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.NextWeek),
	)
	row3 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Teachers),
	)
	row4 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Settings),
	)

	buttons = append(buttons, row1)
	if isPremium {
		buttons = append(buttons, row2)
		buttons = append(buttons, row3)
	}
	buttons = append(buttons, row4)
	keyboard := tgbotapi.NewReplyKeyboard(buttons...)

	return &keyboard
}

func (b *Bot) buildPremiumKeyboard(isPremium bool) *tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton
	row1 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Prem.One),
		tgbotapi.NewKeyboardButton(configurator.Cfg.Prem.Three),
	)
	row2 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Prem.Six),
		tgbotapi.NewKeyboardButton(configurator.Cfg.Prem.Twelve),
	)
	row3 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Ваша подписка активна."),
	)
	row4 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Left),
	)

	if isPremium {
		buttons = append(buttons, row3)
	} else {
		buttons = append(buttons, row1)
		buttons = append(buttons, row2)
	}
	buttons = append(buttons, row4)
	keyboard := tgbotapi.NewReplyKeyboard(buttons...)

	return &keyboard
}

func (b *Bot) sendGroupsKeyboard(chatId int64, user *users.User, pageOffset int) error {
	msg := tgbotapi.NewMessage(chatId, "Выберите вашу группу.")

	user.U.LastActionValue += pageOffset
	err := user.Save()
	if err != nil {
		return err
	}

	keyboard, err := b.buildGroupsKeyboard(b.db.GetAllDistinctField, user, "groups", "tt", fmt.Sprintf("%d", user.U.LastActionValue), fmt.Sprintf("%d", configurator.Cfg.PageSize))
	if err != nil {
		return err
	}

	msg.ReplyMarkup = keyboard
	if user.U.LastActionValue == configurator.Cfg.PageSize {
		msg.Text = "Выберите вашу группу."
	}
	message, err := b.bot.Send(msg)

	if err != nil {
		return err
	}
	log.Println("Message sent: ", message)

	return nil
}

type dataDrainer func(args ...string) ([]string, error)

func (b *Bot) buildGroupsKeyboard(dd dataDrainer, u *users.User, args ...string) (*tgbotapi.ReplyKeyboardMarkup, error) {
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

func (b *Bot) buildSettingsKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton
	row1 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Premium),
	)
	row2 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Info),
	)
	row3 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Left),
	)
	if configurator.Cfg.Premium {
		buttons = append(buttons, row1)
	}
	buttons = append(buttons, row2)
	buttons = append(buttons, row3)
	keyboard := tgbotapi.NewReplyKeyboard(buttons...)

	return &keyboard
}

func (b *Bot) buildTeachersKeyboard(teachers []string) *tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	for i := 0; i < len(teachers); i++ {
		row := tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(teachers[i]),
		)
		rows = append(rows, row)
	}
	back := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Left),
	)
	rows = append(rows, back)
	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	return &keyboard
}

func (b *Bot) sendSettingsKeyboard(user *users.User) error {
	user.U.LastAction = "settings"
	user.U.LastActionValue = 0
	err := user.Save()
	if err != nil {
		return err
	}

	settingsKeyboard := b.buildSettingsKeyboard()
	msg := tgbotapi.NewMessage(user.U.Id, configurator.Cfg.Consts.Settings)
	msg.ReplyMarkup = settingsKeyboard
	m, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	log.Println(m)

	return nil
}
