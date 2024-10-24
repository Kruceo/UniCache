package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicache/memory"
	"unicache/utils"
)

var excludePaths []string = strings.Split(utils.GetEnvOrDefault("EXCLUDED_PATHS", ""), ",")

func main() {

	cacheTimeout, err := strconv.ParseInt(utils.GetEnvOrDefault("CACHE_TIMEOUT", "120"), 10, 32)

	if err != nil {
		log.Printf("Not possible parse cache timeout (%s): %v", utils.GetEnvOrDefault("CACHE_TIMEOUT", "120"), err)
		return
	}

	fmt.Printf("Cache timeout defined to %dS\n", cacheTimeout)

	// convert to millis
	cacheTimeout *= 1000

	backendAddress := utils.GetEnvOrDefault("POINT_ADDRESS", "localhost")

	backendPort := utils.GetEnvOrDefault("POINT_PORT", "80")

	backendProtocol := utils.GetEnvOrDefault("POINT_PROTOCOL", "http")

	var mu sync.Mutex
	var cache = make(map[string]*memory.Cache)

	go memory.StartCacheCleanerService(cache, &mu)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqKey := r.URL.Path + "?" + r.URL.RawQuery
		if r.Method != "GET" {
			log.Printf("(%s) %s\n", r.Method, reqKey)

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading request body: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			req, err := http.NewRequest(r.Method, fmt.Sprintf("%s://%s:%s%s?%s", backendProtocol, backendAddress, backendPort, r.URL.Path, r.URL.RawQuery), bytes.NewReader(bodyBytes))
			if err != nil {
				log.Printf("Error creating request: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Set req headers to fetch REAL API
			for s, v := range r.Header {
				value := v[0]
				req.Header.Set(s, value)
			}
			//------------------------------------------------
			client := &http.Client{}

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error on fetch: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// set response headers from real API
			fmt.Println("Internal response headers")
			for s, v := range resp.Header {
				value := v[0]
				fmt.Println(" |_", s, value)
				w.Header().Set(s, value)
			}
			// --------------------------------
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response body: %v\n", err)
				return
			}

			w.Write(body)
			return
		}

		// Use cache instead
		mu.Lock()
		if _, exists := cache[reqKey]; exists {
			cached := cache[reqKey]
			if time.Now().UnixMilli()-int64(cached.Timestamp) < cacheTimeout {
				fmt.Printf("(%s) Using cache %s\n", r.Method, reqKey)
				for s, v := range cached.Headers {
					value := v[0]
					w.Header().Set(s, value)
				}
				w.Write(cached.Data)
				cached.Access++
				cache[reqKey] = cached
				mu.Unlock()
				return
			}
		}
		mu.Unlock()

		req, err := http.NewRequest(r.Method, fmt.Sprintf("%s://%s:%s%s?%s", backendProtocol, backendAddress, backendPort, r.URL.Path, r.URL.RawQuery), nil)
		if err != nil {
			log.Printf("Error while create request %v\n", err)
			return
		}

		// Set req headers to fetch REAL API
		for s, v := range r.Header {
			value := v[0]
			// dont set these browser cache relationed headers
			if s != "If-None-Match" && s != "ETag" {
				req.Header.Set(s, value)
			}
		}
		//------------------------------------------------

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error on fetch: %v\n", err)
			return
		}

		// set response headers from real API
		for s, v := range resp.Header {
			value := v[0]
			w.Header().Set(s, value)
		}
		// --------------------------------

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Erro ao ler", err)
			return
		}

		w.Write(body)

		// set cache
		if !arrayContains(excludePaths, req.URL.Path) {
			mu.Lock()
			cache[reqKey] = &memory.Cache{
				Data:      body,
				Headers:   resp.Header,
				Timestamp: int(time.Now().UnixMilli()),
				Access:    1,
			}
			mu.Unlock()
		}

		// verify cache length

	})

	fmt.Println("Listening 3030")
	http.ListenAndServe(":3030", nil)
}

func arrayContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
