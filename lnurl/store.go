package lnurl

type InMemoryStore struct {
	db map[string]string
}

func (s *InMemoryStore) Get(key string) string {
	return s.db[key]
}

func (s *InMemoryStore) Set(key, value string) {
	s.db[key] = value
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{db: map[string]string{}}
}
