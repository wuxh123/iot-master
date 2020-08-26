package dtu

type Link interface {
	Open() error
	Close() error
}
