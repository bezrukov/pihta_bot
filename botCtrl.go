package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
	"math/rand"
	"fmt"
	"github.com/bezrukov/pihta_bot/modules/dealMock"
	"github.com/bezrukov/pihta_bot/modules/account"
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

var retryKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Попробовать ещё раз", "advice"),
	),
)

var refillKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пополнить", "refill_350"),
	),
)

type User struct {
	Balance *account.Balance
}

var lastAdviceIndex = -1

type botCtrl struct {
	Users map[int64]User
}

func newBotCtrl() *botCtrl {
	return &botCtrl{Users: make(map[int64]User, 10)}
}

func (ctrl *botCtrl) run(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	users := ctrl.Users

	var chatId int64
	if update.Message != nil {
		chatId = update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	}

	user, ok := users[chatId]

	if !ok {
		users[chatId] = User{Balance: account.NewBalance()}
		user = users[chatId]
	}

	balance := user.Balance

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
					tgbotapi.NewInlineKeyboardButtonData("\xF0\x9F\x94\xBC Вверх", "up"),
					tgbotapi.NewInlineKeyboardButtonData("\xF0\x9F\x94\xBD Вниз", "down"),
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
				"\xF0\x9F\x95\x90 Сделка заключена!",
			)
			bot.Send(msg)

			deal := dealMock.NewDeal()
			deal.Start(6)

			for {
				time.Sleep(time.Second * 1)
				reportMsg, isFinish := deal.Process(balance, 100)

				report := tgbotapi.NewMessage(
					update.CallbackQuery.Message.Chat.ID,
					reportMsg,
				)

				if isFinish {
					msg := tgbotapi.NewMessage(
						update.CallbackQuery.Message.Chat.ID,
						reportMsg)
					bot.Send(msg)

					msg = tgbotapi.NewMessage(
						update.CallbackQuery.Message.Chat.ID,
						fmt.Sprintf("\nТвой баланс: %vр.", balance.Current()))
					msg.ReplyMarkup = &retryKeyboard
					bot.Send(msg)
					break
				}

				bot.Send(report)

			}
		case "refill_350":
			refill(350)
		}

		return
	}

	if update.Message == nil {
		return
	}

	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Виталий приветсвтует тебя, " + update.Message.Chat.FirstName + "!")
		msg.ReplyMarkup = &menuKeyboard
		bot.Send(msg)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "На балансе маловато денег. Виталя рекомендует пополниться")
		msg.ReplyMarkup = refillKeyboard
		bot.Send(msg)
	}

	switch update.Message.Text {
	case "Быстрая сделка":
		if balance.Current() < 100 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "На балансе маловато денег. Виталя рекомендует пополниться")
			msg.ReplyMarkup = refillKeyboard
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Поехали!")
		msg.ReplyMarkup = &backKeyboard
		bot.Send(msg)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, getRandomAdvice())
		msg.ReplyMarkup = &vitaliyKeyboard

		bot.Send(msg)
	case "Пополнение":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите кнопку пополнить")
		msg.ReplyMarkup = &refillKeyboard
		bot.Send(msg)
	case "Помощь":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "\xF0\x9F\x9A\x82 Виталя уже выехал к тебе!")
		bot.Send(msg)
	case "Назад":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите пункт")
		msg.ReplyMarkup = menuKeyboard
		bot.Send(msg)
	}
}

func (ctrl *botCtrl) init(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		go ctrl.run(update, bot)
	}
}

func getRandomAdvice() string {
	var advices = []string{
		"\xF0\x9F\x94\xA5 Виталя считает, что стоит попробовать GBPUSD. Не робей, открывай сделку скорей!",
		"\xF0\x9F\x94\xA5 Виталя считает, что ты засиделся без дела. Вот тебе подходящий актив -  Gold. Пора вернуться в торги!",
		"\xF0\x9F\x94\xA5 Виталя заметил, что Bitcoin вырос на 5% за последнее время. Давай попробуем заработать на этом",
		"\xF0\x9F\x94\xA5 Виталя заметил, что Silver упал на 5% за последнее время. Давай попробуем заработать на этом",
		"\xF0\x9F\x94\xA5 Витале кажется, что ты еще не пробовал GBPJPY. Не теряй эту возможность, открывай сделку.",
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
