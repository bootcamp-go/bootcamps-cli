package confirm

import (
	"fmt"
	"strings"
)

func Ask(msg string) bool {
	fmt.Printf("%s [y/n]: ", msg)

	var answer string
	fmt.Scan(&answer)

	return strings.ToLower(answer) == "y"
}
