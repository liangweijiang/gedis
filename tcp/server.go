package tcp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/liangweijiang/gedis/lib/logger"

	"github.com/liangweijiang/gedis/interface/tcp"
)

// Config 存储tcp服务的属性
type Config struct {
	Address    string        `yaml:"address"`
	MaxConnect uint32        `yaml:"max_connect"`
	Timeout    time.Duration `yaml:"timeout"`
}

//ListenAndServeWithSignal 绑定端口和处理请求，阻塞直到收到停止信号
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Infof("bind: %s, start listening...", cfg.Address)
	ListenAndServe(listener, handler, closeChan)
	return nil
}

//ListenAndServe 绑定端口并处理请求，阻塞直到关闭
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan chan struct{}) {
	go func() {
		<-closeChan
		logger.Info("shutting down...")
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accept link")
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.Handle(ctx, conn)
		}()
	}
	wg.Wait()
}
