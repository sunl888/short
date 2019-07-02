package dao

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	xtime "github.com/bilibili/kratos/pkg/time"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/wq1019/short/internal/contants"
	"github.com/wq1019/short/internal/dao/orm"
	"sync"
	"time"
)

// Dao dao.
type Dao struct {
	db          *gorm.DB
	redis       *redis.Pool
	redisExpire int32
	//subCreateDomainFunc func(dao Dao)
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
	dao = &Dao{
		// mysql
		db: orm.NewMySQL(dc.ShortDomain),
		// redis
		redis:       redis.NewPool(rc.ShortDomain, redis.DialDatabase(1)),
		redisExpire: int32(time.Duration(rc.RedisExpire) / time.Second),
	}

	// redis 订阅
	var sub Subscriber
	ctx := context.Background()
	sub.Connect(dao.redis.Get(ctx))
	// 订阅创建短链接消息
	sub.Subscribe(contants.CreateDomainChannel, func(channel string, message []byte) {
		log.Info("收到 createDomain 频道的消息: %s; 正文: %s\n", channel, string(message))
	})

	return
}

func (d *Dao) PubSub(ctx context.Context) {
	c := d.redis.Get(ctx)
	defer c.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	psc := redis.PubSubConn{Conn: c}
	go func() {
		defer wg.Done()
		for {
			switch n := psc.Receive().(type) {
			case redis.Message:
				fmt.Printf("Message: %s %s\n", n.Channel, n.Data)
			case redis.PMessage:
				fmt.Printf("PMessage: %s %s %s\n", n.Pattern, n.Channel, n.Data)
			case redis.Subscription:
				fmt.Printf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
				if n.Count == 0 {
					return
				}
			case error:
				fmt.Printf("error: %v\n", n)
				return
			}
		}
	}()
	// 订阅创建短链接消息
	if err := psc.Subscribe(contants.CreateDomainChannel); err != nil {
		fmt.Println("订阅失败")
		err = errors.WithMessage(err, "订阅失败")
		panic(err)
	}
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
