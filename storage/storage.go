package storage

import (
	"bufio"
	"io"
	"time"
)

const (
	NeverExpires time.Duration = -1
	DefaultTTL   time.Duration = 1 * time.Minute
)

type Storage interface {
	Get(path string) (*bufio.Reader, error)
	Save(path string, in io.Reader) error
	Delete(path string) error
}
