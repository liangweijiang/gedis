package main

import (
	"github.com/liangweijiang/gedis/lib/logger"
	tcpserver2 "github.com/liangweijiang/gedis/network/tcpserver"
)

func main() {
	logger.Setup(logger.Settings{
		Path:      "./",
		Name:      "redis",
		WithColor: false,
		WithJson:  true,
	})
	logger.Info("success")
	logger.Info("success")
	logger.Info("success")
	logger.Info("success")
	logger.Info("success")
	s := tcpserver2.NewTcpSever(&tcpserver2.Config{Address: "0.0.0.0:9999"}, tcpserver2.NewEchoHandler())
	err := s.Start()
	if err != nil {
		panic(err)
	}

}
