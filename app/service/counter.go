package service

import (
	"context"
	"github.com/gogf/gf/util/grand"
	"sb-counter/app/dao"
	"sb-counter/app/model"
)

type CounterStruct struct {
	ctx *context.Context
}

func Counter() *CounterStruct {
	ctx := context.Background()
	return &CounterStruct{
		&ctx,
	}
}

func (c *CounterStruct) Add(db, ip string) (id string) {
	id = grand.S(32, true)
	dao.Counter.Ctx(*c.ctx).Data(model.Counter{Id: id, Db: db, Ip: ip}).Insert()
	return
}
