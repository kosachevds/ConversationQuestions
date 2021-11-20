package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var keyboardFlags = map[int64]bool{}

func main() {
	rand.Seed(time.Now().Unix())

	bot, updates, err := initBot()
	if err != nil {
		log.Panic(err)
	}

	questionsGist := os.Getenv("QUESTIONS_GIST_ID")
	accessToken := os.Getenv("QUESTIONS_GIST_TOKEN")
	questions, err := downloadQuestions(questionsGist, accessToken)
	if err != nil {
		log.Printf("Questions loading error: %s", err)
	}

	for update := range updates {
		processMessage(bot, update.Message, questions)
	}
}

func processMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, questions []string) {
	if message == nil {
		return
	}
	logMessage(message)

	deleteMessage(bot, message)

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

func initBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	token := os.Getenv("BOT_ACCESS_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Debug = true

	var updates tgbotapi.UpdatesChannel = nil
	heroku_url := os.Getenv("HEROKU_URL")
	if len(heroku_url) > 0 {
		updates, err = initWebhookUpdatesChan(bot, heroku_url)
		if err != nil {
			log.Printf("Launch via webhook error: %s", err)
		}
	}
	if updates == nil {
		ucfg := tgbotapi.NewUpdate(0)
		ucfg.Timeout = 60
		updates = bot.GetUpdatesChan(ucfg)
	}

	return bot, updates, nil
}

func deleteMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	config := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	_, err := bot.Request(config)
	if err != nil {
		log.Panic(err)
	}
}

func initWebhookUpdatesChan(bot *tgbotapi.BotAPI, app_url string) (tgbotapi.UpdatesChannel, error) {
	port := os.Getenv("PORT")
	url := app_url + port + bot.Token
	config, err := tgbotapi.NewWebhook(url)
	if err != nil {
		return nil, err
	}
	_, err = bot.Request(config)
	if err != nil {
		return nil, err
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		return nil, err
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServe(":"+port, nil)
	return updates, nil
}
