package kv

import (
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore"
)

const (
	prefixLen = 8
)

// KV provides key-value storage
type KV struct {
	db *badger.DB
}

// New creates a new KV store. This should be invoked by AK itself, not modules
func New(dbPath string) (KV, error) {
	options := badger.DefaultOptions(dbPath).
		WithLogger(nullLogger{}).
		WithValueLogFileSize(1024 * 1024)
	db, err := badger.Open(options)
	if err != nil {
		return KV{}, fmt.Errorf("opening database: %w", err)
	}
	return KV{
		db: db,
	}, nil
}

// Close the KV
func (kv KV) Close() error {
	// TODO: GC on close
	opts := kv.db.Opts()
	backupPath := filepath.Clean(opts.Dir) + ".backup.gz"
	closing := func(err error) error {
		closeErr := kv.db.Close()
		if closeErr != nil {
			err = errors.Join(
				err,
				fmt.Errorf("closing database: %w", closeErr),
			)
			return err
		}
		return err
	}

	outfh, err := os.Create(backupPath)
	if err != nil {
		return closing(fmt.Errorf("creating backup file %s: %w", backupPath, err))
	}
	defer outfh.Close()

	gzw := gzip.NewWriter(outfh)

	if _, err := kv.db.Backup(gzw, 0); err != nil {
		return closing(fmt.Errorf("backing up database to %s: %w", backupPath, err))
	}

	if err := gzw.Close(); err != nil {
		return closing(fmt.Errorf("finishing backup compression of %s: %w", backupPath, err))
	}
	if err := outfh.Sync(); err != nil {
		return closing(fmt.Errorf("syncing database backup %s: %w", backupPath, err))
	}

	if err := kv.db.Close(); err != nil {
		return fmt.Errorf("closing database: %w", err)
	}
	return nil
}

// WithPrefix creates a wrapped KV where all keys are forced to have a given
// prefix. This segregates values from each module
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

// Set a value in the KV
func (p KVPrefix) Set(k, b []byte) error {
	return p.db.Update(func(txn *badger.Txn) error {
		return txn.Set(append(p.prefix, k...), b)
	})
}

// SetProto marshals a proto and sets it as a value in the KV
func (p KVPrefix) SetProto(k []byte, m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}
	return p.Set(k, b)
}

// Get a value. If no value is found with this key, akcore.ErrNotFound is
// returned.
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

// GetProto wraps Get, unmarshalling a retrieved proto into m
func (p KVPrefix) GetProto(k []byte, m proto.Message) error {
	b, err := p.Get(k)
	if err != nil {
		return err
	}
	return proto.Unmarshal(b, m)
}

// Delete a value from the store, if the key is present
func (p KVPrefix) Delete(k []byte) error {
	return p.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(append(p.prefix, k...))
	})
}

// List keys matching a given prefix.
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

// For tests
func NewMemory() (KV, error) {
	options := badger.DefaultOptions("").
		WithLogger(nullLogger{}).
		WithInMemory(true)
	db, err := badger.Open(options)
	if err != nil {
		return KV{}, fmt.Errorf("opening database: %w", err)
	}
	return KV{
		db: db,
	}, nil
}

// badger wants to log stuff; we provide a do-nothing logger to discard the logs
type nullLogger struct{}

func (nullLogger) Errorf(string, ...interface{})   {}
func (nullLogger) Warningf(string, ...interface{}) {}
func (nullLogger) Infof(string, ...interface{})    {}
func (nullLogger) Debugf(string, ...interface{})   {}
