/*
	This Go package comes with no warranty. It's a naïve implementation of RFC1319
	(http://tools.ietf.org/html/rfc1319) with corrections from the errata
	(http://www.rfc-editor.org/errata_search.php?rfc=1319), for learning purposes.
	I make no guarantee for its correctness. However, feel free to use it anyway
	you see fit. Remember that md2 is obsolete as it is considered insecure,
	so don't use this for anything other than learning purposes.
*/

package md2

import (
	"fmt"
	"testing"
)

type md2Test struct {
	sum   string
	input string
}

var testStrings = []md2Test{
	/* Standard */
	{"8350e5a3e24c153df2275c9f80692773", ""},
	{"32ec01ec4a6dac72c0ab96fb34c0b5d1", "a"},
	{"da853b0d3f88d99b30283a69e6ded6bb", "abc"},
	{"ab4f496bfb2a530b219ff33031fe06b0", "message digest"},
	{"9e2533cf38702269d573d0638d41e25f", "15-characters--"},
	{"e8796ac22354cadc1642550d87eb2e0b", "16-characters---"}, // Fails if correction from errata is not implemented
	{"4e8ddff3650292ab5a4108c3aa47940b", "abcdefghijklmnopqrstuvwxyz"},
	{"da33def2a42df13975352846c30338cd", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"},
	{"d5976f79d83d3a0dc9806c3c66f3efd8", "12345678901234567890123456789012345678901234567890123456789012345678901234567890"},
}

func TestStrings(t *testing.T) {
	fmt.Printf("\n\n")
	for i := 0; i < len(testStrings); i++ {
		test := testStrings[i]
		sum := fmt.Sprintf("%x", Sum([]byte(test.input)))
		if test.sum == sum {
			fmt.Printf("--- PASS: Sum function: md2(%s) = %s expected %s\n", test.input, sum, test.sum)
		} else {
			t.Fatalf("Sum function: md2(%s) = %s expected %s\n", test.input, sum, test.sum)
		}

	}
}

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sum([]byte(""))
	}
}
