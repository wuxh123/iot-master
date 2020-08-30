package rpc

import (
	"context"
	"github.com/zgwit/dtu-admin/plugin"
	"log"
)

type pluginServer struct {
}

func (r *pluginServer) Register(ctx context.Context, req *plugin.RegisterReq) (*plugin.RegisterResp, error) {
	return nil, nil
}

func (r *pluginServer) Start(stream plugin.Plugin_StartServer) error{
	//stream.Send()
	//stream.Recv()
	for  {
		resp, err := stream.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		//TODO 处理
		log.Println(resp)
	}

	return nil
}