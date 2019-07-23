package event

import "github.com/bilibili/kratos/pkg/cache/redis"

type SubscribeCallback func(channel string, message []byte)

type Sub struct {
	Conn           *redis.PubSubConn
	SubCallbackMap map[string]SubscribeCallback
}

type Subscriber interface {
	Connect()
	Subscribe(conn *redis.Conn)
}

func (*Sub) Connect() {

}

func (*Sub) Subscribe(conn *redis.Conn) {

}
