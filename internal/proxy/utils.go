package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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

// func handleResponseBody(body io.Reader) ([]byte, error) {

//     reader, err := charset.NewReader(body, "text/html")
//     if err != nil {
//         return nil, err
//     }

//     bodyBytes, err := io.ReadAll(reader)
//     if err != nil {
//         return nil, err
//     }

// 	log.Println("HTML Content: ", (bodyBytes))

//     return bodyBytes, nil
// }

func rewriteURL(urlStr, proxyHost, baseURL string, scheme string) string {

	if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
		return fmt.Sprintf("http://%s/?url=%s", proxyHost, url.QueryEscape(urlStr))
	}

	if strings.HasPrefix(urlStr, "//") {
		return fmt.Sprintf("http://%s/?url=%s", proxyHost, url.QueryEscape(scheme+":"+urlStr))
	}

	if strings.HasPrefix(urlStr, "/") {
		return fmt.Sprintf("http://%s/?url=%s", proxyHost, url.QueryEscape(baseURL+urlStr))
	}

	return fmt.Sprintf("http://%s/?url=%s", proxyHost, url.QueryEscape(baseURL+"/"+urlStr))
}

func handleRelativeLinks(body []byte, proxyHost, originalURL string) ([]byte, error) {

	bodyStr := string(body[:])
	// log.Println("HANDLING LINKS ON ", bodyStr)

	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return nil, err
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	re := regexp.MustCompile(`(href|src)="([^"]+)"`)

	bodyStr = re.ReplaceAllStringFunc(bodyStr, func(match string) string {

		matches := re.FindStringSubmatch(match)

		if len(matches) > 2 {
			attr := matches[1]
			oldURL := matches[2]
			newURL := rewriteURL(oldURL, proxyHost, baseURL, parsedURL.Scheme)
			return fmt.Sprintf(`%s="%s"`, attr, newURL)
		}

		return match
	})

	return []byte(bodyStr), nil
}
