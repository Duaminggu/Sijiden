package session

import (
	"sync"
)

type SessionStore struct {
	data map[string]int    // session_id → user_id
	csrf map[string]string // session_id → csrf_token
	mu   sync.RWMutex
}

func NewStore() *SessionStore {
	return &SessionStore{
		data: make(map[string]int),
		csrf: make(map[string]string), // ← ini wajib ditambah
	}
}

func (s *SessionStore) Set(sessionID string, userID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[sessionID] = userID
}

func (s *SessionStore) Get(sessionID string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userID, ok := s.data[sessionID]
	return userID, ok
}

func (s *SessionStore) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, sessionID)
}

func (s *SessionStore) SetCSRF(sessionID, token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.csrf[sessionID] = token
}

func (s *SessionStore) GetCSRF(sessionID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	token, ok := s.csrf[sessionID]
	return token, ok
}
