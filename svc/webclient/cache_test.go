package webclient_test

/*
func TestGet(t *testing.T) {
	t.Parallel()
	requireT := require.New(t)

	server := map[string]webcache.FakeFileEntry{
		"12345": {
			ContentType: "image/jpeg",
			Data:        []byte("blarg blarg honk"),
		},
	}
	client := webcache.NewFileServerClient(server)

	tmpDir := t.TempDir()

	cache, err := webcache.New(tmpDir, client)
	requireT.NoError(err, "creating new cache")

	result, err := cache.Get("not-found")
	requireT.ErrorContains(err, "non-200")
	requireT.Zero(result)

	result, err = cache.Get("12345")
	requireT.NoError(err, "getting 12345")
	requireT.Equal("5994471abb01112afcc18159f6cc74b4.jpg", result)

	stored, err := os.ReadFile(filepath.Join(tmpDir, result))
	requireT.NoError(err, "reading cached data")
	requireT.Equal(stored, server["12345"].Data)

	// verify that the cache uses the cached result
	delete(server, "12345")
	// we can still retrieve it, even though it's ont on the server
	result, err = cache.Get("12345")
	requireT.NoError(err, "getting 12345")
	requireT.Equal("5994471abb01112afcc18159f6cc74b4.jpg", result)
}
*/
