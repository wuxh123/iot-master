package types

type Device interface {
	Read(name string) (interface{}, error)
	Write(name string, value interface{}) error
}
