package store

import "sync"

type Store struct {
	sync.RWMutex
	Data map[string]map[string]int
}

func NewStore() *Store {
	return &Store{
		Data: make(map[string]map[string]int),
	}
}

func (s *Store) AddMessage(sender, receiver string) {
	s.Lock()
	defer s.Unlock()

	senderRange := sender[:len(sender)-5]
	receiverRange := receiver[:len(receiver)-5]

	if s.Data[senderRange] == nil {
		s.Data[senderRange] = make(map[string]int)
	}
	s.Data[senderRange][receiverRange]++
}
