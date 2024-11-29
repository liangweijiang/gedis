package redis

// Connection represents a connection with access client
type Connection interface {
	Write([]byte) (int, error)
	Close() error
	RemoteAddr() string

	SetPassword(string)
	GetPassword() string

	InMultiState() bool

	GetDBIndex() int
	SelectDB(int)
}
