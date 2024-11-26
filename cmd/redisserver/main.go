package main

import (
	"github.com/liangweijiang/gedis/lib/logger"
	"github.com/liangweijiang/gedis/network/tcpserver"
)

func main() {
	logger.Setup(logger.Settings{
		Path:      "./",
		Name:      "redis",
		WithColor: false,
		WithJson:  true,
	})
	s := tcpserver.NewTcpSever(&tcpserver.Config{Address: "0.0.0.0:9999"}, tcpserver.NewEchoHandler())
	err := s.Start()
	if err != nil {
		panic(err)
	}

}
