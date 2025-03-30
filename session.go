package main

import "github.com/looplab/fsm"

type Session struct {
	userID int64
	FSM    *fsm.FSM
}

type SessionHandler struct {
	pool Pool
}

func (h SessionHandler) getSession(userID int64) (*Session, error) {
	if s, ok := h.pool.GetFromPool(userID); ok {
		return s, nil
	}

	f := fsm.NewFSM(
		stateWaitGreeting,
		fsm.Events{
			{Name: "greet", Src: []string{stateWaitGreeting}, Dst: stateWaitQuestion},
			{Name: "answer", Src: []string{stateWaitGreeting}, Dst: stateWaitQuestion},
		},
		fsm.Callbacks{},
	)
	s := &Session{
		userID: userID,
		FSM:    f,
	}
	return s, nil
}
