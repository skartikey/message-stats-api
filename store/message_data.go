package store

import (
	"fmt"
	"sync"
)

type Store struct {
	sync.RWMutex
	sendReceRangeMap map[string]map[string]int
}

func NewStore() *Store {
	return &Store{
		sendReceRangeMap: make(map[string]map[string]int),
	}
}

func (s *Store) AddMessage(sender, receiver string) {
	s.Lock()
	defer s.Unlock()

	senderRange := sender                       // save full sender
	receiverRange := receiver[:len(receiver)-5] // save prefix, trimming last 5 char

	if s.sendReceRangeMap[senderRange] == nil {
		s.sendReceRangeMap[senderRange] = make(map[string]int)
	}
	s.sendReceRangeMap[senderRange][receiverRange]++
}

func (s *Store) PrintMessageCountBySenderAndRange() {
	s.RLock()
	defer s.RUnlock()

	for sender, receiverMap := range s.sendReceRangeMap {
		for receiver, count := range receiverMap {
			fmt.Printf("Sender range: %s, Receiver range: %s, Message count: %d\n", sender, receiver, count)
		}
	}
}
