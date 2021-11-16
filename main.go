package main

import (
	"log"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	token := os.Getenv("FOOBOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	ucfg := tgbotapi.NewUpdate(0)

	ucfg.Timeout = 60
	updatesChan, err := bot.GetUpdatesChan(ucfg)

	for {
		select {
		case update := <-updatesChan:
			processUpdate(bot, update)
		}
	}
}

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userName := update.Message.From.UserName
	chatId := update.Message.Chat.ID
	text := update.Message.Text
	log.Printf("[%s] %d %s", userName, chatId, text)
	msg := tgbotapi.NewMessage(chatId, text)
	bot.Send(msg)
}
