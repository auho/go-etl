package tool

import (
	"fmt"
)

func PrintColor(s string) {
	fmt.Printf("\r%c[K", 0x1B)
	fmt.Printf("\r%c[1;40;32m %s %c[0m", 0x1B, s, 0x1B)
}
