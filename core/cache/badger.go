package cache

import (
	"erp6-be-golang/core/configs"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

type BadgerCache struct {
	db *badger.DB
}

func init() {
	Register("badger", func() (Cache, error) {
		return NewBadgerCache()
	})
}

func NewBadgerCache() (*BadgerCache, error) {
	path := configs.ConfigApps.CacheAddr
	if path == "" {
		path = "./badger"
	}
	opts := badger.DefaultOptions(path).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerCache{db: db}, nil
}

func (b *BadgerCache) Set(key string, value interface{}) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(fmt.Sprint(value)))
	})
}

func (b *BadgerCache) Get(key string) (interface{}, error) {
	var val []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(v []byte) error {
			val = append([]byte{}, v...)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return string(val), nil
}

func (b *BadgerCache) Delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}
