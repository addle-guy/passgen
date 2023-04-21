package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
)

const defaultLength = 16

func main() {
	// Charsets
	const lowercaseCharset = "abcdefghijklmnopqrstuvwxyz"
	const uppercaseCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numCharset = "0123456789"
	const symbolCharset = "!#$%&@?*^~"

	// String config
	var length int
	var exclLowercase bool
	var exclUppercase bool
	var exclNums bool
	var exclSymbols bool

	// Parse flags
	flag.IntVar(&length, "len", defaultLength, "length of random string")
	flag.BoolVar(&exclNums, "no-num", false, "exclude nums")
	flag.BoolVar(&exclLowercase, "no-lower", false, "exclude lowercase")
	flag.BoolVar(&exclUppercase, "no-upper", false, "exclude uppercase")
	flag.BoolVar(&exclSymbols, "no-sym", false, "exclude special symbols")
	flag.Parse()

	// Get charset
	var builder strings.Builder
	if !exclNums {
		builder.WriteString(numCharset)
	}
	if !exclLowercase {
		builder.WriteString(lowercaseCharset)
	}
	if !exclUppercase {
		builder.WriteString(uppercaseCharset)
	}
	if !exclSymbols {
		builder.WriteString(symbolCharset)
	}
	charset := builder.String()
	if len(charset) == 0 {
		log.Fatal("Incorrect set of arguments")
	}

	// Generate random runes from given charset
	builder.Reset()
	outputChan := make(chan rune, length)
	wg := &sync.WaitGroup{}
	for length > 0 {
		wg.Add(1)
		go func(charset string, output chan rune, wg *sync.WaitGroup) {
			defer wg.Done()
			dice := rand.Intn(len([]rune(charset)))
			for n, r := range charset {
				if n == dice {
					output <- r
				}
			}
		}(charset, outputChan, wg)
		length -= 1
	}

	// Wait all done
	wg.Wait()
	close(outputChan)

	// Collect output runes
	for r := range outputChan {
		builder.WriteRune(r)
	}

	// Output string
	fmt.Println(builder.String())
}
