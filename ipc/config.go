package ipc

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	ErrEmptyGroupName = errors.New("empty group name")
)

type Config struct {
	GroupName     string
	Mode          Mode
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	BufferSize    int
	IgnoreDialErr bool
	EnableDebug   bool
}

func (c *Config) validate() error {
	c.GroupName = strings.TrimSpace(c.GroupName)
	if c.GroupName == "" {
		return ErrEmptyGroupName
	}
	if c.Mode == UnknownMode {
		c.Mode = Dual
	}
	if int64(c.ReadTimeout) == 0 {
		c.ReadTimeout = 3 * time.Second
	}
	if int64(c.WriteTimeout) == 0 {
		c.WriteTimeout = 3 * time.Second
	}
	if c.BufferSize == 0 {
		c.BufferSize = 16
	}
	return nil
}

func (c *Config) setReadDeadline(conn net.Conn) error {
	return conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
}

func (c *Config) setWriteDeadline(conn net.Conn) error {
	return conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
}

func (c *Config) log(format string, a ...any) {
	if !c.EnableDebug {
		return
	}
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%v %v\n", time.Now().Format("[2006-01-02 15:04:05.999999]"), msg)
}
