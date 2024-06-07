package core

import (
	"io"
	"syscall"
)

type Client struct {
	io.ReadWriter
	Fd     int
	cqueue RedisCmds
}

func (c Client) Write(b []byte) (int, error) {
	return syscall.Write(c.Fd, b)
}

func (c Client) Read(b []byte) (int, error) {
	return syscall.Read(c.Fd, b)
}


func NewClient(fd int) *Client {
	return &Client{
		Fd:     fd,
		cqueue: make(RedisCmds, 0),
	}
}
