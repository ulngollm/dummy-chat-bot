package main

import (
	"context"
	"time"

	tele "gopkg.in/telebot.v4"
)

// StateManager manages the mapping of states to handlers.
type StateManager struct {
	stateHandlers map[string]tele.HandlerFunc
}

// NewStateManager initializes a new StateManager with state-handler mappings.
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

// GetHandler returns the handler for a given state.
func (sm *StateManager) GetHandler(state string) tele.HandlerFunc {
	return sm.stateHandlers[state]
}

// Handlers for each state

func greetHandler(c tele.Context) error {
	//todo remove. retrieve from c.Get("session")
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

func askHandler(c tele.Context) error {
	session, err := sessionHandler.getSession(c.Chat().ID)
	if err != nil {
		return err
	}
	evt := eventAskQuestion
	if !session.FSM.Can(evt) {
		return nil // No transition possible, return
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
	select {
	case <-ctx.Done():
		if err := session.FSM.Event(context.Background(), eventClose); err != nil {
			return err
		}
	}
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
	session, err := sessionHandler.getSession(c.Chat().ID)
	if err != nil {
		return err
	}
	if !session.FSM.Can(eventClose) {
		return nil // No transition possible, return
	}
	if err := session.FSM.Event(context.Background(), eventClose); err != nil {
		return err
	}

	return c.Send("Чат закрыт. Вы можете задать свой вопрос снова")
}

// GetHandlerForCurrentState returns the handler for the current state of the session.
func (sm *StateManager) GetHandlerForCurrentState(chatID int64) (tele.HandlerFunc, error) {
	session, err := sessionHandler.getSession(chatID)
	if err != nil {
		return nil, err
	}
	currentState := session.getCurrentState()
	return sm.GetHandler(currentState), nil
}

// Define getCurrentState method for Session

// getCurrentState returns the current state of the session.
func (s *Session) getCurrentState() string {
	return s.FSM.Current()
}
