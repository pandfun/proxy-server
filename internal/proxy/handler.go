package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/pandfun/proxy-server/internal/utils"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/test", testHandler).Methods("GET")
	router.HandleFunc("/check", proxyHandler).Methods("GET")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Hello, World!"})
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// get the URL from the path
	rawURL := r.URL.Query().Get("url")
	if rawURL == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid url"))
		return
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse url/invalid url"))
		return
	}

	url := parsedURL.String()
	log.Println("incoming request for", url)

	// make a request to the URL
	res, err := http.Get(url)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Println(res)

	// return the response
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}
