package core

import (
	"mydtu/base"
	"mydtu/protocol/adapter"
)

type Device struct {
	link    base.Link
	adapter adapter.Adapter
}
