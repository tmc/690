// package cracker
package cracker

import (
	"bufio"
	"fmt"
	"github.com/traviscline/690/vigenere"
	"io"
	"runtime"
	"strings"
	"sync"
)

type cracker struct {
	dictionary                 map[string]bool
	cipherText                 string
	keyLength, firstWordLength int
}

func NewCracker() *cracker {
	return &cracker{}
}

func (c *cracker) SetDictionary(r io.Reader) error {
	c.dictionary = make(map[string]bool, 0)
	br := bufio.NewReader(r)
	var (
		s   string
		err error
	)
	for s, err = br.ReadString('\n'); err == nil; s, err = br.ReadString('\n') {
		c.dictionary[strings.TrimSpace(s)] = true
	}
	if err != io.EOF {
		return err
	} else {
		c.dictionary[strings.TrimSpace(s)] = true
	}
	return nil
}

func (c *cracker) CrackVigenere(ciphertext string, keyLength, firstWordLength int) (chan string, error) {
	results := make(chan string)
	if c.dictionary == nil {
		return nil, fmt.Errorf("Dictionary not set.")
	}

	firstWordCiphered := ciphertext[:firstWordLength]

	N := 1
	runtime.GOMAXPROCS(1)
	keys := keyGenerator(keyLength, 'A', 'Z')

	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go c.checkKeys(keys, firstWordCiphered, ciphertext, &wg, results)
		//go func() {
		//	for key := range keys {
		//		p := vigenere.Decrypt(firstWordCiphered, key)
		//		//if found := c.dictionary.Count(p) > 0; found {
		//		if c.dictionary[p] {
		//			results <- vigenere.Decrypt(ciphertext, key)
		//		}
		//	}
		//	wg.Done()
		//}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results, nil
}

func (c *cracker) checkKeys(keys chan string, firstWord, ciphertext string, wg *sync.WaitGroup, results chan string) {
	for key := range keys {
		p := vigenere.Decrypt(firstWord, key)
		if c.dictionary[p] {
			results <- vigenere.Decrypt(ciphertext, key)
		}
	}
	wg.Done()
}

func keyGenerator(length int, min, max byte) chan string {
	out := make(chan string)
	go func() {
		current := make([]byte, length)
		for {
			for i := len(current) - 1; i >= 0; i-- {
				if current[i] < min {
					current[i] = min
				} else if current[i] == max {
					if i == 0 {
						close(out)
						return
					}
					current[i] = min
				} else {
					current[i]++
					break
				}
			}
			out <- string(current)
		}
	}()
	return out
}
