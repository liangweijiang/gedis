package redis

//Reply redis序列化协议消息的接口
type Reply interface {
	ToBytes() []byte
}
