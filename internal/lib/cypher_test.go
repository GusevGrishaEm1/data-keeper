package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCypher(t *testing.T) {
	key := "1234567890123456"
	data := []byte("test")
	encrypted, err := Encrypt(key, data)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := Decrypt(key, encrypted)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, data, decrypted)
}
