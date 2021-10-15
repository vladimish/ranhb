package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/telf01/ranhb/pkg/configurator"
	"github.com/telf01/ranhb/pkg/db/models"
	"github.com/telf01/ranhb/pkg/users"
	"github.com/telf01/ranhb/pkg/utils/date"
	"github.com/telf01/yookassa-go-sdk"
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

func (b *Bot) buildTeachersMessage(tt []models.TT) (string, error) {
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
			res[dday*100+dmonth] = append(res[dday*100+dmonth], fmt.Sprintf("%s\n<b>%d. %s</b>\n    Группа: %s\n    Тип занятия: %s\n    Аудитория: %s", tt[t].Time, n, tt[t].Subject, tt[t].Groups, tt[t].Subject_type, tt[t].Classroom))
		} else {
			res[dday*100+dmonth] = append(res[dday*100+dmonth], fmt.Sprintf("%s\n<b>%d. %s</b>\n    Группа: %s\n    Аудитория: %s", tt[t].Time, n, tt[t].Subject, tt[t].Groups, tt[t].Classroom))
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
		err := user.Init(message.Chat.ID, b.db)
		if err != nil {
			return err
		}
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
	case "settings":
		err := b.handleSettingsMessage(message, user)
		if err != nil {
			return err
		}
	case "teachers":
		err := b.handleTeachersMessage(message, user)
		if err != nil {
			return err
		}
	case "teachers_selection":
		err := b.handleTeacherSelectionMessage(message, user)
		if err != nil {
			return err
		}
	case "premium":
		err := b.handlePremiumMessage(message, user)
		if err != nil {
			return nil
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
		isPremium, err := b.db.CheckPremiumStatus(user.U.Id)
		if err != nil {
			return err
		}
		if isPremium {
			startTime, endTime := date.GetWeekInterval(time.Now())
			err := b.generateTtCallbackMessage(message, user, startTime, endTime, time.Now())
			if err != nil {
				return err
			}
		}
	case configurator.Cfg.Consts.NextWeek:
		isPremium, err := b.db.CheckPremiumStatus(user.U.Id)
		if err != nil {
			return err
		}
		if isPremium {
			startTime, endTime := date.GetWeekInterval(time.Now().AddDate(0, 0, 7))
			err := b.generateTtCallbackMessage(message, user, startTime, endTime, time.Now().AddDate(0, 0, 7))
			if err != nil {
				return err
			}
		}
	case configurator.Cfg.Consts.Teachers:
		user.U.LastAction = "teachers"
		user.U.LastActionValue = 0
		err := user.Save()
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, "Введите фамилию преподавателя")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(configurator.Cfg.Consts.Left),
			),
		)
		b.bot.Send(msg)
	case configurator.Cfg.Consts.Settings:
		user.U.LastAction = "settings"
		user.U.LastActionValue = 0
		err := user.Save()
		if err != nil {
			return err
		}
		err = b.sendSettingsKeyboard(user)

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

