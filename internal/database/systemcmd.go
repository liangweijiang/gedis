package database

import (
	"fmt"
	"github.com/liangweijiang/gedis/config"
	"github.com/liangweijiang/gedis/interfaces/redis"
	"github.com/liangweijiang/gedis/internal/access/protocol"
	"github.com/liangweijiang/gedis/lib/logger"
	"github.com/liangweijiang/gedis/network/tcpserver"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var gedisVersion = "1.0.0"

const (
	commandServer      = "server"
	commandClient      = "client"
	commandMemory      = "memory"
	commandPersistence = "persistence"
	commandStat        = "stat"
	commandReplication = "replication"
	commandCPU         = "cpu"
	commandModule      = "module"
	commandErrorStat   = "errorstat"
	commandCluster     = "cluster"
	commandKeySpace    = "keyspace"
)

func Ping(args [][]byte) redis.Reply {
	if len(args) == 0 {
		return protocol.NewSimpleStringReply("PONG")
	} else if len(args) == 1 {
		return protocol.NewSimpleStringReply(string(args[0]))
	} else {
		return protocol.NewSimpleErrorReply("ERR wrong number of arguments for 'ping' command")
	}
}

func Auth(c redis.Connection, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.NewSimpleErrorReply("ERR wrong number of arguments for 'auth' command")
	}

	passwd := string(args[0])
	if config.GetProperties().RequirePass == "" {
		return protocol.NewSimpleErrorReply("ERR Client sent AUTH, but no password is set")
	}
	if config.GetProperties().RequirePass != passwd {
		return protocol.NewSimpleErrorReply("WRONGPASS invalid username-password pair or user is disabled.")
	}
	c.SetPassword(passwd)
	return protocol.NewSimpleStringReply("OK")
}

func DBSize(s *Server) redis.Reply {
	return protocol.NewIntegerReply(0)
}

func Info(db *Server, args [][]byte) redis.Reply {
	if len(args) == 0 {
		infoCommandList := [...]string{commandServer, commandClient, commandMemory, commandPersistence, commandStat,
			commandReplication, commandCPU, commandModule, commandErrorStat, commandCluster, commandKeySpace}

		var allSection []byte
		for _, command := range infoCommandList {
			allSection = append(allSection, GenGedisInfoString(command, db)...)
		}
		return protocol.NewBulkStringReply(allSection)
	} else if len(args) == 1 {
		command := strings.ToLower(string(args[0]))
		return protocol.NewBulkStringReply(GenGedisInfoString(command, db))
	}
	return protocol.NewSimpleErrorReply("ERR syntax error")
}

func SlaveOf(s *Server, args [][]byte) redis.Reply {
	if strings.ToLower(string(args[0])) == "no" && strings.ToLower(string(args[1])) == "one" {
		return protocol.NewSimpleStringReply("OK")
	}
	host := string(args[0])
	port, err := strconv.Atoi(string(args[1]))
	if err != nil {
		logger.Errorf("parse port failed")
		return protocol.NewSimpleErrorReply("ERR value is not an integer or out of range")
	}
	s.slaveStatus.mutex.Lock()
	atomic.StoreInt32(&s.role, slaveRole)
	s.slaveStatus.masterHost = host
	s.slaveStatus.masterPort = port
	s.slaveStatus.mutex.Unlock()
	go s.slaveStatus.setupMaster()
	return protocol.NewSimpleStringReply("OK")
}

func GenGedisInfoString(command string, s *Server) []byte {
	startUpTimeFromNow := getGedisRunningTime()
	switch command {
	case commandServer:
		s := fmt.Sprintf("# Server\r\n"+
			"gedis_version:%s\r\n"+
			//"gedis_git_sha1:00000000\r\n"+
			//"gedis_git_dirty:0\r\n"+
			//"gedis_build_id:240c2fa521c5f553\r\n"+*//*
			"gedis_mode:%s\r\n"+
			"os:%s %s\r\n"+
			"arch_bits:%d\r\n"+
			//"monotonic_clock:POSIX clock_gettime\r\n"+
			"multiplexing_api:epoll\r\n"+
			"atomicvar_api:c11-builtin\r\n"+
			//"gcc_version:10.2.1\r\n"+
			"process_id:%d\r\n"+
			"process_supervised:no\r\n"+
			"run_id:%s\r\n"+
			"tcp_port:%d\r\n"+
			"uptime_in_seconds:%d\r\n"+
			"uptime_in_days:%d\r\n"+
			//"hz:10\r\n"+
			//"configured_hz:10\r\n"+
			//"lru_clock:4712535\r\n"+
			//"executable:/data/gedis-server\r\n"+
			"config_file:%s\r\n",
			// "io_threads_active:0\r\n",
			gedisVersion, config.GetRunningMode(), runtime.GOOS, runtime.GOARCH, 32<<(^uint(0)>>63),
			os.Getpid(), config.GetProperties().RunID, config.GetProperties().Port, startUpTimeFromNow,
			startUpTimeFromNow/time.Duration(3600*24), config.GetProperties().CfPath)
		return []byte(s)
	case commandClient:
		s := fmt.Sprintf("# Clients\r\n"+
			"connected_clients:%d\r\n",
			//"cluster_connections:0\r\n"+
			//"maxclients:10000\r\n"+
			//"client_recent_max_input_buffer:32\r\n"+
			//"client_recent_max_output_buffer:0\r\n"+
			//"blocked_clients:0\r\n"+
			//"tracking_clients:0\r\n"+
			//"clients_in_timeout_table:0\r\n",
			tcpserver.GetClientCount())
		return []byte(s)
	case commandCluster:
		var s string
		if config.GetRunningMode() == config.ClusterMode {
			s = fmt.Sprintf("# Cluster\r\n"+
				"cluster_enabled:%s\r\n",
				"1",
			)
		} else {
			s = fmt.Sprintf("# Cluster\r\n"+
				"cluster_enabled:%s\r\n",
				"0",
			)
		}
		return []byte(s)
	default:
		return []byte("")
	}
}

func getGedisRunningTime() time.Duration {
	return time.Since(config.GetServerInfo().StartUpTime) / time.Second
}

func isAuthenticated(c redis.Connection) bool {
	if config.GetProperties().RequirePass == "" {
		return true
	}
	return c.GetPassword() == config.GetProperties().RequirePass
}
