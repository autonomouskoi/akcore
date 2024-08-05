package webutil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func ZipOrEnvPath(env string, zipBytes []byte) (http.FileSystem, error) {
	if contentPath := os.Getenv(env); contentPath != "" {
		return http.Dir(contentPath), nil
	}
	zipR, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("creating zip reader: %w", err)
	}
	return http.FS(zipR), nil
}
