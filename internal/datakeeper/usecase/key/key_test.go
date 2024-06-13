package key

import (
	"testing"

	"github.com/stretchr/testify/assert"

	customerr "github.com/GusevGrishaEm1/data-keeper/internal/datakeeper/error"
)

func TestGetKeyForUser(t *testing.T) {
	service := NewKeyService()

	err := service.SetKeyForUser("user123", "some_key")
	assert.NoError(t, err)

	userKey, err := service.GetKeyForUser("user123")
	assert.NoError(t, err)
	assert.Equal(t, "some_key", userKey)

	_, err = service.GetKeyForUser("nonexistent_user")
	assert.Error(t, err)
	assert.Equal(t, customerr.NO_KEY_IN_CONTEXT, err.Error())
}

func TestSetKeyForUser(t *testing.T) {
	service := NewKeyService()

	err := service.SetKeyForUser("user123", "some_key")
	assert.NoError(t, err)

	userKey, ok := service.keys["user123"]
	assert.True(t, ok)
	assert.Equal(t, "some_key", userKey)
}
