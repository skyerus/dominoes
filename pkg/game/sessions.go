package game

import (
	"sync"

	"github.com/skyerus/dominoes/pkg/customerror"
)

// Sessions - store of game sessions
type Sessions struct {
	mux      sync.Mutex
	sessions map[string]*Session
}

// NewSessions - constructor
func NewSessions() *Sessions {
	return &Sessions{sessions: make(map[string]*Session)}
}

// NewSession - add new game session to store
func (s *Sessions) NewSession(id string, numOfPlayers int) (*Session, customerror.Error) {
	session, customErr := NewSession(numOfPlayers)
	if customErr != nil {
		return nil, customErr
	}
	s.mux.Lock()
	s.sessions[id] = session
	s.mux.Unlock()

	return session, nil
}

// DeleteSession - delete game session from store
func (s *Sessions) DeleteSession(id string) {
	s.mux.Lock()
	delete(s.sessions, id)
	s.mux.Unlock()
}
