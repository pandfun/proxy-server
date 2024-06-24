package proxy

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ProxyServer struct {
	addr string
}

func NewProxyServer(addr string) *ProxyServer {
	return &ProxyServer{
		addr: addr,
	}
}

func (p *ProxyServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/proxy/").Subrouter()

	RegisterRoutes(subrouter)

	log.Println("Proxy: Listening on", p.addr)

	return http.ListenAndServe(p.addr, router)
}
