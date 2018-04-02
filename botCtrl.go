package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Быстрая сделка"),
		tgbotapi.NewKeyboardButton("Пополнение"),
		tgbotapi.NewKeyboardButton("Помощь"),
	),
)

type botCtrl struct {
}

func newBotCtrl() *botCtrl {
	return &botCtrl{}
}

func (ctrl *botCtrl) init(token string, userIds *map[int64]User) {
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

		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "assets":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выбирай вверх или вниз")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Вверх", "Up"),
						tgbotapi.NewInlineKeyboardButtonData("Вниз", "Down"),
					),
				)
				bot.Send(msg)
			case "Up":
				fallthrough
			case "Down":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сделка заключена")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				bot.Send(msg)
			}

			continue
		}

		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите пункт")
			msg.ReplyMarkup = numericKeyboard
			bot.Send(msg)
		}

		switch update.Message.Text {
		case "Быстрая сделка":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, " Виталя выбрал все за тебя: "+
				"сумма сделки 50 рублей, и продолжительность 1 минута")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("USD CLP", "assets"),
					tgbotapi.NewInlineKeyboardButtonData("EUR AUD", "assets"),
					tgbotapi.NewInlineKeyboardButtonData("GBP JPY", "assets"),
				),
			)

			bot.Send(msg)
		case "Пополнение":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вас пополнили")
			bot.Send(msg)
		case "Помощь":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Помоги себе сам")

			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Hi", "USD CLP"),
				),
			)
			bot.Send(msg)
		}
	}
}
