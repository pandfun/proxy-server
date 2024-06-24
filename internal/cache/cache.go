package cache

import "time"

// Struct/Node for the cache (linked list)
type CacheElement struct {
	Data    string
	URL     string
	LRUTime time.Time
	Next    *CacheElement
}

func findElement(url string) (*CacheElement, error) {
	// Find the element in the cache
	return nil, nil
}

func addElement(url string, data string) error {
	// Add the element to the cache
	return nil
}

func removeElement() error {
	// Remove the element from the cache
	return nil
}
