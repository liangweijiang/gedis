package tcpserver

import (
	"context"
	"errors"
	"github.com/liangweijiang/gedis/interfaces/network/tcp"
	"github.com/liangweijiang/gedis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Config struct {
	Address string
}

type TcpServer struct {
	listener    net.Listener
	conf        *Config
	clientCount int64
	waitDone    sync.WaitGroup
	quitCh      chan os.Signal
	errCh       chan error
	tcpHandler  tcp.Handler
}

func NewTcpSever(conf *Config, handler tcp.Handler) *TcpServer {
	return &TcpServer{
		conf:        conf,
		clientCount: 0,
		quitCh:      make(chan os.Signal),
		errCh:       make(chan error),
		tcpHandler:  handler,
	}
}

func (s *TcpServer) Start() error {
	signal.Notify(s.quitCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	listener, err := net.Listen("tcp", s.conf.Address)
	if err != nil {
		logger.Errorf("listen error: %v", err)
		return err
	}
	s.listener = listener
	s.listenAndServe()
	return nil
}

func (s *TcpServer) listenAndServe() {
	go func() {
		select {
		case <-s.quitCh:
		case err := <-s.errCh:
			logger.Errorf("listenAndServe error: %v", err)
		}
		_ = s.listener.Close()
		_ = s.tcpHandler.Close()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				time.Sleep(5 * time.Millisecond)
				continue
			}
			s.errCh <- err
			break
		}
		go s.handleConn(conn)
	}
	s.waitDone.Wait()
}

func (s *TcpServer) handleConn(conn net.Conn) {
	atomic.AddInt64(&s.clientCount, 1)
	s.waitDone.Add(1)
	defer func() {
		s.waitDone.Done()
		atomic.AddInt64(&s.clientCount, -1)
	}()
	s.tcpHandler.Handle(context.Background(), conn)
}
