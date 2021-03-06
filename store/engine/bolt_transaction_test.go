package engine

import (
	bolt "github.com/coreos/bbolt"
	"github.com/stretchr/testify/assert"
	"os"
	"shopaholic/store"
	"shopaholic/utils"
	"testing"
	"time"
)

var testDb = "test-transaction.db"

func TestBoltDB_CreateAndList(t *testing.T) {
	defer os.Remove(testDb)
	var b = prepBolt(t)

	res, err := b.List(store.User{ID: "user1"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, utils.Money{2121, "usd"}, res[0].Amount)
	assert.Equal(t, "user1", res[0].User.ID)
	t.Log(res[0].ID)

	_, err = b.Create(store.Transaction{ID: res[0].ID, User: store.User{ID: "user1"}})
	assert.NotNil(t, err)
	assert.Equal(t, "key id-1 already in store", err.Error())

	_, err = b.List(store.User{ID: "user-not-found"})
	assert.EqualError(t, err, `no bucket user-not-found in store`)

	assert.NoError(t, b.Disconnect())
}

func TestBoltDB_New(t *testing.T) {
	_, err := NewBoltDB(bolt.Options{}, "/tmp/no-such-place/tmp.db")
	assert.EqualError(t, err, "failed to make boltdb for /tmp/no-such-place/tmp.db: open /tmp/no-such-place/tmp.db: no such file or directory")
}

func prepBolt(t *testing.T) *BoltDB {
	os.Remove(testDb)

	boltStore, err := NewBoltDB(bolt.Options{}, testDb)
	assert.Nil(t, err)
	b := boltStore

	transaction := store.Transaction{
		ID:        "id-1",
		Amount:    utils.Money{2121, "usd"},
		CreatedAt: time.Date(2017, 12, 20, 15, 18, 22, 0, time.Local),
		User:      store.User{ID: "user1", Name: "user name"},
	}
	_, err = b.Create(transaction)
	assert.Nil(t, err)

	transaction2 := store.Transaction{
		ID:        "id-2",
		Amount:    utils.Money{21212, "usd"},
		CreatedAt: time.Date(2017, 12, 20, 15, 18, 22, 0, time.Local),
		User:      store.User{ID: "user2", Name: "second name"},
	}
	_, err = b.Create(transaction2)
	assert.Nil(t, err)

	return b
}
