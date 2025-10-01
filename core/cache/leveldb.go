package cache

import (
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBCache struct {
	db *leveldb.DB
}

func init() {
	Register("leveldb", func() (Cache, error) {
		return NewLevelDBCache()
	})
}

func NewLevelDBCache() (*LevelDBCache, error) {
	path := os.Getenv("LEVELDB_PATH")
	if path == "" {
		path = "./leveldb"
	}
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBCache{db: db}, nil
}

func (l *LevelDBCache) Set(key string, value interface{}) error {
	return l.db.Put([]byte(key), []byte(fmt.Sprint(value)), nil)
}

func (l *LevelDBCache) Get(key string) (interface{}, error) {
	val, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	return string(val), nil
}

func (l *LevelDBCache) Delete(key string) error {
	return l.db.Delete([]byte(key), nil)
}
