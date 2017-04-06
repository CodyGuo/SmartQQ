package smartqq

import "fmt"

func hash33(str string) string {
	skey := []byte(str)
	e := 0
	for i, n := 0, len(str); n > i; i++ {
		e += (e << 5) + int(skey[i])
	}

	return fmt.Sprint(2147483647 & e)
}
