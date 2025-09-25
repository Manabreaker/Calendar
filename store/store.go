package store

type Store interface {
	Create(value []byte) error
	Read(key string) ([]byte, error)
	Update(key string, value []byte) error
	Delete(key string) error
}
