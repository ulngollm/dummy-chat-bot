package main

import "github.com/looplab/fsm"

type Session struct {
	userID int64
	FSM    *fsm.FSM
}

type SessionHandler struct {
	pool Pool
}

func NewSessionHandler() SessionHandler {
	return SessionHandler{pool: NewPool()}
}

func (h SessionHandler) getSession(userID int64) (*Session, error) {
	if s, ok := h.pool.GetFromPool(userID); ok {
		return s, nil
	}

	f := fsm.NewFSM(
		stateClosed,
		fsm.Events{
			{Name: eventGreet, Src: []string{stateClosed}, Dst: stateWaitGreeting},
			{Name: eventAskQuestion, Src: []string{stateWaitGreeting}, Dst: stateWaitQuestion},
			{Name: eventAskFeedback, Src: []string{stateWaitQuestion}, Dst: stateWaitFeedback},
			{Name: eventClose, Src: []string{stateWaitFeedback}, Dst: stateClosed},
		},
		fsm.Callbacks{},
	)
	s := &Session{
		userID: userID,
		FSM:    f,
	}
	h.pool.AddToPool(s)
	return s, nil
}
