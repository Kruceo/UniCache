package memory

import (
	"fmt"
	"sync"
	"time"
)

func StartCacheCleanerService(cacheMap map[string]*Cache, maxMemoryUsage int, mu *sync.Mutex) {
	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ticker.C:
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
					fmt.Println("delete", key)
					delete(cacheMap, key)
				}
			}
			// delete(cacheMap)

			mu.Unlock()
		}
	}
}
