package telegram

import (
	"errors"
	"fmt"
	"github.com/vladimish/ranhb/internal/configurator"
	"github.com/vladimish/ranhb/internal/users"
	"github.com/vladimish/yookassa-go-sdk"
	"log"
	"time"
)

func (b *Bot) createPaymentConfig(user *users.User, months int) *yookassa.PaymentConfig {
	var price int
	var description string
	switch months {
	case 1:
		price = configurator.Cfg.Prices.One
		description = configurator.Cfg.Prem.One
	case 3:
		price = configurator.Cfg.Prices.Three
		description = configurator.Cfg.Prem.Three
	case 6:
		price = configurator.Cfg.Prices.Six
		description = configurator.Cfg.Prem.Six
	case 12:
		price = configurator.Cfg.Prices.Twelve
		description = configurator.Cfg.Prem.Twelve
	default:
		price = 100

	}
	amountPrice := yookassa.Amount{
		Value:    fmt.Sprintf("%d", price),
		Currency: "RUB",
	}

	cfg := yookassa.NewPaymentConfig(
		amountPrice,
		yookassa.Redirect{
			Type:      yookassa.TypeRedirect,
			Locale:    "ru_RU",
			Enforce:   false,
			ReturnURL: "https://t.me/ranh_ranepa_bot",
		})

	cfg.Description = fmt.Sprintf("%d", user.U.Id)

	cartItem := yookassa.Item{
		Description:    description,
		Quantity:       "1",
		Amount:         amountPrice,
		VatCode:        1,
		PaymentSubject: "service",
		PaymentMode:    "full_prepayment",
	}
	cart := make([]yookassa.Item, 1)
	cart[0] = cartItem
	cfg.Receipt = &yookassa.Receipt{}
	cfg.Receipt.Items = cart
	return cfg
}

func (b *Bot) CheckPaymentUpdates(requestAmount int, waitTime time.Duration, paymentID string) (*yookassa.Payment, error) {
	log.Println("Check payment entry")
	p, err := b.kassa.GetPayment(paymentID)
	if err != nil {
		return nil, err
	}

	for i := 0; i < requestAmount; i++ {
		p, err = b.kassa.GetPayment(paymentID)
		if err != nil {
			return nil, err
		}
		if p.Status == yookassa.WaitingForCapture {
			err := b.kassa.AcceptSpending(paymentID)
			if err != nil {
				return nil, err
			}
			log.Printf("Payment %s accepted.", paymentID)

			p, err := b.kassa.GetPayment(paymentID)
			if err != nil {
				return nil, err
			}

			return p, nil
		}
		// log.Println("Sleeping")
		time.Sleep(waitTime)
	}

	return nil, errors.New("USER DIDN'T PAYED")
}
