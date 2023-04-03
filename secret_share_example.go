package mpc

import (
	"fmt"
)

func ExampleSplit() {
	secrets := make([]byte, 256)
	for i := uint8(0); int(i) < 256; i++ {
		secrets[i] = i
	}
	fmt.Println(secrets)
}
