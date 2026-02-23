package kv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/autonomouskoi/akcore/storage/kv"
	"github.com/stretchr/testify/require"
)

func TestBackupRestor(t *testing.T) {
	t.Parallel()

	testData := t.TempDir()
	dbPath := filepath.Join(testData, "kvdb")

	kvdb, err := kv.New(dbPath)
	require.NoError(t, err, "creating initial DB")

	kvp := kvdb.WithPrefix([8]byte{})
	testKey := []byte("test-key")
	testValue := []byte("test-value")
	require.NoError(t, kvp.Set(testKey, testValue), "setting value")
	got, err := kvp.Get(testKey)
	require.NoError(t, err, "getting initial value")
	require.Equal(t, testValue, got)
	require.NoError(t, kvdb.Close(), "closing")

	require.NoError(t, os.RemoveAll(dbPath), "deleting db")

	kvdb, err = kv.New(dbPath)
	require.NoError(t, err, "recovering DB")
	kvp = kvdb.WithPrefix([8]byte{})
	got, err = kvp.Get(testKey)
	require.NoError(t, err, "getting recovered value")
	require.Equal(t, testValue, got)
	require.NoError(t, kvdb.Close(), "closing")
}
