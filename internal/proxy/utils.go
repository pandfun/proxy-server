package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)


func validateURL(rawURL string) (int, error) {

	parsedUrl, err := url.Parse(rawURL)
	if err != nil {
        log.Println("Failed to validate URL : ", rawURL)
		return http.StatusInternalServerError, err
	}

	// If the url has no scheme or host, it's invalid
	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return http.StatusBadRequest, fmt.Errorf("url has no scheme or host")
	}

    return http.StatusOK, nil
}


func copyHeader(dstHeader, srcHeader http.Header) {

    for key, values := range srcHeader {
        for _, headerValue := range values {

            dstHeader.Add(key, headerValue)
        }
    }
}
