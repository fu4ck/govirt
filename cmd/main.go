package main

import (
	"flag"
	"govirt/runtime"
	"govirt/server"
	"log"
)

var (
	// 配置文件路径 --config /path/path  或者 -c /path/path
	cfg = flag.String("config", "./conf/govirt.conf", "virt config file path.")
)

func main() {

	flag.Parse()

	// init config
	cfg, err := runtime.InitConfig(*cfg)
	if err != nil {
		log.Fatalf("init config error, 【error】: %+v", err)
	}

	// new go virt service
	govirt := runtime.NewGoVirt(cfg)

	server.StartRouter(govirt)
}
