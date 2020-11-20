package core

import (
	"git.zgwit.com/zgwit/dtu-admin/base"
	"git.zgwit.com/zgwit/dtu-admin/protocol/adapter"
)

type Device struct {
	link    base.Link
	adapter adapter.Adapter
}
