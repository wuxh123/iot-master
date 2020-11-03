package core

import (
	"git.zgwit.com/zgwit/iot-admin/base"
	"git.zgwit.com/zgwit/iot-admin/protocol/adapter"
)

type Device struct {
	link    base.Link
	adapter adapter.Adapter
}
