package color

import "fmt"

func Print(color string, message string) {
	var colorRegex string
	switch color {
	case "red":
		colorRegex = "\033[31m❌ "
	case "green":
		colorRegex = "\033[32m✅ "
	case "yellow":
		colorRegex = "\033[33m"
	case "blue":
		colorRegex = "\033[34m"
	case "magenta":
		colorRegex = "\033[35m"
	case "cyan":
		colorRegex = "\033[36m🚀 "
	default:
		colorRegex = "\033[37m"
	}
	fmt.Printf("\x1b[1m%s%s%s\x1b[0m\n", colorRegex, message, "\033[0m")
}
