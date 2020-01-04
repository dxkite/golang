// 基于HTTP代理的连接器
package proxy

import (
	"bufio"
	"bytes"
	"dxkite.cn/GoProxy/config"
	"dxkite.cn/GoProxy/pac"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HTTPConnect struct {
	Proxy     string
	connector Connector
	Timeout   time.Duration
	wrapper   Wrapper
	User      string
	Password  string
	Mac       bool
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

	var macAddr []string
	if c.Mac == true {
		if mac, er := GetMac(); er == nil {
			macAddr = mac
		}
	}
	sort.Strings(macAddr)
	log.Println("handshake ->", "CONNECT "+address)
	_, we := conn.Write(createRequest(address, c.User, c.Password, strings.Join(macAddr, ",")))
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

func getRespond(b []byte) (code int, message string) {
	s := string(b)
	lr := strings.Index(s, "\r\n")
	h := s[:lr]
	t := strings.SplitN(h, " ", 3)
	code, _ = strconv.Atoi(t[1])
	message = t[2]
	return
}

func getHost(host string) string {
	if strings.Index(host, ":") > 0 {
		return host
	}
	return host + ":80"
}

// 处理HTTP代理协议
func httpTunnel(conn net.Conn, dial DialFunc) (mac, host string, up, down int64, err error) {
	defer warnError(conn.Close)
	data, re := readData(conn)
	if re != nil {
		if _, we := conn.Write([]byte(fmt.Sprintf("HTTP/1.1 400 Bad Request\r\nContent-Length: %d\r\n\r\n%v", len(re.Error()), re))); we != nil {
			err = errors.New(fmt.Sprintf("%v: %v", re, we))
		} else {
			err = we
		}
		return
	}
	request, er := http.ReadRequest(bufio.NewReader(bytes.NewReader(data)))
	if er != nil {
		err = er
		return
	}
	host = getHost(request.Host)
	if host == conn.LocalAddr().String() {
		log.Println("request self as http, respond as pac file", request.Method, request.URL.Path)
		if n, err :=pac.WritePacResponse(conn, config.GetConfig().PacFile, conn.LocalAddr().String()); err != nil {
			log.Println("request err", err)
		} else {
			down = int64(n)
		}
		warnError(conn.Close)
		return
	}
	to, de := dial("tcp", host)
	log.Println("dial", host)
	if to != nil {
		defer warnError(to.Close)
	}
	if de != nil {
		if _, we := conn.Write([]byte(fmt.Sprintf("HTTP/1.1 502 Bad Gateway\r\nContent-Length: %d\r\n\r\n%v", len(de.Error()), de))); we != nil {
			err = errors.New(fmt.Sprintf("%v: %v", we, de))
		} else {
			err = de
		}
		return
	}
	mac = request.Header.Get("Mac")
	if config.GetConfig().Auth == true {
		log.Println("auth enable")
		if user, pass, ok := ProxyAuth(request); ok {
			log.Println("auth enable", user, pass)
			if er := config.AuthUser(user, pass, strings.Split(mac, ",")); er != nil {
				if _, we := conn.Write([]byte(fmt.Sprintf("406 Not Acceptable\r\nContent-Length: %d\r\n\r\n%v", len(er.Error()), er))); we != nil {
					err = we
				}
				return
			}
			log.Println("user", user, "login", mac)
		} else {
			err = errors.New("basic auth error")
			if _, we := conn.Write([]byte(fmt.Sprintf("406 Not Acceptable\r\nContent-Length: %d\r\n\r\n%v", len(err.Error()), err))); we != nil {
				err = errors.New(fmt.Sprintf("%v: %v", err, we))
			}
			return
		}
	}
	if request.Method == "CONNECT" {
		if _, we := conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); we != nil {
			err = we
			log.Println(we)
			warnError(conn.Close)
			warnError(to.Close)
			return
		}
	} else {
		if n, we := to.Write(data); we != nil {
			err = we
			up += int64(n)
			return
		} else {
			up += int64(n)
		}
	}
	log.Println("make tunnel", host)
	upr, down := makeTunnel(conn, to)
	up += upr
	return
}

// 开启HTTP代理服务器
func StartHTTP(address string) {
	startServe(httpTunnel, nil, address, nil, nil)
}

// 开启HTTP代理服务器
func StartHTTPWrapper(address string, wrapper Wrapper) {
	startServe(httpTunnel, nil, address, wrapper, nil)
}

// 开启HTTP代理服务器
func StartHTTPWrapperConnect(address string, wrapper Wrapper, connector Connector) {
	startServe(httpTunnel, nil, address, wrapper, connector)
}

// 开启HTTP代理服务器
func StartHTTPWrapperConnectListen(address string, wrapper Wrapper, connector Connector, listenFunc ListenFunc) {
	startServe(httpTunnel, listenFunc, address, wrapper, connector)
}

// 创建请求
func createRequest(host, username, password, mac string) []byte {
	request := "CONNECT " + host + " HTTP/1.1\r\n"
	if len(mac) > 0 {
		request += "Mac: " + mac + "\r\n"
	}
	if len(username) > 0 {
		credentials := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		request += "Proxy-Authorization: Basic " + credentials + "\r\n"
	}
	return []byte(request + "\r\n")
}

func ProxyAuth(r *http.Request) (username, password string, ok bool) {
	auth := r.Header.Get("Proxy-Authorization")
	if auth == "" {
		return
	}
	return parseBasicAuth(auth)
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
