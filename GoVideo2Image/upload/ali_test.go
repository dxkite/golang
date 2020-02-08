package upload

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestAli_Upload(t *testing.T) {
	if data, err := ioutil.ReadFile("./test/cdn-video-0000.jpg"); err == nil {
		res, er := Upload("ali", &FileObject{
			Name: "cdn.jpg",
			Data: data,
		})
		if er != nil {
			t.Error(er)
		} else {
			fmt.Println("uploaded", res.Url)
		}
	} else {
		t.Error(err)
	}
}
