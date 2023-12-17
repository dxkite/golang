package upload

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestAli_Upload(t *testing.T) {
	if data, err := ioutil.ReadFile("./test/1.png"); err == nil {
		res, er := Upload("ali", &FileObject{
			Name: "upload.png",
			Data: data,
		})
		if er != nil {
			t.Error("uploaded real image error", er)
		} else {
			fmt.Println("uploaded real image", res.Url)
		}
		if data, err := ioutil.ReadFile("./test/2-fake.png"); err == nil {
			res, er := Upload("ali", &FileObject{
				Name: "upload.png",
				Data: data,
			})
			if er != nil {
				t.Error("uploaded fake image error", er)
			} else {
				fmt.Println("uploaded fake image", res.Url)
			}
		} else {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}
}
