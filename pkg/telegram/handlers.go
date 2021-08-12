package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/configurator"
	"github.com/telf01/ranhb/pkg/db/models"
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
	} else {
		return InvalidCallbackErr
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
		startTime, endTime = b.getWeekInterval(d)
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

	// Answer to callback
	cbcfg := tgbotapi.NewCallback(query.ID, "OK")
	resp, err := b.bot.AnswerCallbackQuery(cbcfg)
	if err != nil {
		return err
	}
	log.Println(resp)

	return nil
}

// buildTtMessage creates string from tt slice.
func (b *Bot) buildTtMessage(tt []models.TT) (string, error) {
	str := ""
	res := make(map[int][]string, 0)
	n := 1
	for t := range tt {
		dday, err := strconv.Atoi(tt[t].Day)
		if err != nil {
			return "", err
		}
		dmonth, err := strconv.Atoi(tt[t].Month)
		if err != nil {
			return "", err
		}

		weekdayNumber, err := strconv.Atoi(tt[t].Day_of_week)
		if err != nil {
			return "", err
		}
		weekday := date.IntToWeekday(weekdayNumber)

		if len(res[dday*100+dmonth]) == 0 {
			res[dday*100+dmonth] = []string{}
			res[dday*100+dmonth] = append(res[dday*100+dmonth], fmt.Sprintf("<u>%s %02d.%02d</u>\n\n", weekday, dday, dmonth))
		}

		if t >= 1 {
			if tt[t-1].Time != tt[t].Time {
				n++
			}
		}

		if tt[t].Subject_type != "" {
			res[dday*100+dmonth] = append(res[dday*100+dmonth], fmt.Sprintf("%s\n<b>%d. %s</b>\n    Преподаватель: %s\n    Тип занятия: %s\n    Аудитория: %s", tt[t].Time, n, tt[t].Subject, tt[t].Teacher, tt[t].Subject_type, tt[t].Classroom))
		} else {
			res[dday*100+dmonth] = append(res[dday*100+dmonth], fmt.Sprintf("%s\n<b>%d. %s</b>\n    Преподаватель: %s\n    Аудитория: %s", tt[t].Time, n, tt[t].Subject, tt[t].Teacher, tt[t].Classroom))
		}
		if tt[t].Subgroup != "" {
			res[dday*100+dmonth] = append(res[dday*100+dmonth], fmt.Sprintf("\n    Подгруппа: %s\n\n", tt[t].Subgroup))
		} else {
			res[dday*100+dmonth] = append(res[dday*100+dmonth], "\n\n")
		}
	}

	for i := range res {
		str += strings.Join(res[i], "")
	}

	return str, nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	// Search for a user in a database.
	user, err := users.Get(message.Chat.ID, b.db)
	if err != nil {
		return err
	}

	// Create new user if it's not registered yet.
	if user.U == nil {
		user.Init(message.Chat.ID, b.db)
		user.U.LastActionValue = -configurator.Cfg.PageSize
	}

	// Choose how to process messages based on the last user action.
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

// handleMenuMessage processes the message sent by the user from the tt view menu.
func (b *Bot) handleMenuMessage(message *tgbotapi.Message, user *users.User) error {
	switch message.Text {
	case configurator.Cfg.Consts.Today:
		t := time.Now()
		err := b.generateTtCallbackMessage(message, user, t, t, t)
		if err != nil {
			return err
		}
	case configurator.Cfg.Consts.Tomorrow:
		t := time.Now().AddDate(0, 0, 1)
		err := b.generateTtCallbackMessage(message, user, t, t, t)
		if err != nil {
			return err
		}
	case configurator.Cfg.Consts.ThisWeek:
		startTime, endTime := b.getWeekInterval(time.Now())
		err := b.generateTtCallbackMessage(message, user, startTime, endTime, time.Now())
		if err != nil {
			return err
		}
	case configurator.Cfg.Consts.NextWeek:
		startTime, endTime := b.getWeekInterval(time.Now().AddDate(0, 0, 7))
		err := b.generateTtCallbackMessage(message, user, startTime, endTime, time.Now().AddDate(0, 0, 7))
		if err != nil {
			return err
		}

	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда.")
		response, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
		log.Println(response)
	}
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

	if fullMsgString == initMsgString{
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

// getWeekInterval returns two times
// which corresponds for monday and sunday
// of a d week.
func (b *Bot) getWeekInterval(d time.Time) (startTime time.Time, endTime time.Time) {
	for d.Weekday() != time.Monday {
		d = d.AddDate(0, 0, -1)
	}
	startTime = d
	endTime = startTime.Add(6 * 24 * time.Hour)

	return startTime, endTime
}

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
