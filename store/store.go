package store

type Store interface {
	Create(key string, value []byte) error
	Read(key string) ([]byte, error)
	Update(key string, value []byte) error
	Delete(key string) error
}
