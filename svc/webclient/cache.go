package webclient

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type cache struct {
	cacheDir string
	client   *http.Client
	lock     sync.Mutex
	files    map[string]string
}

func newCache(cacheDir string, client *http.Client) (*cache, error) {
	c := &cache{
		cacheDir: cacheDir,
		client:   client,
		files:    map[string]string{},
	}

	des, err := os.ReadDir(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(cacheDir, 0700); err != nil {
				return nil, fmt.Errorf("creating cache dir %s: %w", cacheDir, err)
			}
			return c, nil
		}
		return nil, fmt.Errorf("reading FS: %w", err)
	}

	for _, de := range des {
		if de.IsDir() {
			continue
		}
		name := de.Name()
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		c.files[base] = name
	}
	return c, nil
}

func (c *cache) Get(url string) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	hash := sha256.Sum256([]byte(url))
	sum := hex.EncodeToString(hash[:16])
	if filename, present := c.files[sum]; present {
		return filename, nil
	}

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	resp, err := c.client.Do(r)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %d/%s", resp.StatusCode, resp.Status)
	}

	ct := resp.Header.Get("Content-Type")
	var ext string
	switch ct {
	case "image/apng":
		ext = ".apng"
	case "image/gif":
		ext = ".gif"
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/svg+xml":
		ext = ".svg"
	case "image/webp":
		ext = ".webp"
	}

	filename := sum + ext
	outfh, err := os.Create(filepath.Join(c.cacheDir, filename))
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}
	defer outfh.Close()

	if _, err := io.Copy(outfh, resp.Body); err != nil {
		return "", fmt.Errorf("copying file data: %w", err)
	}
	if err := outfh.Sync(); err != nil {
		return "", fmt.Errorf("syncing file: %w", err)
	}

	c.files[sum] = filename

	return filename, nil
}
