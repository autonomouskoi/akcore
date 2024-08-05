package kv

import (
	"testing"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/require"
)

func TestKVPrefix(t *testing.T) {
	t.Parallel()
	requireT := require.New(t)

	kv, err := New(t.TempDir())
	requireT.NoError(err, "opening db")
	t.Cleanup(func() { kv.Close() })

	prefixArr := [...]byte{0, 1, 2, 3, 4, 5, 6, 7}

	kvp := kv.WithPrefix(prefixArr)

	test1 := []byte("test1")

	// set via kvprefix
	requireT.NoError(kvp.Set(test1, test1), "setting test1")
	// get via kvprefix
	got, err := kvp.Get(test1)
	requireT.NoError(err, "getting test1")
	requireT.Equal(test1, got)
	// get directly from the db to verify it's set with the prefix
	err = kv.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(prefixArr[:], test1...))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			got = make([]byte, len(val))
			copy(got, val)
			return nil
		})
	})
	requireT.NoError(err, "getting test1")
	requireT.Equal(test1, got)

	// set a new value
	testV2 := []byte("testv2")
	requireT.NoError(kvp.Set(test1, testV2), "setting test1")
	got, err = kvp.Get(test1)
	requireT.NoError(err, "getting test1")
	requireT.Equal(testV2, got)

	// set more values
	boop := []byte("boop")
	test2 := []byte("test2")
	test3 := []byte("test3")
	zoom := []byte("zoom")
	for _, key := range [][]byte{
		boop, test2, test3, zoom,
	} {
		value := []byte("foo")
		requireT.NoError(kvp.Set(key, value), "setting ", string(key))
		got, err := kvp.Get(key)
		requireT.NoError(err, "getting ", string(key))
		requireT.Equal(value, got)
	}

	// list all keys
	want := [][]byte{boop, test1, test2, test3, zoom}
	gotKeys, err := kvp.List(nil)
	requireT.NoError(err, "listing all keys")
	requireT.Equal(want, gotKeys)

	// list only test keys
	want = [][]byte{test1, test2, test3}
	gotKeys, err = kvp.List([]byte("test"))
	requireT.NoError(err, "listing test keys")
	requireT.Equal(want, gotKeys)
}
