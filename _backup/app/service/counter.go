package service

import (
	"context"
	"sb-counter/internal/service/internal/dao"
	"sb-counter/app/model"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
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
retry:
	id = grand.S(64, true)
	_, err := dao.Counter.Ctx(*c.ctx).Data(model.Counter{Id: id, Db: db, Ip: ip}).Insert()
	if err != nil {
		goto retry
	}
	return
}

func (c *CounterStruct) GetAll(db string) (v int) {
	r, _ := dao.Counter.Ctx(*c.ctx).Where("db", db).Count()
	return r
}

func (c *CounterStruct) GetDay(db string, offset int) (v int) {
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day()-offset+2, 0, 0, 0, 0, today.Location())
	yesterday := today.Add(time.Hour * -24)
	todayG := gtime.NewFromTime(today)
	yesterdayG := gtime.NewFromTime(yesterday)
	r, err := dao.Counter.Ctx(*c.ctx).
		Where("db", db).
		WhereBetween("created_at", yesterdayG, todayG).
		Count()
	g.Log().Debug(r, err, todayG, yesterdayG)
	return r
}
