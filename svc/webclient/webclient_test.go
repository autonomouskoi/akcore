package webclient_test

/*
func TestStaticDownload(t *testing.T) {
	t.Parallel()

	cacheDir := t.TempDir()
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	server := map[string]FakeFileEntry{
		"12345": {
			ContentType: "image/jpeg",
			Data:        []byte("blarg blarg honk"),
		},
	}

	wc, err := webclient.New(&modutil.Deps{
		Log:        log,
		HttpClient: NewFileServerClient(server),
		CachePath:  cacheDir,
	})
	require.NoError(t, err, "creating web client")

	msg := &bus.BusMessage{}
	wc.MarshalMessage(msg, &svc.WebclientStaticDownloadRequest{URL: "12345"})
	require.Nil(t, msg.Error)

	reply := wc.HandleRequestStaticDownload(msg)
	require.NotNil(t, reply)
	require.Nil(t, reply.Error)

	resp := &svc.WebclientStaticDownloadResponse{}
	require.Nil(t, wc.UnmarshalMessage(reply, resp))
	require.Equal(t, "/s/webclient/c/5994471abb01112afcc18159f6cc74b4.jpg", resp.Path)

	reply = wc.HandleRequestStaticDownload(msg)
	require.NotNil(t, reply)
	require.Nil(t, reply.Error)
	resp.Reset()
	require.Nil(t, wc.UnmarshalMessage(reply, resp))
	require.Equal(t, "/s/webclient/c/5994471abb01112afcc18159f6cc74b4.jpg", resp.Path)

	wc.MarshalMessage(msg, &svc.WebclientStaticDownloadRequest{URL: "not-found"})
	require.Nil(t, msg.Error)
	reply = wc.HandleRequestStaticDownload(msg)
	require.NotNil(t, reply)
	require.NotNil(t, reply.Error)
}
*/
