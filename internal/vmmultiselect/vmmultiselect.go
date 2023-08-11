package vmmultiselect

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/nickmy9729/vmmultiselect/config"
	"github.com/nickmy9729/vmmultiselect/internal/api"
)

func ExecuteVMSelect(cfg *config.Config, groupName string, headers http.Header) ([]byte, error) {
	endpoints, exists := cfg.Groups[groupName]
	if !exists {
		return nil, fmt.Errorf("group not found")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var combinedData []byte

	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			data, err := api.MakeAPIRequest(endpoint, headers)
			if err != nil {
				// Handle error
				return
			}
			mu.Lock()
			combinedData = append(combinedData, data...)
			mu.Unlock()
		}(endpoint)
	}

	wg.Wait()
	return combinedData, nil
}

