package storage

type S3 struct{}

func NewS3() *S3 {
	return &S3{}
}

func (s3 *S3) GetFile(path string) ([]byte, error) {
	return nil, nil
}

func (s3 *S3) PutFile(path string, data []byte) error {
	return nil
}

func (s3 *S3) DeleteFile(path string) error {
	return nil
}
