package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/pandfun/proxy-server/utils"
)

// Handler for all invalid paths
func invalidPathHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("New incoming request", r.URL)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Invalid path"})
}


// Handle incoming proxy request
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	targetUrl := r.URL.Query().Get("url")
	if targetUrl == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("url not found in the request"))
		return
	}

	log.Print("Incoming request ", targetUrl)

	// parse the given url to check validity
	parsedUrl, err := url.Parse(targetUrl)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to parse the url"))
		log.Print("Failed to parse the url", targetUrl)
		return
	}

	// If the url has no scheme or host, it's invalid
	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid url"))
		return
	}

	// get the data from the remote web server
	resp, err := http.Get(targetUrl)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to make get request to target url"))
		log.Print("Failed to make get request", targetUrl)
		return
	}

	defer resp.Body.Close()

	utils.WriteResponse(w, resp)
}
