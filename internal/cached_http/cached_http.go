package cached_http

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/amirulabu/pokemon-store-backend/internal/database"
)

var cacheMutex sync.Mutex

func CacheAndRetrieve(url string, db *database.DB) ([]byte, error) {
	// Check if the response is already cached
	cachedResponse, err := db.GetCachedResponse(url)
	if err == nil {
		fmt.Println("Data retrieved from cache.")
		return cachedResponse.Payload, nil
	}

	// Fetch data from the API
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response payload
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Insert the response into the cache
	cacheMutex.Lock()
	_, err = db.InsertCachedResponse(url, payload)
	cacheMutex.Unlock()

	if err != nil {
		return nil, err
	}

	fmt.Println("Data cached and retrieved successfully.")
	return payload, nil
}
