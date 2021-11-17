package main

import (
	"log"
	"math/rand"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("BOT_ACCESS_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	questionsGist := os.Getenv("QUESTIONS_GIST_ID")
	accessToken := os.Getenv("QUESTIONS_GIST_TOKEN")
	questions, err := downloadQuestions(questionsGist, accessToken)
	if err != nil {
		log.Printf("Questions loading error: %s", err)
	}

	ucfg := tgbotapi.NewUpdate(0)

	ucfg.Timeout = 60
	updatesChan := bot.GetUpdatesChan(ucfg)

	for update := range updatesChan {
		if update.Message == nil {
			continue
		}
		processUpdate(bot, update, questions)
	}
}

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, questions []string) {
	userName := update.Message.From.UserName
	chatId := update.Message.Chat.ID
	text := update.Message.Text
	log.Printf("[%s] %d %s", userName, chatId, text)

	var answerMessage string
	if len(questions) == 0 {
		answerMessage = "Answers unavailable"
	} else {
		answerMessage = questions[rand.Intn(len(questions))]
	}
	msg := tgbotapi.NewMessage(chatId, answerMessage)
	bot.Send(msg)
}
