package key

import (
	"math/rand"
	"sync"
	"time"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/error"
)

type Service struct {
	keys map[string]string
	mu   sync.RWMutex
}

func NewKeyService() *Service {
	rand.NewSource(int64(time.Now().Nanosecond()))
	return &Service{keys: make(map[string]string)}
}

func (s *Service) GetKeyForUser(user string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key, ok := s.keys[user]
	if !ok {
		return "", customerr.Error(customerr.NO_KEY_IN_CONTEXT)
	}
	return key, nil
}

func (s *Service) SetKeyForUser(user string, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys[user] = key
	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (s *Service) GenerateKey() (string, error) {
	const length = 32
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b), nil
}
