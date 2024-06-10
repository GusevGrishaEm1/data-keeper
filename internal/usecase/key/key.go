package key

import (
	"crypto/rand"
	"sync"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/error"
)

type keyService struct {
	keys map[string]string
	m    sync.Mutex
}

func NewKeyService() *keyService {
	return &keyService{
		keys: make(map[string]string),
	}
}

func (s *keyService) GetKeyForUser(user string) (string, error) {
	key, ok := s.keys[user]
	if !ok {
		return "", customerr.Error(customerr.NO_KEY_IN_CONTEXT)
	}
	return key, nil
}

func (s *keyService) SetKeyForUser(user string, key string) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.keys[user] = key
	return nil
}

func (s *keyService) GenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return string(key), nil
}
