package memory

import "fmt"

func VerifyMemUsage(g map[string]*Cache) int {
	allLength := 0
	for _, value := range g {
		allLength += len(value.Data)
		// for _, v := range value.Headers {
		// 	allLength += len([]byte((v[0])))
		// }
	}
	fmt.Println(allLength, "bytes")
	return allLength
}
