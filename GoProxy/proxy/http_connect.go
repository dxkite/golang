// 基于HTTP代理的连接器
package proxy

import (
	"errors"
	"log"
	"net"
	"time"
)

type HTTPConnect struct {
	Proxy     string
	connector Connector
	Timeout   time.Duration
	wrapper   Wrapper
}

// 建立链接
func (c HTTPConnect) Dial(network, address string) (conn net.Conn, err error) {
	// 如果有下级
	if c.connector == nil {
		conn, err = net.DialTimeout(network, c.Proxy, c.Timeout)
		log.Println("create proxy to", c.Proxy)
	} else {
		conn, err = c.connector.Dial(network, address)
	}
	if err != nil {
		log.Println("dial error", err)
		return
	}
	if c.wrapper != nil {
		log.Println("use wrapper")
		conn = c.wrapper.Wrapper(conn)
	}
	log.Println("handshake ->", "CONNECT "+address)
	_, we := conn.Write([]byte("CONNECT " + address + " HTTP/1.1\r\n\r\n"))
	if we != nil {
		err = we
		return
	}
	data, err := readData(conn)
	// HTTP 协议握手
	code, msg := getRespond(data)
	if code != 200 {
		err = errors.New("http handshake error: " + msg)
	} else {
		log.Println("handshake <-", code, msg)
	}
	return conn, err
}

func (c HTTPConnect) NextConnect(connector Connector) Connector {
	c.connector = connector
	return c
}

func (c HTTPConnect) SetWrapper(wrapper Wrapper) Connector {
	c.wrapper = wrapper
	return c
}

// HTTP包装器
func NewHTTPConnect(proxy string, timeout time.Duration) HTTPConnect {
	return HTTPConnect{
		Proxy:   proxy,
		Timeout: timeout,
	}
}
