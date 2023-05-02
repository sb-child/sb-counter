package database

import (
	"context"
	"sb-counter/internal/consts"
	"sb-counter/internal/model/entity"
	"sb-counter/internal/service"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
)

type sDatabase struct {
}

func init() {
	service.RegisterDatabase(New())
}
func New() *sDatabase {
	return &sDatabase{}
}

func (s *sDatabase) Add(ctx context.Context, db string, ip string) error {
	tx, err := g.DB("default").Begin(ctx)
	if err != nil {
		g.Log().Error(ctx, "Database(Add): failed to start transaction:", err)
		return nil
	}
	id := grand.S(64, true)
	_, err = tx.Model("counter").Data(entity.Counter{Id: id, Db: db, Ip: ip}).Insert()
	if err != nil {
		g.Log().Error(ctx, "Database(Add): failed to insert:", db, err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *sDatabase) FetchData(ctx context.Context, db string) *consts.CounterData {
	cd := consts.CounterData{}
	tx, err := g.DB("default").Begin(ctx)
	if err != nil {
		g.Log().Error(ctx, "Database(FetchData): failed to start transaction:", err)
		return nil
	}
	data_all, err := tx.Model("counter").Where("db", db).Count()
	if err != nil {
		g.Log().Error(ctx, "Database(FetchData): failed to get data_all for user", db, err)
		tx.Rollback()
		return nil
	}
	cd.All = data_all
	for _, i := range []int{1, 2, 3} {
		today := time.Now()
		today = time.Date(today.Year(), today.Month(), today.Day()-i+2, 0, 0, 0, 0, today.Location())
		yesterday := today.Add(time.Hour * -24)
		todayG := gtime.NewFromTime(today)
		yesterdayG := gtime.NewFromTime(yesterday)
		r, err := tx.Model("counter").
			Where("db", db).
			WhereBetween("created_at", yesterdayG, todayG).
			Count()
		if err != nil {
			g.Log().Error(ctx, "Database(FetchData): failed to get time range for user", i, db, err)
			tx.Rollback()
			return nil
		}
		switch i {
		case 1:
			cd.Today = r
		case 2:
			cd.Yesterday = r
		case 3:
			cd.BeforeYesterday = r
		}
	}
	tx.Commit()
	return &cd
}
