package proxy

import (
	"encoding/hex"
	"io"
	"net"
	"strings"
)

const maxLength = 8 * 1024 * 1024 // 8M

type DialFunc func(network, address string) (net.Conn, error)
type ListenFunc func(network, address string) (net.Listener, error)
type TunnelFunc func(conn net.Conn, dial DialFunc) (mac, host string, up, down int64, err error)

// 网络连接
// 用来做代理与网络数据交互的部分
type Connector interface {
	// 请求
	Dial(network, address string) (conn net.Conn, err error)
	// 下一个连接
	NextConnect(connector Connector) Connector
	// 设置包装器
	SetWrapper(wrapper Wrapper) Connector
}

// 网络包装器
type Wrapper interface {
	Wrapper(conn net.Conn) net.Conn
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

// 获取MAC地址
func getMac() (addr []string, err error) {
	if it, ie := net.Interfaces(); ie == nil {
		for _, itr := range it {
			ad := strings.ToUpper(hex.EncodeToString(itr.HardwareAddr))
			if len(ad) > 0 {
				addr = append(addr, ad)
			}
		}
	} else {
		err = ie
	}
	return
}
