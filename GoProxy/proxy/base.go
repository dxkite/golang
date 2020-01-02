package proxy

import (
	"net"
)

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

type Wrapper interface {
	Wrapper(conn net.Conn) net.Conn
}
