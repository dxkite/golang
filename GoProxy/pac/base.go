package pac

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// 保存PAC文件
func WritePacResponse(writer io.Writer, pacFile, proxy string) (int, error) {
	data, err := ioutil.ReadFile(pacFile)
	if err != nil {
		return 0, err
	}
	var respond = "HTTP/1.1 200 OK\r\n"
	respond += "Content-Type: application/x-ns-proxy-autoconfig\r\n"
	pacTxt := strings.Replace(string(data), "__PROXY__", "PROXY "+proxy, -1)
	respond += fmt.Sprintf("Content-Length: %d\r\n", len(pacTxt))
	respond += "\r\n"
	respond += pacTxt
	return writer.Write([]byte(respond))
}

func warnError(fun func() (err error)) {
	if err := fun(); err != nil {
		log.Println(err)
	}
}
