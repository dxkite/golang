package proxy

import (
	"log"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func TestNewTLSListen(t *testing.T) {
	StartHTTPWrapperConnectListen(":8080",
		nil,
		NewHTTPConnect("127.0.0.1:12639", 3*time.Second),
		NewTLSListen("../conf/server.pem", "../conf/server.key"))

	//var xor = 0x11
	//// STL代理
	//StartHTTPWrapperConnectListen(":8080",
	//	NewXORWrapper(byte(xor)),
	//	NewHTTPConnect("127.0.0.1:12639", 3*time.Second),
	//	NewTLSListen("conf/server.pem", "conf/server.key"))

}

func TestNewTSLConnect(t *testing.T) {

	// 代理链接到加密代理
	StartHTTPWrapperConnect(":8888", nil,
		NewTLSConnect(":8080", 3*time.Second))

	//var xor = 0x11
	//// 代理链接到加密代理
	//StartHTTPWrapperConnect(":8888", nil,
	//	NewTLSConnect(":8080", 3*time.Second).SetWrapper(NewXORWrapper(byte(xor))))
}
