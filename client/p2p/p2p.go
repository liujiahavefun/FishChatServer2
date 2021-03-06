package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"FishChatServer2/codec"
	"FishChatServer2/libnet"
	"FishChatServer2/protocol/external"
	"time"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func checkErr(err error) {
	if err != nil {
		glog.Error(err)
	}
}

var _currentID int64

func clientLoop(session *libnet.Session, protobuf *codec.ProtobufProtocol) {
	var err error
	var clientMsg *libnet.Session
	err = session.Send(&external.ReqAccessServer{
		Cmd: external.ReqAccessServerCMD,
	})
	checkErr(err)
	rsp, err := session.Receive()
	checkErr(err)
	glog.Info(string(rsp))
	if rsp != nil {
		baseCMD := &external.Base{}
		if err = proto.Unmarshal(rsp, baseCMD); err != nil {
			glog.Error(err)
		}
		switch baseCMD.Cmd {
		case external.ReqAccessServerCMD:
			resSelectMsgServerForClientPB := &external.ResSelectAccessServerForClient{}
			if err = proto.Unmarshal(rsp, resSelectMsgServerForClientPB); err != nil {
				glog.Error(err)
			}
			glog.Info(resSelectMsgServerForClientPB)
			glog.Info(resSelectMsgServerForClientPB.Addr)
			clientMsg, err = libnet.Connect("tcp", resSelectMsgServerForClientPB.Addr, protobuf, 0)
			checkErr(err)
		}
	}
	fmt.Print("输入我的id :")
	var myID int64
	if _, err := fmt.Scanf("%d\n", &myID); err != nil {
		glog.Error(err.Error())
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(fmt.Sprintf("%d", myID)))
	md5Ctx.Write([]byte("!0h#?123(ABM"))
	cipherStr := md5Ctx.Sum(nil)
	calcToken := hex.EncodeToString(cipherStr)
	err = clientMsg.Send(&external.ReqLogin{
		Cmd:   external.LoginCMD,
		UID:   myID,
		Token: calcToken,
	})
	checkErr(err)
	rsp, err = clientMsg.Receive()
	checkErr(err)
	glog.Info(string(rsp))
	go func() {
		for {
			err = clientMsg.Send(&external.ReqPing{
				Cmd: external.PingCMD,
				UID: myID,
			})
			checkErr(err)
			time.Sleep(5 * time.Second)
		}
	}()
	// glog.Info(string(rsp))
	go func() {
		for {
			rsp, err := clientMsg.Receive()
			if err != nil {
				glog.Error(err.Error())
			}
			if rsp != nil {
				baseCMD := &external.Base{}
				if err = proto.Unmarshal(rsp, baseCMD); err != nil {
					continue
				}
				switch baseCMD.Cmd {
				case external.LoginCMD:
					resLogin := &external.ResLogin{}
					if err = proto.Unmarshal(rsp, resLogin); err != nil {
						glog.Error(err)
					}
					fmt.Printf("收到登录返回: 返回码[%d]", resLogin.ErrCode)
					fmt.Println()
				case external.PingCMD:

				case external.SendP2PMsgCMD:
					resSendP2PMsg := &external.ResSendP2PMsg{}
					if err = proto.Unmarshal(rsp, resSendP2PMsg); err != nil {
						glog.Error(err)
					}
					fmt.Printf("收到点对点消息: 返回码[%d], 对方ID[%d], 消息内容[%s]", resSendP2PMsg.ErrCode, resSendP2PMsg.SourceUID, resSendP2PMsg.Msg)
					fmt.Println()
				case external.NotifyCMD:
					glog.Info("recive NotifyCMD")
					resNotify := &external.ResNotify{}
					if err = proto.Unmarshal(rsp, resNotify); err != nil {
						glog.Error(err)
					}
					fmt.Println(resNotify.CurrentID)
					err = clientMsg.Send(&external.ReqSyncMsg{
						Cmd:       external.SyncMsgCMD,
						UID:       myID,
						CurrentID: _currentID,
					})
					// _currentID = resNotify.CurrentID
				case external.SyncMsgCMD:
					resSyncMsg := &external.ResSyncMsg{}
					if err = proto.Unmarshal(rsp, resSyncMsg); err != nil {
						glog.Error(err)
					}
					fmt.Println(resSyncMsg.Msgs)
					for _, msg := range resSyncMsg.Msgs {
						fmt.Printf("收到点对点消息: 消息类型[%s], 对方ID[%d], 消息内容[%s]", msg.MsgType, msg.SourceUID, msg.Msg)
						fmt.Println()
					}
					_currentID = resSyncMsg.CurrentID
				}
			}
		}
	}()
	go func() {
		for {
			err = clientMsg.Send(&external.ReqSyncMsg{
				Cmd:       external.SyncMsgCMD,
				UID:       myID,
				CurrentID: _currentID,
			})
			time.Sleep(10 * time.Second)
		}
	}()
	for {
		glog.Info("send p2p msg")
		var targetID int64
		fmt.Print("输入对方的id :")
		if _, err = fmt.Scanf("%d\n", &targetID); err != nil {
			glog.Error(err.Error())
		}
		var msg string
		fmt.Print("输入你想说的话 :")
		if _, err = fmt.Scanf("%s\n", &msg); err != nil {
			glog.Error(err.Error())
		}
		err = clientMsg.Send(&external.ReqSendP2PMsg{
			Cmd:       external.SendP2PMsgCMD,
			SourceUID: myID,
			TargetUID: targetID,
			Msg:       msg,
		})
	}
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:10000", "server address")
	flag.Parse()
	protobuf := codec.Protobuf()
	client, err := libnet.Connect("tcp", addr, protobuf, 0)
	checkErr(err)
	clientLoop(client, protobuf)
}
