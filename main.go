package main

import (
	"flag"

	"github.com/zhjx922/zcache/server"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "11211", "您要监听的端口号是啥?")
	flag.Parse()
	server.Start(port)
}
