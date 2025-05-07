package dal

import (
	"github.com/DeamvatLab/email"
	"github.com/DeamvatLab/email/main/core"
	"github.com/DeamvatLab/email/main/dal/redis"
)

type IDataAccess interface {
	GetEmailAccount(id string) (*email.EmailAccount, error)
}

func NewDAL() IDataAccess {
	r := new(redis.RedisDAL)
	r.Init(core.ConfigProvider)
	return r
}
