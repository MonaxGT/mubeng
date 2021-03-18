package storage

type Database interface {
	Open() error
	Load() ([]string, error)
}
