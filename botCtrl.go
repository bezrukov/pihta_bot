package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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
		if update.Message == nil {
			continue
		}

		//fmt.Println(update.Message.From.ID)
		//if


		// log.Printf("%#v", update)
		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "nah idi")
			// отправляем
			bot.Send(msg)
		}

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		//
		//bot.Send(msg)
	}
}