func (b *Bot) handleSettingsMessage(message *tgbotapi.Message, user *users.User) error {
	switch message.Text {
	case configurator.Cfg.Consts.Left:
		err := b.sendMenuKeyboard(user)
		if err != nil {
			return err
		}
	case configurator.Cfg.Consts.Info:
		err := b.sendInfoMessage(user.U.Id)
		if err != nil {
			return err
		}
	case configurator.Cfg.Consts.Premium:
		err := b.sendPremiumKeyboard(user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) sendInfoMessage(id int64) error {
	text := "Ranh bot v1.0\n\nРазработано Владимиром Мишаковым\nПри поддержке Александра Тарасюка\n\nПо всем вопросам и предложением обращайтесь на 01.vladimir.mishakov@gmail.com или @telf01"
	msg := tgbotapi.NewMessage(id, text)
	m, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	log.Println(m)
	return nil
}

func (b *Bot) handleTeachersMessage(message *tgbotapi.Message, user *users.User) error {
	switch message.Text {
	case configurator.Cfg.Consts.Left:
		user.U.LastAction = "menu"
		user.U.LastActionValue = 0
		err := user.Save()
		if err != nil {
			return err
		}
		err = b.sendMenuKeyboard(user)
		if err != nil {
			return err
		}
	default:
		teachers, err := b.db.GetTeachers(message.Text)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, "")
		if len(teachers) > 0 {
			msg.Text = "Выберите преподавателя из списка"
			msg.ReplyMarkup = b.buildTeachersKeyboard(teachers)
			user.U.LastAction = "teachers_selection"
			user.Save()
		} else {
			msg.Text = "Преподаватель не найден"
		}

		b.bot.Send(msg)
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

			err = b.sendMenuKeyboard(user)
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

func (b *Bot) sendMenuKeyboard(user *users.User) error {
	msg := tgbotapi.NewMessage(user.U.Id, "Меню")
	user.U.LastAction = "menu"
	err := user.Save()
	if err != nil {
		return err
	}

	isUserPremium, err := b.db.CheckPremiumStatus(user.U.Id)
	if err != nil {
		return err
	}

	keyboard := b.buildMenuKeyboard(isUserPremium)
	msg.ReplyMarkup = keyboard

	message, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	log.Println("Message sent: ", message)

	return nil
}

func (b *Bot) sendPremiumKeyboard(user *users.User) error {
	if !user.U.IsTermsOfUseAccepted {
		txt := "Совершая покупку, вы соглашаетесь с <a href=\"" + configurator.Cfg.TermsOfUseUrl + "\">условиями предоставления услуг</a>."
		msg := tgbotapi.NewMessage(user.U.Id, txt)
		msg.ParseMode = "HTML"
		answer, err := b.bot.Send(msg)
		if err != nil {
			return err
		}
		user.U.IsTermsOfUseAccepted = true
		err = user.Save()
		if err != nil {
			return err
		}
		log.Println(answer)
	}

	msg := tgbotapi.NewMessage(user.U.Id, "Премиум доступ откроет тебе возможность воспользоваться уникальными функциями нашего Бота. Есть какая-то задолженность? Либо надо встретиться с конкретным преподавателем по какому-то срочному вопросу, и ты не знаешь, где его искать? Не беда! Теперь ты можешь быть в курсе всех событий: видеть своё расписание на неделю и также расписание конкретного преподавателя. В будущем мы добавим ещё огромное количество новых и интересных фишек. Оставайся с нами, и будь всегда на шаг впереди!")
	user.U.LastAction = "premium"
	user.U.LastActionValue = 0
	err := user.Save()
	if err != nil {
		return err
	}

	isUserPremium, err := b.db.CheckPremiumStatus(user.U.Id)
	if err != nil {
		return err
	}

	keyboard := b.buildPremiumKeyboard(isUserPremium)
	msg.ReplyMarkup = keyboard

	message, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	log.Println("Message sent: ", message)

	return nil
}

func (b *Bot) handleTeacherSelectionMessage(message *tgbotapi.Message, user *users.User) error {
	switch message.Text {
	case configurator.Cfg.Consts.Left:
		user.U.LastAction = "menu"
		err := user.Save()
		if err != nil {
			return err
		}

		err = b.sendMenuKeyboard(user)
		if err != nil {
			return err
		}
	default:
		exists, err := b.db.IsTeacherExists(message.Text)
		if err != nil {
			return err
		}
		if exists {
			msg, err := b.buildTeacherTtMessage(message, user, time.Now())
			if err != nil {
				return err
			}

			m, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
			log.Println(m)

			err = b.db.SaveTeacherCallback(message.Chat.ID, m.MessageID, time.Now().Day(), int(time.Now().Month()), message.Text)
			if err != nil {
				return err
			}
		} else {
			user.U.LastAction = "teachers"
			user.Save()
			err := b.handleTeachersMessage(message, user)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Bot) buildTeacherTtMessage(message *tgbotapi.Message, user *users.User, date time.Time) (*tgbotapi.MessageConfig, error) {
	tts, err := b.db.GetTeachersLessons(message.Text, user.U.Group, date.Day(), int(date.Month()))
	if err != nil {
		return nil, err
	}

	var msgString string
	if len(tts) == 0 {
		msgString = message.Text + "\n"
		msgString += "Нет занятий."
	} else {
		msgString += message.Text + "\n"
		msgString, err = b.buildTtMessage(tts)
		if err != nil {
			return nil, err
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, msgString)
	msg.ReplyMarkup = b.buildTeacherKeyboard(date)
	msg.ParseMode = "HTML"

	return &msg, nil
}

func (b *Bot) handlePremiumMessage(message *tgbotapi.Message, user *users.User) error {
	switch message.Text {
	case configurator.Cfg.Consts.Left:
		err := b.sendSettingsKeyboard(user)
		if err != nil {
			return err
		}
	case configurator.Cfg.Prem.One:
		err := b.sendBuyMessage(user, 1)
		if err != nil {
			return err
		}
	case configurator.Cfg.Prem.Three:
		err := b.sendBuyMessage(user, 3)
		if err != nil {
			return err
		}
	case configurator.Cfg.Prem.Six:
		err := b.sendBuyMessage(user, 6)
		if err != nil {
			return err
		}
	case configurator.Cfg.Prem.Twelve:
		err := b.sendBuyMessage(user, 12)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) sendBuyMessage(user *users.User, months int) error {
	pcfg := b.createPaymentConfig(user, months)
	payment, err := b.kassa.SendPaymentConfig(pcfg)
	if err != nil {
		return err
	}

	err = b.db.SavePayment(*payment, user.U.Id)
	if err != nil {
		return err
	}

	url := "Заказа " + payment.Id + " создан. \nОплатите его по ссылке: " + fmt.Sprintf("%v", payment.Confirmation.(map[string]interface{})["confirmation_url"])

	msg := tgbotapi.NewMessage(user.U.Id, url)

	message, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	log.Println("Message sent: ", message)

	updatedPayment, err := b.CheckPaymentUpdates(100, time.Second*10, payment.Id)
	if err != nil {
		log.Println("[ERROR] PAYMENT" + payment.Id + "FAILED")
		errorMessage := tgbotapi.NewMessage(user.U.Id, "Произошла ошибка при оплате заказа "+payment.Id+". В случае, если заказ был оплачен, свяжитесь с нами по адресу, указанному в справочной информации и укажите в письме номер заказа.")
		_, err = b.bot.Send(errorMessage)
		if err != nil {
			log.Println(err)
		}
	}
	if updatedPayment.Status == yookassa.Succeeded {
		log.Printf("User %d bought %s\n", user.U.Id, payment.Description)
		err = b.db.AddPremium(user.U.Id, time.Hour*24*31*time.Duration(months))
		if err != nil {
			log.Println(err)
		}
		successMessage := tgbotapi.NewMessage(user.U.Id, "Платёж прошёл успешно!")
		msg, err := b.bot.Send(successMessage)
		if err != nil{
			return err
		}
		log.Println(msg)
	}

	return nil
}
