package storage

type FileStorage struct{}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (fs *FileStorage) GetFile(path string) ([]byte, error) {
	return nil, nil
}

func (fs *FileStorage) PutFile(path string, data []byte) error {
	return nil
}

func (fs *FileStorage) DeleteFile(path string) error {
	return nil
}
