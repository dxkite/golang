// 实现HTTP代理服务器
package proxy

import (
	"io"
	"log"
	"net"
)

// 开启代理
func startServe(tunnel TunnelFunc, listen ListenFunc, address string, wrapper Wrapper, connector Connector) {
	if listen == nil {
		listen = net.Listen
	}
	listener, err := listen("tcp", address)
	defer warnError(listener.Close)

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
			go runTunnel(conn, dial, tunnel)
		} else {
			log.Println("use wrapper")
			wrapped := wrapper.Wrapper(conn)
			go runTunnel(wrapped, dial, tunnel)
		}
	}
}

func runTunnel(conn net.Conn, dial DialFunc, tunnelFunc TunnelFunc) {
	if mac, host, up, down, err := tunnelFunc(conn, dial); err == nil {
		log.Printf("connection closed: %s %s ↑%db, ↓%db", mac, host, up, down)
	} else {
		log.Printf("connection closed: %s %s ↑%db, ↓%db", mac, host, up, down)
		log.Printf("connection error: %v", err)
	}
	return
}

// 构建隧道
func makeTunnel(src, dest io.ReadWriteCloser) (up, down int64) {
	// src.Write(dest.Read())
	var _up = make(chan int64)
	var _down = make(chan int64)

	go func() {
		defer warnError(src.Close)
		defer warnError(dest.Close)
		if n, err := io.Copy(src, dest); err != nil {
			log.Println("makeTunnel:", err)
			_down <- n
		} else {
			_down <- n
		}
	}()
	// dest.Write(src.Read())
	go func() {
		defer warnError(src.Close)
		defer warnError(dest.Close)
		if n, err := io.Copy(dest, src); err != nil {
			log.Println("makeTunnel:", err)
			_up <- n
		} else {
			_up <- n
		}
	}()
	up = <-_up
	down = <-_down
	return
}

func warnError(fun func() (err error)) {
	if err := fun(); err != nil {
		log.Println(err)
	}
}
