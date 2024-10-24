package memory

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"unicache/utils"
)

func StartCacheCleanerService(cacheMap map[string]*Cache, mu *sync.Mutex) {
	cacheCleanerInterval, err := strconv.ParseInt(utils.GetEnvOrDefault("CACHE_CLEANER_INTERVAL", "60"), 10, 32)

	if err != nil {
		log.Printf("Not possible parse cache cleaner interval (%s): %v", utils.GetEnvOrDefault("CACHE_CLEANER_INTERVAL", "60"), err)
		return
	}
	fmt.Printf("Cache cleaner interval to %dS\n", cacheCleanerInterval)

	ticker := time.NewTicker(time.Duration(cacheCleanerInterval) * time.Second)

	for range ticker.C {
		mu.Lock()

		sum := 0.0
		count := 0.0

		for key, value := range cacheMap {
			fmt.Println(key, value.Access)
			sum += float64(value.Access)
			count++
		}
		median := sum / count

		for key, value := range cacheMap {
			if value.Access < int(median) {
				fmt.Printf("delete %s, median %00f, access %d\n", key, median, value.Access)
				delete(cacheMap, key)
			}
		}

		mu.Unlock()
	}
}
