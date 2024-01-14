package store

type IStore interface {
	Set() bool
	Get() // how input and out put should be?
}
