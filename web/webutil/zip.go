package webutil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"os"
)

// ZipOrEnvPath returns a FileSystem to serve web content from. If the name env
// refers to is not empty it is assumed to be a valid filesystem path and is
// used to provide content. If it is empty, zipBytes is assumed to be a zip file
// and is used to provide content.
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
