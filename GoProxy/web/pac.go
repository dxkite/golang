package web

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strings"
)

func StartPacServe(address, pacFile, proxy string) {
	r := gin.Default()
	r.GET("/pac.txt", func(c *gin.Context) {
		data, err := ioutil.ReadFile(pacFile)
		if err != nil {
			log.Fatal(err)
		}
		pacTxt := strings.Replace(string(data), "__PROXY__", proxy, -1)
		if _, err := c.Writer.Write([]byte(pacTxt)); err != nil {
			log.Println(err)
		}
	})
	if err := r.Run(address); err != nil {
		log.Println(err)
	}
}
