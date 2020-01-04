package main

import (
	"dxkite.cn/GoProxy/config"
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
	var filename = flag.String("conf", "conf/client.yml", "the conf file")
	var help = flag.Bool("help", false, "the file name be input")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	conf, err := config.LoadConfig(*filename)

	if err != nil {
		log.Fatalln("read conf file error", err)
	}

	var wrapper = proxy.NewXORWrapper(byte(conf.XorKey))
	var timeout = time.Second * time.Duration(conf.Timeout)

	if conf.Timeout <= 0 {
		timeout = time.Second * 3
	}

	if len(conf.RuntimeLog) > 0 {
		_ = os.MkdirAll(path.Dir(conf.RuntimeLog), os.ModePerm)
		f, err := os.OpenFile(conf.RuntimeLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Fatalf("error log file: %v", err)
		}
		defer func() { _ = f.Close() }()
		log.SetOutput(io.MultiWriter(os.Stderr, f))
	}

	if conf.Mode == "client" {
		listen := proxy.GetRealProxy(conf.Listen)
		log.Println("client mode start:", listen)
		log.Println("client mode pac:", "http://"+listen+"/pac.txt")
		if conf.AutoPac {
			log.Println("enable auto pac", listen)
			go pac.AutoSetPac("http://"+listen+"/pac.txt", conf.PacFileBackup)
		}
		proxy.StartHTTPWrapperConnect(conf.Listen, nil,
			proxy.NewTLSConnect(conf.Server, timeout).SetWrapper(wrapper))
	} else {
		log.Println("server mode start:", conf.Listen)
		_, err := config.LoadUserConfig(conf.UserFile)
		if err != nil {
			log.Fatalln("read user conf error", err)
		}
		var connector proxy.Connector
		if len(conf.HTTPProxy) > 0 {
			connector = proxy.NewHTTPConnect(conf.HTTPProxy, timeout)
		}
		proxy.StartTLSWrapperConnectListen(conf.Listen, wrapper,
			connector,
			proxy.NewTLSListen(conf.CertFile, conf.KeyFile))
	}
}
