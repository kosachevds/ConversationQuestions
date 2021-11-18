package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var keyboardFlags = map[int64]bool{}

func main() {
	rand.Seed(time.Now().Unix())

	bot, err := initBot()
	if err != nil {
		log.Panic(err)
	}

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

	deletionMessage := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	bot.Request(deletionMessage)

	if message.Command() != "ask" {
		return
	}

	answer := tgbotapi.NewMessage(message.Chat.ID, "")
	if len(questions) == 0 {
		answer.Text = "Answers unavailable"
	} else {
		answer.Text = questions[rand.Intn(len(questions))]
	}

	if _, ok := keyboardFlags[message.Chat.ID]; !ok {
		keyboardFlags[message.Chat.ID] = true
		answer.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/ask")))
	}

	if _, err := bot.Send(answer); err != nil {
		log.Panic(err)
	}
}

func logMessage(message *tgbotapi.Message) {
	log.Printf(
		"[%s] %d %s",
		message.From.UserName,
		message.Chat.ID,
		message.Text,
	)
}

func initBot() (*tgbotapi.BotAPI, error) {
	token := os.Getenv("BOT_ACCESS_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = true
	return bot, nil
}
