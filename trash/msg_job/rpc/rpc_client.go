package rpc

import (
	"github.com/golang/glog"
	"FishChatServer2/jobs/msg_job/rpc/client"
)

type RPCClient struct {
	AccessServer *client.AccessServerRPCCli
}

func NewRPCClient() (c *RPCClient, err error) {
	accessServer, err := client.NewAccessServerRPCCli()
	if err != nil {
		glog.Error(err)
		return
	}
	c = &RPCClient{
		AccessServer: accessServer,
	}
	return
}
