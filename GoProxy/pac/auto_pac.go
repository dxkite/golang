// +build !windows

package pac

import (
	"log"
)

func AutoSetPac(pacFile, proxy string) {
	log.Fatalln("auto set pac only support windows")
}
