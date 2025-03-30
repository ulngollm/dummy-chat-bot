package main

import (
	"context"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	"log"
	"os"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv.Load: %s", err)
		return
	}
}

const (
	stateWaitGreeting string = "greeting"
	stateWaitQuestion string = "question"
	stateClosed       string = "closed"
	stateWaitFeedback string = "feedback"
)

const (
	eventGreet       string = "greet"
	eventAskQuestion string = "answer"
	eventClose       string = "closeChat"
	eventAskFeedback string = "feedback"
)

var sessionHandler SessionHandler

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

	bot.Handle(tele.OnText, greet)
	bot.Handle(tele.OnText, ask)
	bot.Handle(tele.OnText, closeChat)
	bot.Handle(tele.OnText, askFeedback)

	bot.Start()
}

func greet(c tele.Context) error {
	session, err := sessionHandler.getSession(c.Chat().ID)
	if err != nil {
		return err
	}
	if !session.FSM.Can(eventGreet) {
		return nil
	}
	if err := session.FSM.Event(context.Background(), eventGreet); err != nil {
		return err
	}
	return c.Send("Добрый день! О чем вы хотите узнать?")
}

func ask(c tele.Context) error {
	session, err := sessionHandler.getSession(c.Chat().ID)
	if err != nil {
		return err
	}
	if !session.FSM.Can(eventAskQuestion) {
		return nil
	}
	if err := session.FSM.Event(context.Background(), eventAskFeedback); err != nil {
		return err
	}
	return c.Send("Мы передали ваш вопрос")
}

func askFeedback(c tele.Context) error {
	session, err := sessionHandler.getSession(c.Chat().ID)
	if err != nil {
		return err
	}
	if !session.FSM.Can("ask") {
		return nil
	}
	if err := session.FSM.Event(context.Background(), eventAskFeedback); err != nil {
		return err
	}
	return c.Send("Напишите нам оценку пжаста")
}

func closeChat(c tele.Context) error {
	session, err := sessionHandler.getSession(c.Chat().ID)
	if err != nil {
		return err
	}
	if !session.FSM.Can("ask") {
		return nil
	}
	if err := session.FSM.Event(context.Background(), eventClose); err != nil {
		return err
	}

	return c.Send("Чат закрыт. Вы можете задать свой вопрос снова")
}
