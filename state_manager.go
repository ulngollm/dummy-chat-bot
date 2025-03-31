package main

import (
	"context"
	"fmt"
	"time"

	tele "gopkg.in/telebot.v4"
)

type StateManager struct {
	stateHandlers map[string]tele.HandlerFunc
}

func NewStateManager() *StateManager {
	return &StateManager{
		stateHandlers: map[string]tele.HandlerFunc{
			stateClosed:          greetHandler,
			stateGreeted:         askHandler,
			stateGetQuestion:     askFeedbackHandler,
			stateRequestFeedback: getFeedbackHandler,
			stateGetFeedback:     closeChatHandler,
		},
	}
}

func (sm *StateManager) GetHandler(state string) tele.HandlerFunc {
	return sm.stateHandlers[state]
}

func greetHandler(c tele.Context) error {
	session := c.Get("session").(*Session)
	e := eventGreet
	if !session.FSM.Can(e) {
		if err := session.FSM.Event(context.Background(), eventClose); err != nil {
			return err
		}
	}
	if err := session.FSM.Event(context.Background(), e); err != nil {
		return err
	}
	return c.Send("Добрый день! О чем вы хотите узнать?")
}

func askHandler(c tele.Context) error {
	session := c.Get("session").(*Session)
	evt := eventAskQuestion
	if !session.FSM.Can(evt) {
		return nil
	}
	if err := session.FSM.Event(context.Background(), evt); err != nil {
		return err
	}
	return c.Send("Мы передали ваш вопрос")
}

func askFeedbackHandler(c tele.Context) error {
	session, err := sessionHandler.getSession(c.Chat().ID)
	if err != nil {
		return err
	}
	evt := eventAskFeedback
	if !session.FSM.Can(evt) {
		return nil // No transition possible, return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := session.FSM.Event(ctx, evt); err != nil {
		return err
	}
	go func() {
		select {
		case <-time.After(time.Minute):
			if err := session.FSM.Event(context.Background(), eventClose); err != nil {
				// Log the error if needed
			}
			fmt.Println("чат закрыт после ожидания 1 мин")
		}
	}()
	return c.Send("Напишите нам оценку пжаста")
}

func getFeedbackHandler(c tele.Context) error {
	session, err := sessionHandler.getSession(c.Chat().ID)
	if err != nil {
		return err
	}
	evt := eventCollectFeedback
	if !session.FSM.Can(evt) {
		return nil // No transition possible, return
	}
	if err := session.FSM.Event(context.Background(), evt); err != nil {
		return err
	}
	return c.Send("Спасибо за обратную связь!")
}

func closeChatHandler(c tele.Context) error {
	session := c.Get("session").(*Session)
	e := eventClose
	if !session.FSM.Can(e) {
		return nil
	}
	if err := session.FSM.Event(context.Background(), e); err != nil {
		return err
	}
	return c.Send("Чат закрыт. Вы можете задать свой вопрос снова")
}

func (sm *StateManager) GetHandlerForCurrentState(currentState string) (tele.HandlerFunc, error) {
	return sm.GetHandler(currentState), nil
}

func (s *Session) getCurrentState() string {
	return s.FSM.Current()
}
