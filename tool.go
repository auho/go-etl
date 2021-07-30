package go_etl

import (
	"fmt"
)

func RemoveReplicaSliceString(slc []string) []string {
	result := make([]string, 0)
	tempMap := make(map[string]bool, len(slc))
	for _, e := range slc {
		if tempMap[e] == false {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}

func PrintColor(s string) {
	fmt.Printf("\r%c[1;40;32m %s %c[0m", 0x1B, s, 0x1B)
}
