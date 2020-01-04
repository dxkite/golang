// +build windows

package pac

import (
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func AutoSetPac(pacUri,pacBackFile string) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer warnError(k.Close)
	configUrl, _, err := k.GetStringValue("AutoConfigURL");
	var exist = true
	if err != nil {
		exist = false
		if err != registry.ErrNotExist {
			log.Fatal(err)
		}
	}
	if exist {
		log.Println("got raw pac", configUrl)
		if err := ioutil.WriteFile(pacBackFile, []byte(configUrl), os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	if err := k.SetStringValue("AutoConfigURL", pacUri); err != nil {
		log.Fatal(err)
	}
	log.Println("set AutoConfigURL", pacUri)
	<-signals
	log.Println("recover AutoConfigURL")
	if exist {
		if err := k.SetStringValue("AutoConfigURL", configUrl); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := k.DeleteValue("AutoConfigURL"); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("auto config finish")
	os.Exit(0)
}
