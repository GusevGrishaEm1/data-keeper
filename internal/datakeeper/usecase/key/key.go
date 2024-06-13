package key

import (
	"math/rand"
	"sync"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/error"
)

type keyService struct {
	keys map[string]string
	m    sync.Mutex
}

func NewKeyService() *keyService {
	rand.NewSource(34)
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (s *keyService) GenerateKey() (string, error) {
	const length = 32
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b), nil
}
