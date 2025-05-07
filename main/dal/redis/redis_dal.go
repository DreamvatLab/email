package redis

import (
	"context"
	"encoding/json"

	"github.com/DeamvatLab/email"
	"github.com/DeamvatLab/email/main/core"
	"github.com/DreamvatLab/go/xconfig"
	"github.com/DreamvatLab/go/xerr"
	"github.com/DreamvatLab/go/xredis"
	"github.com/redis/go-redis/v9"
)

const (
	_key = "account:Emails"
)

type RedisDAL struct {
	client redis.UniversalClient
}

func (o *RedisDAL) Init(config xconfig.IConfigProvider) {
	var rc *xredis.RedisConfig
	core.ConfigProvider.GetStruct("Redis", &rc)
	o.client = xredis.NewClient(rc)
}

// GetEmailAccount implements dal.IDataAccess.
func (o *RedisDAL) GetEmailAccount(id string) (*email.EmailAccount, error) {
	jsonBytes, err := o.client.HGet(context.Background(), _key, id).Bytes()
	if err != nil {
		return nil, xerr.WithStack(err)
	}

	var r *email.EmailAccount
	err = json.Unmarshal(jsonBytes, &r)
	if err != nil {
		return nil, xerr.WithStack(err)
	}

	return r, nil
}
