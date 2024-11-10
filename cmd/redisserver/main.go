package main

import (
	"github.com/liangweijiang/gedis/pkg/lib/logger"
	"github.com/liangweijiang/gedis/pkg/network/tcpserver"
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
	s := tcpserver.NewTcpSever(&tcpserver.Config{Address: "0.0.0.0:9999"}, tcpserver.NewEchoHandler())
	err := s.Start()
	if err != nil {
		panic(err)
	}

}
