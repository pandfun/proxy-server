package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"

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

	status, err := validateURL(targetUrl)
	if err != nil {
		utils.WriteError(w, status, err)
	}

	proxyReq, err := http.NewRequest(r.Method, targetUrl, r.Body)
	if err != nil {
		log.Println("Failed to create a new request")
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create a new request"))
		return
	}

	copyHeader(proxyReq.Header, r.Header)

	client := &http.Client{}


	resp, err := client.Do(proxyReq)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	defer resp.Body.Close()


	copyHeader(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)


	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	w.Write(respBody)

	// if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
	// 	w.Write(respBody)
	// 	return
	// }

	// modifiedBody := handleRelativeLinks(respBody, r.Host, targetUrl)
	// w.Write(modifiedBody)
}
