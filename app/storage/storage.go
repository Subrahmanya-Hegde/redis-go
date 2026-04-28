package storage

import "time"

const (
	TypeString = "string"
	TypeList   = "list"
)

type Value struct {
	Type   string
	String string
	List   []string
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
