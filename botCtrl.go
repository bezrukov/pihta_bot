package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"pihta_bot/modules/dealMock"
	"pihta_bot/modules/account"
	"time"
	"math/rand"
	"fmt"
)

var menuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Быстрая сделка"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Пополнение"),
		tgbotapi.NewKeyboardButton("Помощь"),
	),
)

var backKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Назад"),
	),
)

var vitaliyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Торговать", "deal"),
		tgbotapi.NewInlineKeyboardButtonData("Другой совет", "advice"),
	),
)

var refillKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("100р", "refill_100"),
		tgbotapi.NewInlineKeyboardButtonData("300р", "refill_300"),
		tgbotapi.NewInlineKeyboardButtonData("500р", "refill_500"),
	),
)

type User struct {
	TelegramId  int
	PlatformId  int64
	Name        string
	PhoneNumber string
}

var lastAdviceIndex = -1

type botCtrl struct {
	Users []User
}

func newBotCtrl() *botCtrl {
	return &botCtrl{}
}

func (ctrl *botCtrl) init(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	balance := account.NewBalance()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		refill := func(value int){
			balance.Inc(value)
			msg := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				fmt.Sprintf("Успешное пополнение. Ваш баланс %vр", balance.Current()),
			)
			bot.Send(msg)

			msg1 := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Поехали!")
			msg1.ReplyMarkup = &backKeyboard
			bot.Send(msg1)

			msg2 := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, getRandomAdvice())
			msg2.ReplyMarkup = &vitaliyKeyboard

			bot.Send(msg2)
		}

		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "deal":
				msg := tgbotapi.NewEditMessageText(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					"Куда пойдёт актив через 1 минуту?")
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Вверх", "up"),
						tgbotapi.NewInlineKeyboardButtonData("Вниз", "down"),
					),
				)
				msg.ReplyMarkup = &keyboard

				bot.Send(msg)
			case "advice":
				msg := tgbotapi.NewEditMessageText(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					getRandomAdvice())
				msg.ReplyMarkup = &vitaliyKeyboard

				bot.Send(msg)
			case "up":
				fallthrough
			case "down":
				msg := tgbotapi.NewEditMessageText(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					"Сделка заключена!",
				)
				bot.Send(msg)

				deal := dealMock.NewDeal()
				deal.Start(5)

				for {
					time.Sleep(time.Second * 1)
					reportMsg, isFinish := deal.Process(balance, 40)

					report := tgbotapi.NewMessage(
						update.CallbackQuery.Message.Chat.ID,
						reportMsg,
					)
					bot.Send(report)

					if isFinish {
						balanceMsg := tgbotapi.NewMessage(
							update.CallbackQuery.Message.Chat.ID,
							fmt.Sprintf("Твой баланс: %vр.", balance.Current()),
						)
						bot.Send(balanceMsg)
						msg := tgbotapi.NewMessage(
							update.CallbackQuery.Message.Chat.ID,
							getRandomAdvice())
						msg.ReplyMarkup = &vitaliyKeyboard

						bot.Send(msg)
						break
					}
				}
			case "refill_100":
				refill(100)
			case "refill_300":
				refill(300)
			case "refill_500":
				refill(500)
			}

			continue
		}

		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, Гришка!")
			msg.ReplyMarkup = &menuKeyboard
			bot.Send(msg)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "На балансе маловато денег. Рекомендую пополниться")
			msg.ReplyMarkup = refillKeyboard
			bot.Send(msg)
		}

		switch update.Message.Text {
		case "Быстрая сделка":
			if balance.Current() < 40 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "На балансе маловато денег. Рекомендую пополниться")
				msg.ReplyMarkup = refillKeyboard
				bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Поехали!")
			msg.ReplyMarkup = &backKeyboard
			bot.Send(msg)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, getRandomAdvice())
			msg.ReplyMarkup = &vitaliyKeyboard

			bot.Send(msg)
		case "Пополнение":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите сумму пополнения")
			msg.ReplyMarkup = &refillKeyboard
			bot.Send(msg)
		case "Помощь":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Помоги себе сам")
			bot.Send(msg)
		case "Назад":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите пункт")
			msg.ReplyMarkup = menuKeyboard
			bot.Send(msg)
		}
	}
}

func getRandomAdvice() string {
	var advices = []string{
		"Виталя считает, что стоит попробовать GBPUSD. Не робей, открывай сделку скорей!",
		"Виталя считает, что ты засиделся без дела. Вот тебе подходящий актив -  Gold. Пора вернуться в торги!",
		"Виталя заметил, что Bitcoin вырос на 5% за последнее время. Давай попробуем заработать на этом",
		"Виталя заметил, что Silver упал на 5% за последнее время. Давай попробуем заработать на этом",
		"Виталя кажется, что ты еще не пробовал GBPJPY. Не теряй эту возможность, открывай сделку.",
	}

	var isEqual = true
	var index = -1

	for isEqual {
		index = rand.Intn(len(advices))
		if index != lastAdviceIndex {
			lastAdviceIndex = index
			isEqual = false
		}
	}

	return advices[index]
}
