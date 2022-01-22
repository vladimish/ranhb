package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vladimish/ranhb/pkg/configurator"
	"time"
)

func (b *Bot) buildDayKeyboard(date time.Time) *tgbotapi.InlineKeyboardMarkup {
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

	return &ttKeyboard
}

func (b *Bot) buildWeekKeyboard(date time.Time) *tgbotapi.InlineKeyboardMarkup {
	oldYesterday := fmt.Sprintf("%02d.%02d", date.AddDate(0, 0, -7).Day(), date.AddDate(0, 0, -7).Month())
	newTomorrow := fmt.Sprintf("%02d.%02d", date.AddDate(0, 0, 7).Day(), date.AddDate(0, 0, 7).Month())

	ttKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(configurator.Cfg.Consts.VeryLeft+" "+oldYesterday, "week/-7"),
			tgbotapi.NewInlineKeyboardButtonData(newTomorrow+" "+configurator.Cfg.Consts.VeryRight, "week/7"),
		),
	)

	return &ttKeyboard
}

func (b *Bot) generatePrivacyKeyboard() *tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Продолжить", acceptString),
	)
	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отказаться", declineString),
	)

	ttKeyboard := tgbotapi.NewInlineKeyboardMarkup(row1, row2)

	return &ttKeyboard
}

func (b *Bot) buildTeacherKeyboard(date time.Time) *tgbotapi.InlineKeyboardMarkup {
	yesterday := fmt.Sprintf("%02d.%02d", date.Add(-24*time.Hour).Day(), date.Add(-24*time.Hour).Month())
	tomorrow := fmt.Sprintf("%02d.%02d", date.Add(24*time.Hour).Day(), date.Add(24*time.Hour).Month())

	ttKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(configurator.Cfg.Consts.Left+" "+yesterday, "teacher/day/-1"),
			tgbotapi.NewInlineKeyboardButtonData(tomorrow+" "+configurator.Cfg.Consts.Right, "teacher/day/1"),
		),
	)

	return &ttKeyboard
}
