package store

type Store interface {
	Set() bool
	Get()
}
