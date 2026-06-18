package stringutils

import "fmt"

func init() {
	fmt.Println("[stringutils/reverse.go] Package stringutils initialized!")
}

// Reverse returns the reversed version of the input string.
// Exported function.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
