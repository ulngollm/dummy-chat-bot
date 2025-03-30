package main

type Pool struct {
	data map[int64]*Session
}

func NewPool() Pool {
	return Pool{
		data: make(map[int64]*Session),
	}
}

func (p *Pool) AddToPool(session *Session) {
	p.data[session.userID] = session
}

func (p *Pool) GetFromPool(id int64) (*Session, bool) {
	v, ok := p.data[id]
	if ok {
		return v, true
	}
	return &Session{}, false
}
