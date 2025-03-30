package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv.Load: %s", err)
		return
	}
}

const (
	stateGreeted         string = "greeting"
	stateGetQuestion     string = "question"
	stateClosed          string = "closed"
	stateRequestFeedback string = "req_feedback"
	stateGetFeedback     string = "get_feedback"
)

const (
	eventGreet           string = "greet"
	eventAskQuestion     string = "answer"
	eventClose           string = "closeChat"
	eventAskFeedback     string = "feedback"
	eventCollectFeedback string = "collect"
)

var sessionHandler SessionHandler
var stateManager *StateManager

func main() {
	t, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatalf("bot token is empty")
		return
	}

	pref := tele.Settings{
		Token:  t,
		Poller: &tele.LongPoller{Timeout: 2 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %s", err)
		return
	}

	sessionHandler = NewSessionHandler()
	stateManager = NewStateManager()

	bot.Handle(tele.OnText, handle)

	bot.Start()
}

func handle(c tele.Context) error {
	handler, err := stateManager.GetHandlerForCurrentState(c.Chat().ID)
	if err != nil {
		return err
	}
	return handler(c)
}
