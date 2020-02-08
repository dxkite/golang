package hello

import (
	"fmt"
	"log"
)

func init()  {
	log.Println("demo.go init")
}

func Hello() {
	fmt.Println("Hello, dxkite")
}
