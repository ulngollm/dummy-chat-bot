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

	//todo вынести работу со стейтами в StateManager
	f := fsm.NewFSM(
		stateClosed,
		fsm.Events{
			{Name: eventGreet, Src: []string{stateClosed}, Dst: stateGreeted},
			{Name: eventAskQuestion, Src: []string{stateGreeted}, Dst: stateGetQuestion},
			{Name: eventAskFeedback, Src: []string{stateGetQuestion}, Dst: stateRequestFeedback},
			{Name: eventCollectFeedback, Src: []string{stateRequestFeedback}, Dst: stateClosed},
			{Name: eventClose, Src: []string{stateRequestFeedback, stateGetFeedback}, Dst: stateClosed},
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
