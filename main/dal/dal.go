package dal

import (
	"github.com/DreamvatLab/email"
	"github.com/DreamvatLab/email/main/core"
	"github.com/DreamvatLab/email/main/dal/redis"
)

type IDataAccess interface {
	GetEmailAccount(id string) (*email.EmailAccount, error)
}

func NewDAL() IDataAccess {
	r := new(redis.RedisDAL)
	r.Init(core.ConfigProvider)
	return r
}
