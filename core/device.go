package core

import (
	"git.zgwit.com/zgwit/MyDTU/base"
	"git.zgwit.com/zgwit/MyDTU/protocol/adapter"
)

type Device struct {
	link    base.Link
	adapter adapter.Adapter
}
