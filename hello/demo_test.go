package hello

import (
	"encoding/hex"
	"log"
	"testing"
)

func init()  {
	log.Println("demo_test.go init")
}

func TestHello(t *testing.T) {
	log.Println("test run")
	log.Println("\n"+hex.Dump([]byte("main.go initmain.go initmain.go initmain.go initmain.g\no initmain.go initmain.go initmain.go init")))

}