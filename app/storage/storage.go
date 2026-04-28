package storage

import "time"

const (
	TypeString = "string"
	TypeList   = "list"
)

type Data struct {
	Type   string
	String string
	List   []string
	Expiry time.Time
}
type Store struct {
	data map[string]Data
}

func NewStore() *Store {
	return &Store{data: make(map[string]Data)}
}

func (s *Store) Get(key string) (Data, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *Store) Set(key string, value Data) {
	s.data[key] = value
}
