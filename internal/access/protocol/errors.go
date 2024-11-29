package protocol

// ArgNumErrReply represents wrong number of arguments for command
type ArgNumErrReply struct {
	Cmd string
}

// NewArgNumErrReply represents wrong number of arguments for command
func NewArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

// ToBytes marshals redis.Reply
func (r *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + r.Cmd + "' command\r\n")
}

func (r *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + r.Cmd + "' command"
}
