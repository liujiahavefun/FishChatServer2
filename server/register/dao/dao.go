package dao

import (
	// "github.com/golang/glog"
	"FishChatServer2/common/dao/xredis"
	"FishChatServer2/server/register/conf"
)

type Dao struct {
	redis *xredis.Pool
	Mysql *Mysql
}

func NewDao() (dao *Dao) {
	mysql := NewMysql()
	dao = &Dao{
		redis: xredis.NewPool(conf.Conf.Redis.Redis),
		Mysql: mysql,
	}
	return
}
