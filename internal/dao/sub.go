package dao

import (
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/log"
)

type SubscribeCallback func(channel string, message []byte)

type Subscriber struct {
	client      redis.PubSubConn             // redis 连接
	callbackMap map[string]SubscribeCallback // 订阅回调
}

func (c *Subscriber) Connect(conn redis.Conn) {
	c.client = redis.PubSubConn{Conn: conn}
	c.callbackMap = make(map[string]SubscribeCallback)
	go func() {
		for {
			switch n := c.client.Receive().(type) {
			case redis.Message:
				log.Info("Message: %s %s\n", n.Channel, n.Data)
				c.callbackMap[n.Channel](n.Channel, n.Data)
			case redis.PMessage:
				//TODO
				log.Info("PMessage: %s %s %s\n", n.Pattern, n.Channel, n.Data)
			case redis.Subscription:
				//TODO
				log.Info("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
				if n.Count == 0 {
					return
				}
			case error:
				log.Error("error: %v\n", n)
				return
			}
		}
	}()
}

func (c *Subscriber) Close() {
	err := c.client.Close()
	if err != nil {
		log.Error("redis 关闭失败")
	}
}

// 订阅
func (c *Subscriber) Subscribe(channel interface{}, callback SubscribeCallback) {
	err := c.client.Subscribe(channel)
	if err != nil {
		log.Error("redis 订阅失败")
	}
	c.callbackMap[channel.(string)] = callback
}
