package telegram

import (
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

var acceptString = "accept"
var declineString = "decline"

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
		startTime, endTime := date.GetWeekInterval(time.Now())
		err := b.generateTtCallbackMessage(message, user, startTime, endTime, time.Now())
		if err != nil {
			return err
		}
	case configurator.Cfg.Consts.NextWeek:
		startTime, endTime := date.GetWeekInterval(time.Now().AddDate(0, 0, 7))
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

func (b *Bot) sendPrivacyNote(chatId int64) error {
	msgText := "Для начала работы ознакомьтесь с <a href=\"" + configurator.Cfg.PrivacyUrl + "\">политикой конфиденциальности</a>. Если согласны — нажмите <b>Продолжить</b>"
	msgConfig := tgbotapi.NewMessage(chatId, msgText)
	msgConfig.ParseMode = "HTML"
	msgConfig.ReplyMarkup = b.generatePrivacyKeyboard()

	m, err := b.bot.Send(msgConfig)
	if err != nil {
		return err
	}
	log.Println("message sent: ", m)

	return nil
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
