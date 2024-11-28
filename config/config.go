package config

import (
	"bufio"
	"github.com/liangweijiang/gedis/lib/logger"
	"github.com/liangweijiang/gedis/lib/utils"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	ClusterMode    = "cluster"
	StandaloneMode = "standalone"
)

type ServerProperties struct {
	Bind  string `conf:"bind"`
	Port  int    `conf:"port"`
	Dir   string `conf:"dir"`
	RunID string `conf:"runid"`

	Databases int `conf:"databases"`

	AppendOnly     bool   `conf:"appendonly"`
	AppendFilename string `conf:"appendfilename"`
	AppendFsync    string `conf:"appendfsync"`
	RequirePass    string `conf:"requirepass,omitempty"`

	ClusterEnabled string `conf:"cluster-enabled"` // Not used at present.

	CfPath string `conf:"cf,omitempty"`
}

type ServerInfo struct {
	StartUpTime time.Time
}

var properties *ServerProperties
var serverInfo *ServerInfo

func GetProperties() *ServerProperties {
	return properties
}

func GetServerInfo() *ServerInfo {
	return serverInfo
}

func SetupConfig(configFilename string) {
	file, err := os.Open(configFilename)
	if err != nil {
		logger.Fatal(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Errorf("file close err:%+v", err)
			return
		}
	}(file)

	properties = parse(file)
	properties.RunID = utils.RandString(40)
	configFilePath, err := filepath.Abs(configFilename)
	if err != nil {
		return
	}
	properties.CfPath = configFilePath
	if properties.Dir == "" {
		properties.Dir = "."
	}
}

func parse(r io.Reader) *ServerProperties {
	config := new(ServerProperties)

	rawMap := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && strings.TrimLeft(line, " ")[0] == '#' {
			continue
		}
		pivot := strings.IndexAny(line, " ")
		if pivot > 0 && pivot < len(line)-1 {
			key := line[0:pivot]
			value := strings.Trim(line[pivot+1:], " ")
			rawMap[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	n := t.Elem().NumField()
	for i := 0; i < n; i++ {
		field := t.Elem().Field(i)
		fieldVal := v.Elem().Field(i)
		key, ok := field.Tag.Lookup("conf")
		if !ok || strings.Trim(key, " ") == "" {
			key = field.Name
		}
		value, ok := rawMap[key]
		if ok {
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal.SetString(value)
			case reflect.Int:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					fieldVal.SetInt(intValue)
				}
			case reflect.Bool:
				boolValue := "yes" == value
				fieldVal.SetBool(boolValue)
			case reflect.Slice:
				if field.Type.Elem().Kind() == reflect.String {
					slice := strings.Split(value, ",")
					fieldVal.Set(reflect.ValueOf(slice))
				}
			default:

			}
		}
	}
	return config
}

func GetRunningMode() string {
	if properties.ClusterEnabled == "yes" {
		return ClusterMode
	} else {
		return StandaloneMode
	}
}
