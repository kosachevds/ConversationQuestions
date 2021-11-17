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
		processMessage(bot, update.Message, questions)
	}
}

func processMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, questions []string) {
	if message == nil {
		return
	}
	logMessage(message)
	if !message.IsCommand() {
		return
	}

	if message.Command() != "ask" {
		return
	}

	answer := tgbotapi.NewMessage(message.Chat.ID, "")
	if len(questions) == 0 {
		answer.Text = "Answers unavailable"
	} else {
		answer.Text = questions[rand.Intn(len(questions))]
	}
	bot.Send(answer)
}

func logMessage(message *tgbotapi.Message) {
	log.Printf(
		"[%s] %d %s",
		message.From.UserName,
		message.Chat.ID,
		message.Text,
	)
}
