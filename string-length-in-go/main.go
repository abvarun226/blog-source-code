package main

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

func main() {
	s3 := "Well done 👍🏼"
	fmt.Printf("s3 = \"%s\"\n", s3)
	fmt.Printf("len(s3)                         = %d\n", len(s3))                         // bytes count
	fmt.Printf("utf8.RuneCountInString(s3)      = %d\n", utf8.RuneCountInString(s3))      // runes count
	fmt.Printf("uniseg.GraphemeClusterCount(s3) = %d\n", uniseg.GraphemeClusterCount(s3)) // characters/graphemes count

	for _, i := range []rune(s3) {
		fmt.Printf("%v ", i)
	}

	fmt.Println()

	s1 := string(65) // integer 65 is interpreted as unicode code point and character represented by the code point 65 is "A"
	fmt.Printf("s1 =  %s, length = %d\n", s1, len(s1))

	s2 := strconv.FormatInt(65, 10) // integer 65 is converted to string "65" here.
	fmt.Printf("s2 = %s, length = %d\n", s2, len(s2))
}
