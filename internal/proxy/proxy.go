package proxy

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Proxy struct {
	addr string
}

// Proxy server constructor
func NewProxyServer(addr string) *Proxy {
	return &Proxy{addr: addr}
}

// Start the proxy server
func (p *Proxy) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/proxy", proxyHandler)
	router.HandleFunc("/", invalidPathHandler)

	log.Printf("Server is running on %v", p.addr)
	
	return http.ListenAndServe(p.addr, router)
}
