package session

import (
	"sync"
)

type SessionStore struct {
	data map[string]int               // session_id → user_id
	csrf map[string]string            // session_id → csrf_token
	kv   map[string]map[string]string // session_id → { key: value }
	mu   sync.RWMutex
}

func NewStore() *SessionStore {
	return &SessionStore{
		data: make(map[string]int),
		csrf: make(map[string]string),
		kv:   make(map[string]map[string]string),
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
	delete(s.csrf, sessionID)
	delete(s.kv, sessionID)
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

func (s *SessionStore) ValidateCSRF(sessionID, token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expectedToken, ok := s.csrf[sessionID]
	if !ok || expectedToken == "" {
		return false
	}

	return expectedToken == token
}

func (s *SessionStore) SetValue(sessionID, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.kv[sessionID]; !ok {
		s.kv[sessionID] = make(map[string]string)
	}
	s.kv[sessionID][key] = value
}

func (s *SessionStore) GetValue(sessionID, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if session, ok := s.kv[sessionID]; ok {
		val, exists := session[key]
		return val, exists
	}
	return "", false
}
