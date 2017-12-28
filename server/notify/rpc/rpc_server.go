package rpc

import (
	"github.com/golang/glog"
	"FishChatServer2/common/ecode"
	"FishChatServer2/protocol/rpc"
	"FishChatServer2/server/notify/conf"
	"FishChatServer2/server/notify/dao"
	sd "FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	dao       *dao.Dao
	rpcClient *RPCClient
}

func (s *RPCServer) Notify(ctx context.Context, in *rpc.NFNotifyMsgReq) (res *rpc.NFNotifyMsgRes, err error) {
	glog.Info("notify recive Notify")
	userMsgID, err := s.dao.Mysql.GetUserMsgID(ctx, in.TargetUID)
	if err != nil {
		glog.Error(err)
		return
	}
	_, err = s.dao.Mysql.UpdateUserMsgID(ctx, in.TargetUID, userMsgID.CurrentMsgID, in.TotalID)
	if err != nil {
		glog.Error(err)
		return
	}
	sendNotifyReqRPC := &rpc.ASSendNotifyReq{
		UID:              in.TargetUID,
		CurrentID:        userMsgID.CurrentMsgID,
		TotalID:          in.TotalID,
		AccessServerAddr: in.AccessServerAddr,
	}
	_, err = s.rpcClient.Access.SendNotify(ctx, sendNotifyReqRPC)
	if err != nil {
		glog.Error(err)
		return
	}
	res = &rpc.NFNotifyMsgRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func RPCServerInit(rpcClient *RPCClient) {
	glog.Info("[notify] rpc server init: ", conf.Conf.RPCServer.Addr)
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	err = sd.Register(conf.Conf.ServiceDiscoveryServer.ServiceName, conf.Conf.ServiceDiscoveryServer.RPCAddr, conf.Conf.ServiceDiscoveryServer.EtcdAddr, conf.Conf.ServiceDiscoveryServer.Interval, conf.Conf.ServiceDiscoveryServer.TTL)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpcServer := &RPCServer{
		dao:       dao.NewDao(),
		rpcClient: rpcClient,
	}
	rpc.RegisterNotifyServerRPCServer(s, rpcServer)
	s.Serve(lis)
}
