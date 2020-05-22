package utils

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Conf struct {
	Addr           string
	MaxCli         uint32
	MaxWorker      uint32
	MaxHandleQueue uint32
	MaxSendQueue   uint32
	SendTimeOut    uint32
	LogPath        string
	LogLevel       uint32
}

func (c *Conf) Load(file string) {
	viper.SetConfigFile(strings.Split(filepath.Base(file), ".")[0])
	viper.AddConfigPath(filepath.Dir(file))
	if err := viper.ReadInConfig(); err != nil {
		c.SetDefalt()
	}
	cfg := viper.GetViper()
	if c.Addr = cfg.GetString("tcpServer.addr"); c.Addr == "" {
		c.Addr = DefaultAddr
	}
	if c.MaxCli = cfg.GetUint32("tcpServer.max_cli"); c.MaxCli == 0 {
		c.MaxCli = DefaultMaxCli
	}
	if c.MaxWorker = cfg.GetUint32("tcpServer.max_worker"); c.MaxWorker == 0 {
		c.MaxWorker = DefaultMaxWorker
	}
	if c.MaxHandleQueue = cfg.GetUint32("tcpServer.max_handle_queue"); c.MaxHandleQueue == 0 {
		c.MaxHandleQueue = DefaultMaxHandleQueue
	}
	if c.MaxSendQueue = cfg.GetUint32("tcpServer.max_send_queue"); c.MaxSendQueue == 0 {
		c.MaxSendQueue = DefaultMaxSendQueue
	}
	if c.SendTimeOut = cfg.GetUint32("tcpServer.send_timeout"); c.SendTimeOut == 0 {
		c.SendTimeOut = DefaultSendTimeOut
	}
	if c.LogPath = cfg.GetString("tcpServer.logpath"); c.LogPath == "" {
		c.LogPath = DefaultLogPath
	}
	if c.LogLevel = cfg.GetUint32("tcpServer.loglevel"); c.LogLevel == 0 {
		c.LogLevel = DefaultLogLevel
	}
}

func (c *Conf) SetDefalt() {
	c.Addr = DefaultAddr
	c.MaxCli = DefaultMaxCli
	c.MaxWorker = DefaultMaxWorker
	c.MaxHandleQueue = DefaultMaxHandleQueue
	c.MaxSendQueue = DefaultMaxSendQueue
	c.SendTimeOut = DefaultSendTimeOut
	c.LogPath = DefaultLogPath
	c.LogLevel = DefaultLogLevel
}
