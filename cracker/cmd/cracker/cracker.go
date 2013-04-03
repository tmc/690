package main

import (
	"flag"
	"fmt"
	"github.com/traviscline/690/cracker"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
)

var (
	ciphertext      = flag.String("ciphertext", "", "enciphered text")
	dictionary      = flag.String("dict", "", "path to newline-separated dictionary file")
	keyLength       = flag.Int("keyLength", 0, "length of the key")
	firstWordLength = flag.Int("firstWordLength", 0, "length of the first word")
	definite        = flag.Bool("definite", false, "look for a definite match (all words matching)")
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *ciphertext == "" || *dictionary == "" || *keyLength == 0 || *firstWordLength == 0 {
		flag.Usage()
		os.Exit(1)
	}

	c := cracker.NewCracker()
	dictPath, _ := filepath.Abs(*dictionary)
	f, err := os.Open(dictPath)
	if err != nil {
		log.Fatalln(err)
	}
	c.SetDictionary(f)

	results, err := c.CrackVigenere(*ciphertext, *keyLength, *firstWordLength)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	for result := range results {
		fmt.Println("Key:  ", result.Key)
		fmt.Println("Plain:", result.Plaintext)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}
