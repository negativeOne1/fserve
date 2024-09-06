package storage

import (
	"bufio"
	"os"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{basePath}
}

func (fs *FileStorage) GetFile(path string) (*bufio.Reader, error) {
	f, err := os.Open(fs.basePath + path)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(f), nil
}

func (fs *FileStorage) PutFile(path string, reader *bufio.Reader) error {
	f, err := os.Create(fs.basePath + path)
	if err != nil {
		return err
	}

	_, err = f.ReadFrom(reader)
	return err
}

func (fs *FileStorage) DeleteFile(path string) error {
	return nil
}
