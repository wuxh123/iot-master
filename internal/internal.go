package internal

import (
	"git.zgwit.com/zgwit/iot-admin/internal/channel"
	"git.zgwit.com/zgwit/iot-admin/internal/conf"
)

func Start() error  {
	//加载配置
	err := conf.Load()
	if err != nil {
		return err
	}

	//启动总线
	err = channel.StartDBus(conf.Config.DBus.Addr)
	if err != nil {
		return err
	}

	//恢复之前的链接
	err = channel.Recovery()
	if err != nil {
		return err
	}

	return nil
}
