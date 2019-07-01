package dao

import (
	"context"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	xtime "github.com/bilibili/kratos/pkg/time"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/short/internal/dao/orm"
	"time"
)

// Dao dao.
type Dao struct {
	db          *gorm.DB
	redis       *redis.Pool
	redisExpire int32
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// New new a dao and return.
func New() (dao *Dao) {
	var (
		dc struct {
			shortDomain *orm.Config
		}
		rc struct {
			shortDomain *redis.Config
			DemoExpire  xtime.Duration
		}
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))
	dao = &Dao{
		// mysql
		db: orm.NewMySQL(dc.shortDomain),
		// redis
		redis:       redis.NewPool(rc.shortDomain),
		redisExpire: int32(time.Duration(rc.DemoExpire) / time.Second),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	if d.db != nil {
		_ = d.db.Close()
	}
	if d.redis != nil {
		_ = d.redis.Close()
	}
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	if err = d.pingRedis(ctx); err != nil {
		log.Error("ping redis error(%v)", err)
		return
	}
	if err = d.db.DB().Ping(); err != nil {
		log.Error("ping mysql error(%v)", err)
		return
	}
	return
}

func (d *Dao) pingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}
