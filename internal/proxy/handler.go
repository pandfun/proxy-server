package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pandfun/proxy-server/internal/cache"
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

	// Check the it's cached
	if elem, ok := cache.LRU.Get(targetUrl); ok {

		log.Println("Resource already in cache: ", targetUrl)

		// If the reponse has expired
		if time.Now().After(elem.Expiration) {

			log.Println("Remove resource: already expired")
			cache.LRU.Remove(targetUrl)
		} else {

			log.Println("Returing cached data for ", targetUrl)

			for key, value := range elem.Headers {
				w.Header()[key] = value
			}

			w.WriteHeader(http.StatusOK)
			w.Write(elem.Value)
			return
		}
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

	// If not a html response, no need to process it further
	if !(strings.Contains(resp.Header.Get("Content-Type"), "text/html")) {

		if ok, exp := checkCanCache(resp); ok {
			log.Println("New cache entry for: ", targetUrl)
			cache.LRU.Set(targetUrl, []byte(respBody), resp.Header, exp)
		} else if !ok {
			log.Println("Not caching data for ", targetUrl)
		}

		w.Write(respBody)
		return
	}

	modifiedBody, err := handleRelativeLinks(respBody, r.Host, targetUrl)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if ok, exp := checkCanCache(resp); ok {
		log.Println("New cache entry for: ", targetUrl)
		cache.LRU.Set(targetUrl, []byte(modifiedBody), resp.Header, exp)
	} else if !ok {
		log.Println("Not caching data for ", targetUrl)
	}

	w.Write(modifiedBody)
}

func checkCanCache(resp *http.Response) (bool, time.Time) {

	cacheControl := resp.Header.Get("Cache-Control")
	if cacheControl != "" {

		if strings.Contains(cacheControl, "no-store") {
			log.Println("Found 'no-store' in cache control")
			return false, time.Time{}
		}

		if strings.Contains(cacheControl, "no-cache") {
			log.Println("Found 'no-cache' in cache control")
			return false, time.Time{}
		}

		if strings.Contains(cacheControl, "private") {
			log.Println("Found 'private' in cache control")
			return false, time.Time{}
		}

		if strings.Contains(cacheControl, "max-age=") {
			parts := strings.Split(cacheControl, "max-age=")
			if len(parts) >= 1 {
				maxAge, err := strconv.Atoi(strings.Split(parts[1], ",")[0])

				if err != nil {
					return true, time.Now().Add(time.Duration(maxAge) * time.Second)
				}
			}
		}
	}

	expires := resp.Header.Get("Expires")
	if expires != "" {
		expirationTime, err := time.Parse(time.RFC1123, expires)
		if err != nil {
			return true, expirationTime
		}
	}

	// If no cache directives found, make the element expire in 24hrs
	return true, time.Now().Add(24 * time.Hour)
}
