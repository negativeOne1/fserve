package storage

import "bufio"

type Storage interface {
	GetFile(path string) (*bufio.Reader, error)
	PutFile(path string, in *bufio.Reader) error
	DeleteFile(path string) error
}
