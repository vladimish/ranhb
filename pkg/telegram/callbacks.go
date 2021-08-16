package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/configurator"
	"github.com/telf01/ranhb/pkg/users"
	"github.com/telf01/ranhb/pkg/utils/date"
	"log"
	"strconv"
	"strings"
	"time"
)

var InvalidCallbackErr = errors.New("INVALID CALLBACK")

func (b *Bot) handleCallback(query *tgbotapi.CallbackQuery) error {
	queryParts := strings.Split(query.Data, "/")
	if len(queryParts) == 2 {
		err := b.handleTtCallback(query, queryParts[0], queryParts[1])
		if err != nil {
			return err
		}
	} else if len(queryParts) == 1 {
		err := b.handlePrivacyCallback(query)
		if err != nil {
			return err
		}
	} else {
		return InvalidCallbackErr
	}

	return nil
}

func (b *Bot) handlePrivacyCallback(query *tgbotapi.CallbackQuery) error {
	user, err := users.Get(query.Message.Chat.ID, b.db)
	if err != nil {
		return err
	}

	err = user.Init(query.Message.Chat.ID, b.db)
	if err != nil {
		return err
	}

	if query.Data == acceptString {
		user.U.IsPrivacyAccepted = true
		user.Save()
		b.handleStartCommand(query.Message.Chat.ID)
	} else {
		user.U.IsPrivacyAccepted = false
		user.Save()
	}

	err = b.answerToCallback(query.ID, "OK")
	if err != nil {
		return err
	}

	return nil
}

// handleTtCallback processing callback
// gotten from tt menu messages
func (b *Bot) handleTtCallback(query *tgbotapi.CallbackQuery, queryType string, queryData string) error {
	// Get integer value of callback data.
	daysToSkip, err := strconv.Atoi(queryData)
	if err != nil {
		return err
	}

	// Find callback in database and get its date.
	day, month, err := b.db.GetCallback(query.Message.Chat.ID, query.Message.MessageID)
	if err != nil {
		return err
	}

	user, err := users.Get(query.Message.Chat.ID, b.db)
	if err != nil {
		return err
	}

	// Get required time interval.
	d := time.Date(time.Now().Year(), time.Month(month), day, 0, 0, 0, 0, time.FixedZone(configurator.Cfg.TimeZone, 0))
	d = d.AddDate(0, 0, daysToSkip)
	var startTime, endTime time.Time
	switch queryType {
	case "week":
		startTime, endTime = date.GetWeekInterval(d)
	case "day":
		startTime = d
		endTime = d
	default:
		return InvalidCallbackErr
	}

	initMsgString := user.U.Group + "\n"
	fullMsgString := initMsgString

	// Build message.
	for i := startTime; i.Unix() <= endTime.Unix(); i = i.AddDate(0, 0, 1) {
		tts, err := b.db.GetSpecificTt(user.U.Group, i.Day(), int(i.Month()))
		if err != nil {
			return err
		}

		msgString, err := b.buildTtMessage(tts)
		if err != nil {
			return err
		}

		fullMsgString += msgString
	}

	if fullMsgString == initMsgString {
		fullMsgString += "Занятий нет."
	}

	// Create and send message.
	nmsg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, fullMsgString)
	nmsg.ParseMode = "HTML"
	m, err := b.bot.Send(nmsg)
	if err != nil {
		return err
	}
	log.Println(m)

	// Add keyboard to edited message.
	var keyboard *tgbotapi.InlineKeyboardMarkup
	switch queryType {
	case "day":
		keyboard, err = b.buildDayKeyboard(d)
		if err != nil {
			return err
		}
	case "week":
		keyboard, err = b.buildWeekKeyboard(d)
		if err != nil {
			return err
		}
	default:
		return InvalidCallbackErr
	}

	nmarkup := tgbotapi.NewEditMessageReplyMarkup(query.Message.Chat.ID, query.Message.MessageID, *keyboard)
	m, err = b.bot.Send(nmarkup)
	if err != nil {
		return err
	}
	log.Println(m)

	err = b.db.UpdateCallback(m.Chat.ID, m.MessageID, d.Day(), int(d.Month()))
	if err != nil {
		return err
	}

	err = b.answerToCallback(query.ID, "OK")
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) answerToCallback(queryID string, msg string) error {
	// Answer to callback
	cbcfg := tgbotapi.NewCallback(queryID, "OK")
	resp, err := b.bot.AnswerCallbackQuery(cbcfg)
	if err != nil {
		return err
	}
	log.Println(resp)

	return nil
}

func (b *Bot) generateTtCallbackMessage(message *tgbotapi.Message, user *users.User, startTime time.Time, endTime time.Time, exactTime time.Time) error {
	initMsgString := user.U.Group + "\n"
	fullMsgString := initMsgString

	// Select which keyboard to use
	// if all times are equal then there
	// is only one day, and we need
	// to use day keyboard
	var keyboard *tgbotapi.InlineKeyboardMarkup
	if startTime == endTime && endTime == exactTime {
		var err error
		keyboard, err = b.buildDayKeyboard(exactTime)
		if err != nil {
			return err
		}
	} else {
		var err error
		keyboard, err = b.buildWeekKeyboard(exactTime)
		if err != nil {
			return err
		}
	}

	// Build tt message for every provided day.
	for i := startTime; i.Unix() <= endTime.Unix(); i = i.AddDate(0, 0, 1) {
		tts, err := b.db.GetSpecificTt(user.U.Group, i.Day(), int(i.Month()))
		if err != nil {
			return err
		}

		msgString, err := b.buildTtMessage(tts)
		if err != nil {
			return err
		}

		fullMsgString += msgString
	}

	if fullMsgString == initMsgString {
		fullMsgString += "Занятий нет."
	}

	// Make message config.
	msg := tgbotapi.NewMessage(message.Chat.ID, fullMsgString)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "HTML"

	m, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	// Add callback to database for further tracking.
	err = b.db.SaveCallback(message.Chat.ID, m.MessageID, exactTime.Day(), int(exactTime.Month()))
	if err != nil {
		return err
	}

	log.Println(m)

	return nil
}