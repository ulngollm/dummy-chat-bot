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

	sessionHandler = NewSessionHandler()

	bot.Handle(tele.OnText, greet, closeChat, askFeedback, ask)

	bot.Start()
}

func greet(c tele.Context) error {
	session, err := sessionHandler.getSession(c.Chat().ID)
	if err != nil {
		return err
	}
	evt := eventGreet
	if !session.FSM.Can(evt) {
		if err := session.FSM.Event(context.Background(), eventClose); err != nil {
			return err
		}
	}
	if err := session.FSM.Event(context.Background(), evt); err != nil {
		return err
	}
	return c.Send("Добрый день! О чем вы хотите узнать?")
}

func ask(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		session, err := sessionHandler.getSession(c.Chat().ID)
		if err != nil {
			return err
		}
		evt := eventAskQuestion
		if !session.FSM.Can(evt) {
			return next(c)
		}
		if err := session.FSM.Event(context.Background(), evt); err != nil {
			return err
		}
		return c.Send("Мы передали ваш вопрос")
	}
}

func askFeedback(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		session, err := sessionHandler.getSession(c.Chat().ID)
		if err != nil {
			return err
		}
		evt := eventAskFeedback
		if !session.FSM.Can(evt) {
			return next(c)
		}
		//todo use context with cancel
		if err := session.FSM.Event(context.Background(), evt); err != nil {
			return err
		}
		return c.Send("Напишите нам оценку пжаста")
	}
}

func closeChat(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		session, err := sessionHandler.getSession(c.Chat().ID)
		if err != nil {
			return err
		}
		if !session.FSM.Can(eventClose) {
			return next(c)
		}
		if err := session.FSM.Event(context.Background(), eventClose); err != nil {
			return err
		}

		return c.Send("Чат закрыт. Вы можете задать свой вопрос снова")
	}
}
