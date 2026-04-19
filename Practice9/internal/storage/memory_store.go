package storage

import "sync"

type Record struct {
	Status     string
	StatusCode int
	Body       []byte
}

type Store struct {
	mu    sync.Mutex
	store map[string]*Record
}

func NewStore() *Store {
	return &Store{
		store: make(map[string]*Record),
	}
}

// Atomic: get OR create processing
func (s *Store) GetOrCreate(key string) (*Record, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if rec, ok := s.store[key]; ok {
		return rec, true
	}

	rec := &Record{Status: "processing"}
	s.store[key] = rec
	return rec, false
}

func (s *Store) SetCompleted(key string, statusCode int, body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[key] = &Record{
		Status:     "completed",
		StatusCode: statusCode,
		Body:       body,
	}
}
