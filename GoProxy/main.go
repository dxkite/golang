package main

import (
	config2 "dxkite.cn/GoProxy/config"
	"dxkite.cn/GoProxy/pac"
	"dxkite.cn/GoProxy/proxy"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	var filename = flag.String("conf", "conf/client.yml", "the config file")
	var help = flag.Bool("help", false, "the file name be input")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	config, err := config2.LoadConfig(*filename)

	if err != nil {
		log.Fatalln("read config file error", err)
	}

	var wrapper = proxy.NewXORWrapper(byte(config.XorKey))
	var timeout = time.Second * time.Duration(config.Timeout)

	if config.Timeout <= 0 {
		timeout = time.Second * 3
	}

	if len(config.RuntimeLog) > 0 {
		_ = os.MkdirAll(path.Dir(config.RuntimeLog), os.ModePerm)
		f, err := os.OpenFile(config.RuntimeLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Fatalf("error log file: %v", err)
		}
		defer func() { _ = f.Close() }()
		log.SetOutput(io.MultiWriter(os.Stderr, f))
	}

	if config.Mode == "client" {
		listen := proxy.GetRealProxy(config.Listen)
		log.Println("client mode start:", listen)
		log.Println("client mode pac:", "http://"+listen+"/pac.txt")
		if config.AutoPac {
			log.Println("enable auto pac", listen)
			go pac.AutoSetPac("http://"+listen+"/pac.txt", config.PacFileBackup)
		}
		proxy.StartHTTPWrapperConnect(config.Listen, nil,
			proxy.NewTLSConnect(config.Server, timeout).SetWrapper(wrapper))
	} else {
		log.Println("server mode start:", config.Listen)
		_, err := config2.LoadUserConfig(config.UserFile)
		if err != nil {
			log.Fatalln("read user config error", err)
		}
		var connector proxy.Connector
		if len(config.HTTPProxy) > 0 {
			connector = proxy.NewHTTPConnect(config.HTTPProxy, timeout)
		}
		proxy.StartTLSWrapperConnectListen(config.Listen, wrapper,
			connector,
			proxy.NewTLSListen(config.CertFile, config.KeyFile))
	}
}
