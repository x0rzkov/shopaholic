package engine

import (
	"github.com/stretchr/testify/assert"
	"os"
	"shopaholic/store"
	"testing"
)

func TestBoltDB_Register(t *testing.T) {
	defer os.Remove(testDb)
	b := prepBolt(t)
	user := store.User{
		ID:   "user1",
		Name: "user name",
	}

	userID, _ := b.Register(user)
	assert.Equal(t, "user1", userID)

	_, err := b.Register(user)
	assert.NotNil(t, err)
}

func TestBoltDB_Retrieve(t *testing.T) {
	defer os.Remove(testDb)
	b := prepBolt(t)
	user := store.User{
		ID:   "user1",
		Name: "user name",
	}

	userID, _ := b.Register(user)

	result, _ := b.Details(userID)
	assert.Equal(t, "user1", result.ID)
	assert.Equal(t, "user name", result.Name)
}
