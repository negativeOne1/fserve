package storage

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type FileStorage struct {
	basePath   string
	mu         sync.RWMutex
	janitor    *janitor
	expiration time.Duration
}

func NewFileStorage(
	basePath string,
	defaultExpiration, cleanupInterval time.Duration,
) (*FileStorage, error) {
	fs := &FileStorage{
		basePath:   basePath,
		expiration: defaultExpiration,
	}

	if cleanupInterval > 0 {
		startJanitor(fs, cleanupInterval)
		runtime.SetFinalizer(fs, stopJanitor)
	}

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *FileStorage) Get(path string) (*bufio.Reader, error) {
	fs.mu.RLock()
	f, err := os.Open(fs.basePath + path)
	fs.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	return bufio.NewReader(f), nil
}

func (fs *FileStorage) Save(name string, reader io.Reader) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	p := path.Join(fs.basePath, name)

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Debug().Str("path", p).Msg("Saving file")

	_, err = f.ReadFrom(reader)
	return err
}

func (fs *FileStorage) Delete(path string) error {
	return errors.New("not implemented")
}

// maybe this should be at the janitor
func (fs *FileStorage) DeleteExpired() {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	files, err := os.ReadDir(fs.basePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read directory")
		return
	}

	d := make([]os.DirEntry, 0)

	for _, f := range files {
		i, err := f.Info()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get file info")
			continue
		}
		if time.Now().After(i.ModTime().Add(fs.expiration)) {
			os.Remove(path.Join(fs.basePath, f.Name()))
		}
	}

	for _, f := range d {
		p := path.Join(fs.basePath, f.Name())
		log.Debug().Str("path", p).Msg("Deleting expired file")
		os.Remove(p)
	}
}
