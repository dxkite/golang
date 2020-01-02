package proxy

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"time"
)

func NewTLSListen(certFile, keyFile string) ListenFunc {
	return func(network, address string) (listener net.Listener, e error) {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Println(err)
			return
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		return tls.Listen(network, address, config)
	}
}

type TLSConnect struct {
	Address string
	Timeout time.Duration
	wrapper Wrapper
}

func NewTLSConnect(address string, timeout time.Duration) TLSConnect {
	return TLSConnect{
		Address: address,
		Timeout: timeout,
	}
}

// 建立链接
func (c TLSConnect) Dial(network, address string) (conn net.Conn, err error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err = tls.Dial(network, c.Address, conf)
	if err != nil {
		log.Println("stl dial error", err)
		return
	}
	if c.wrapper != nil {
		conn = c.wrapper.Wrapper(conn)
		log.Println("use wrapper")
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
	return
}

func (c TLSConnect) NextConnect(connector Connector) Connector {
	return c
}

func (c TLSConnect) SetWrapper(wrapper Wrapper) Connector {
	c.wrapper = wrapper
	return c
}
