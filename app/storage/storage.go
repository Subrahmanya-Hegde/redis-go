package storage

import "time"

type Value struct {
	Data   string
	Expiry time.Time
}
type Store struct {
	data map[string]Value
}

func NewStore() *Store {
	return &Store{data: make(map[string]Value)}
}

func (s *Store) Get(key string) (Value, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *Store) Set(key string, value Value) {
	s.data[key] = value
}
