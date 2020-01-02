package proxy

import (
	"net"
)

type XORWrapper struct {
	Value byte
}

type wrappedXORConnect struct {
	net.Conn
	Value byte
}

// 写包装
func (c wrappedXORConnect) Read(b []byte) (n int, err error) {
	n, re := c.Conn.Read(b)
	if re != nil {
		err = re
		return
	}
	for i, v := range b {
		b[i] = v ^ c.Value
	}
	return n, err
}

// 读包装
func (c wrappedXORConnect) Write(b []byte) (n int, err error) {
	for i, v := range b {
		b[i] = v ^ c.Value
	}
	return c.Conn.Write(b)
}

// 包装
func (c XORWrapper) Wrapper(conn net.Conn) net.Conn {
	return wrappedXORConnect{
		Conn:  conn,
		Value: c.Value,
	}
}

// 创建包装器
func NewXORWrapper(xor byte) Wrapper {
	return XORWrapper{
		Value: xor,
	}
}
