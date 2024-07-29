package main

import (
	"log"

	"github.com/pandfun/proxy-server/internal/proxy"
)

func main() {
	proxyServer := proxy.NewProxyServer(":9000")
	
	if err := proxyServer.Run(); err != nil {
		log.Fatal("error in proxy server", err.Error())
	}
}