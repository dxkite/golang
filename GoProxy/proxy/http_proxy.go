// 实现HTTP代理服务器
package proxy

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const maxLength = 8 * 1024 * 1024 // 8M

type DialFunc func(network, address string) (net.Conn, error)
type ListenFunc func(network, address string) (net.Listener, error)

// 开启HTTP代理服务器
func StartHTTP(address string) {
	startHTTPServe(nil, address, nil, nil)
}

// 开启HTTP代理服务器
func StartHTTPWrapper(address string, wrapper Wrapper) {
	startHTTPServe(nil, address, wrapper, nil)
}

// 开启HTTP代理服务器
func StartHTTPWrapperConnect(address string, wrapper Wrapper, connector Connector) {
	startHTTPServe(nil, address, wrapper, connector)
}

// 开启HTTP代理服务器
func StartHTTPWrapperConnectListen(address string, wrapper Wrapper, connector Connector, listenFunc ListenFunc) {
	startHTTPServe(listenFunc, address, wrapper, connector)
}

// 开启代理
func startHTTPServe(listen ListenFunc, address string, wrapper Wrapper, connector Connector) {
	if listen == nil {
		listen = net.Listen
	}
	listener, err := listen("tcp", address)
	defer listener.Close()

	if err != nil {
		log.Fatal("create http proxy error", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		var dial DialFunc
		if connector == nil {
			dial = net.Dial
		} else {
			dial = connector.Dial
			log.Println("use connector dialer")
		}
		if wrapper == nil {
			go parseHTTPConnect(conn, dial)
		} else {
			log.Println("use wrapper")
			wrapped := wrapper.Wrapper(conn)
			go parseHTTPConnect(wrapped, dial)
		}
	}
}

// 处理HTTP代理协议
func parseHTTPConnect(conn net.Conn, dial DialFunc) {
	data, err := readData(conn)
	if err != nil {
		log.Println("read error", err)
		if _, err := conn.Write([]byte(fmt.Sprintf("HTTP/1.1 400 Bad Request\r\n\r\n%s", err))); err != nil {
			log.Println(err)
		}
		return
	}
	request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(data)))
	if err != nil {
		log.Println("read header error", err)
		log.Println("header\n" + hex.Dump(data))
		return
	}
	host := getHost(request.Host)
	log.Println("connect", host)
	to, err := dial("tcp", host)
	if err != nil {
		log.Println("create tunnel error", err)
		if _, err := conn.Write([]byte(fmt.Sprintf("HTTP/1.1 502 Bad Gateway\r\n\r\n%s", err))); err != nil {
			log.Println(err)
		}
		return
	}
	if request.Method == "CONNECT" {
		if _, err := conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
			log.Println(err)
		}
	} else {
		//if _, err := conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n\r\n")); err != nil {
		//	log.Println(err)
		//}
		//log.Println("method not allowed", request.Method, request.Host)
		log.Println("raw http")
		if _, err := to.Write(data); err != nil {
			log.Println(err)
		}
	}
	makeTunnel(conn, to)
	log.Println("connection closed", request.Host)
}

// 构建隧道
func makeTunnel(src, dest io.ReadWriteCloser) {
	var wg sync.WaitGroup
	// src.Write(dest.Read())
	go func() {
		defer src.Close()
		defer dest.Close()
		defer wg.Done()

		wg.Add(1)
		if _, err := io.Copy(src, dest); err != nil {
			log.Println("makeTunnel:", err)
		}
	}()
	// dest.Write(src.Read())
	go func() {
		defer src.Close()
		defer dest.Close()
		defer wg.Done()

		wg.Add(1)
		if _, err := io.Copy(dest, src); err != nil {
			log.Println("makeTunnel:", err)
		}
	}()
	wg.Wait()
}

// 复制数据
func readData(reader io.Reader) (d []byte, err error) {
	var size = 1024
	d = make([]byte, 0)
	buf := make([]byte, size)
	read := 0
	for {
		n, re := reader.Read(buf)
		if re != nil {
			err = re
			return
		}
		d = append(d, buf[:n]...)
		read += n
		if n < size || read >= maxLength {
			return
		}
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
