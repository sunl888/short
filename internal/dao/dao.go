package dao

import (
	"context"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	xtime "github.com/bilibili/kratos/pkg/time"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/short/internal/contants"
	"github.com/wq1019/short/internal/dao/orm"
	"time"
)

// Dao dao.
type Dao struct {
	db          *gorm.DB
	redis       *redis.Pool
	subRedis    *redis.Conn
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
			ShortDomain *orm.Config
		}
		rc struct {
			ShortDomain *redis.Config
			RedisExpire xtime.Duration
		}
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))
	// 获取一个长连接给 sub 使用
	cn, err := redis.Dial(rc.ShortDomain.Proto, rc.ShortDomain.Addr,
		redis.DialPassword(rc.ShortDomain.Auth),
		redis.DialWriteTimeout(time.Duration(rc.ShortDomain.WriteTimeout)),
		redis.DialDatabase(1),
	)
	checkErr(err)
	dao = &Dao{
		// mysql
		db: orm.NewMySQL(dc.ShortDomain),
		// redis
		redis:       redis.NewPool(rc.ShortDomain, redis.DialDatabase(1)),
		redisExpire: int32(time.Duration(rc.RedisExpire) / time.Second),
		subRedis:    &cn,
	}

	// redis 订阅
	var sub Subscriber
	sub.Connect(*dao.subRedis)
	// 订阅创建短链接消息
	sub.Subscribe(contants.CreateDomainChannel, func(channel string, message []byte) {
		log.Info("收到 createDomain 频道的消息: %s; 正文: %s\n", channel, string(message))
	})

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
