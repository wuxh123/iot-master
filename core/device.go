package core

import (
	"git.zgwit.com/iot/mydtu/base"
	"git.zgwit.com/iot/mydtu/protocol/adapter"
)

type Device struct {
	link    base.Link
	adapter adapter.Adapter
}
