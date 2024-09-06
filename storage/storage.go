package storage

type Storage interface {
	GetFile(path string) ([]byte, error)
	PutFile(path string, data []byte) error
	DeleteFile(path string) error
}
