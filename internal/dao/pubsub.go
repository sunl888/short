package dao

// 订阅回调
type SubscribeCallbackFunc func(channel string, message []byte)

type Channel string

// 消息主体
type Message struct {
	// 频道
	Channel Channel

	// 消息正文
	Data []byte
}

type PubSub interface {
	// 订阅消息
	Subscribe(channel Channel, callback SubscribeCallbackFunc) error

	// 发布消息
	Publish(message Message) error

	Close()
}
