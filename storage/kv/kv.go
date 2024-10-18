package kv

import (
	"errors"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore"
)

const (
	prefixLen = 8
)

type KV struct {
	db *badger.DB
}

func New(dbPath string) (KV, error) {
	options := badger.DefaultOptions(dbPath).
		WithLogger(nullLogger{}).
		WithValueLogFileSize(1024 * 128)
	db, err := badger.Open(options)
	if err != nil {
		return KV{}, fmt.Errorf("opening database: %w", err)
	}
	return KV{
		db: db,
	}, nil
}

func (kv KV) Close() error {
	// TODO: GC on close
	return kv.db.Close()
}

func (kv KV) WithPrefix(prefix [prefixLen]byte) *KVPrefix {
	v := make([]byte, prefixLen)
	copy(v, prefix[:])
	return &KVPrefix{
		db:     kv.db,
		prefix: v,
	}
}

type KVPrefix struct {
	db     *badger.DB
	prefix []byte
}

func (p KVPrefix) Set(k, b []byte) error {
	return p.db.Update(func(txn *badger.Txn) error {
		return txn.Set(append(p.prefix, k...), b)
	})
}

func (p KVPrefix) SetProto(k []byte, m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}
	return p.Set(k, b)
}

func (p KVPrefix) Get(k []byte) ([]byte, error) {
	var v []byte
	err := p.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(p.prefix, k...))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			v = make([]byte, len(val))
			copy(v, val)
			return nil
		})
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		err = akcore.ErrNotFound
	}
	return v, err
}

func (p KVPrefix) GetProto(k []byte, m proto.Message) error {
	b, err := p.Get(k)
	if err != nil {
		return err
	}
	return proto.Unmarshal(b, m)
}

func (p KVPrefix) Delete(k []byte) error {
	return p.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(append(p.prefix, k...))
	})
}

func (p KVPrefix) List(prefix []byte) ([][]byte, error) {
	var keys [][]byte
	err := p.db.View(func(txn *badger.Txn) error {
		opts := badger.IteratorOptions{} // no value prefetch
		itr := txn.NewIterator(opts)
		defer itr.Close()
		prefix := append(p.prefix, prefix...)
		for itr.Seek(prefix); itr.ValidForPrefix(prefix); itr.Next() {
			item := itr.Item()
			keySrc := item.Key()
			key := make([]byte, len(keySrc)-prefixLen)
			copy(key, keySrc[prefixLen:])
			keys = append(keys, key)
		}
		return nil
	})
	return keys, err
}

type nullLogger struct{}

func (nullLogger) Errorf(string, ...interface{})   {}
func (nullLogger) Warningf(string, ...interface{}) {}
func (nullLogger) Infof(string, ...interface{})    {}
func (nullLogger) Debugf(string, ...interface{})   {}
