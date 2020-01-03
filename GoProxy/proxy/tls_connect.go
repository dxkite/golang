package proxy

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
	"time"
)

type ConnectMessage struct {
	Host         string   `json:"host"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	HardwareAddr []string `json:"mac_addr"`
}

type Message struct {
	Code int         `json:"err_code"`
	Msg  string      `json:"err_msg"`
	Data interface{} `json:"data"`
}

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
		log.Println("tls dial error", err)
		return
	}
	if c.wrapper != nil {
		conn = c.wrapper.Wrapper(conn)
		log.Println("use wrapper")
	}
	var macAddr []string
	if mac, er := getMac(); er == nil {
		macAddr = mac
	}

	if _, we := sendJsonMsg(conn, ConnectMessage{
		Host:         address,
		Username:     _config.Username,
		Password:     _config.Password,
		HardwareAddr: macAddr,
	}); we != nil {
		err = we
		return
	}

	data, re := readData(conn)
	var msg Message
	if re == nil {
		if er := json.Unmarshal(data, &msg); er != nil {
			err = er
			return
		}
	} else {
		err = re
		return
	}
	if msg.Code != 0 {
		err = errors.New("tls handshake error: " + msg.Msg)
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

// 处理tls
func tlsTunnel(conn net.Conn, dial DialFunc) (mac, host string, up, down int64, err error) {
	defer warnError(conn.Close)
	data, re := readData(conn)
	if re != nil {
		if _, we := sendJsonMsg(conn, Message{Code: 400, Msg: re.Error()}); we != nil {
			err = errors.New(fmt.Sprintf("%v: %v", re, we))
		} else {
			err = we
		}
		return
	}
	var connect ConnectMessage
	if er := json.Unmarshal(data, &connect); er != nil {
		if _, we := sendJsonMsg(conn, Message{Code: 400, Msg: er.Error()}); we != nil {
			err = errors.New(fmt.Sprintf("%v: %v", re, we))
		} else {
			err = we
		}
		return
	}
	sort.Strings(connect.HardwareAddr)
	mac = strings.Join(connect.HardwareAddr, ",")
	if _config.Auth == true {
		log.Println("auth enable")
		if er := AuthUser(connect.Username, connect.Password, connect.HardwareAddr); er != nil {
			if _, we := sendJsonMsg(conn, Message{Code: 401, Msg: er.Error()}); we != nil {
				err = we
				return
			}
			err = er
			return
		}
		log.Println("user", connect.Username, "login", mac)
	}
	host = connect.Host
	to, de := dial("tcp", host)
	log.Println("dial", host)
	if to != nil {
		defer warnError(to.Close)
	}
	if de != nil {
		if _, we := sendJsonMsg(conn, Message{Code: 502, Msg: "connection error: " + de.Error()}); we != nil {
			err = errors.New(fmt.Sprintf("%v: %v", de, we))
		} else {
			err = de
		}
		return
	}
	if _, we := sendJsonMsg(conn, Message{Code: 0, Msg: "ok"}); we != nil {
		err = we
		return
	}
	log.Println("make tunnel", host)
	up, down = makeTunnel(conn, to)
	return
}

func sendJsonMsg(conn net.Conn, message interface{}) (n int, err error) {
	data, er := json.Marshal(message)
	if er != nil {
		log.Println("json encode error:", err)
		err = er
		return
	}
	return conn.Write(data)
}

// 开启HTTP代理服务器
func StartTLS(address string) {
	startServe(tlsTunnel, nil, address, nil, nil)
}

// 开启HTTP代理服务器
func StartTLSWrapper(address string, wrapper Wrapper) {
	startServe(tlsTunnel, nil, address, wrapper, nil)
}

// 开启HTTP代理服务器
func StartTLSWrapperConnect(address string, wrapper Wrapper, connector Connector) {
	startServe(tlsTunnel, nil, address, wrapper, connector)
}

// 开启HTTP代理服务器
func StartTLSWrapperConnectListen(address string, wrapper Wrapper, connector Connector, listenFunc ListenFunc) {
	startServe(tlsTunnel, listenFunc, address, wrapper, connector)
}
