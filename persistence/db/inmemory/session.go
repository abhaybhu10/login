package db

import (
	"errors"

	"github.com/abhaybhu10/login/model"
)

type InMemorySession struct {
	sessions map[string]model.Session
}

func NewInMomorySession() *InMemorySession {
	return &InMemorySession{
		sessions: map[string]model.Session{},
	}
}

func (s *InMemorySession) Save(session model.Session) error {
	s.sessions[session.SessionID] = session
	return nil
}

func (s *InMemorySession) Get(sessionId string) (*model.Session, error) {
	session, ok := s.sessions[sessionId]
	if !ok {
		return nil, errors.New("Session not found")
	}
	return &session, nil
}
