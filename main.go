package main

import (
	"fmt"

	"github.com/liangweijiang/gedis/lib/logger"

	"github.com/liangweijiang/gedis/tcp"
)

func main() {
	logger.InitLog(&logger.Config{
		Path:       "logs",
		Name:       "gedis",
		Ext:        "log",
		TimeFormat: "2006-01-02 15:04:05",
	})

	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", "0.0.0.0", 8080),
	}, tcp.NewEchoHandler())
	if err != nil {
		panic(err)
	}
}
