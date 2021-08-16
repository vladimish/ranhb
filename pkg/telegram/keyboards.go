package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/configurator"
	"github.com/telf01/ranhb/pkg/users"
	"log"
	"time"
)

func (b *Bot) buildDayKeyboard(date time.Time) (*tgbotapi.InlineKeyboardMarkup, error) {
	yesterday := fmt.Sprintf("%02d.%02d", date.Add(-24*time.Hour).Day(), date.Add(-24*time.Hour).Month())
	tomorrow := fmt.Sprintf("%02d.%02d", date.Add(24*time.Hour).Day(), date.Add(24*time.Hour).Month())

	oldYesterday := fmt.Sprintf("%02d.%02d", date.AddDate(0, 0, -7).Day(), date.Add(-24*time.Hour).Month())
	newTomorrow := fmt.Sprintf("%02d.%02d", date.AddDate(0, 0, 7).Day(), date.Add(-24*time.Hour).Month())

	ttKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(configurator.Cfg.Consts.Left+" "+yesterday, "day/-1"),
			tgbotapi.NewInlineKeyboardButtonData(tomorrow+" "+configurator.Cfg.Consts.Right, "day/1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(configurator.Cfg.Consts.VeryLeft+" "+oldYesterday, "day/-7"),
			tgbotapi.NewInlineKeyboardButtonData(newTomorrow+" "+configurator.Cfg.Consts.VeryRight, "day/7"),
		),
	)

	return &ttKeyboard, nil
}

func (b *Bot) buildWeekKeyboard(date time.Time) (*tgbotapi.InlineKeyboardMarkup, error) {
	oldYesterday := fmt.Sprintf("%02d.%02d", date.AddDate(0, 0, -7).Day(), date.AddDate(0, 0, -7).Month())
	newTomorrow := fmt.Sprintf("%02d.%02d", date.AddDate(0, 0, 7).Day(), date.AddDate(0, 0, 7).Month())

	ttKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(configurator.Cfg.Consts.VeryLeft+" "+oldYesterday, "week/-7"),
			tgbotapi.NewInlineKeyboardButtonData(newTomorrow+" "+configurator.Cfg.Consts.VeryRight, "week/7"),
		),
	)

	return &ttKeyboard, nil
}

func (b *Bot) generateMenuKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton
	row1 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Today),
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Tomorrow),
	)
	row2 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.ThisWeek),
		tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.NextWeek),
	)
	buttons = append(buttons, row1)
	buttons = append(buttons, row2)
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

	keyboard, err := b.generateGroupsKeyboard(b.db.GetAllDistinctField, user, "groups", "tt", fmt.Sprintf("%d", user.U.LastActionValue), fmt.Sprintf("%d", configurator.Cfg.PageSize))
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
