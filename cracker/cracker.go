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

	N := 2
	runtime.GOMAXPROCS(2)
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

func bkeyGenerator(keyLength int, start, end byte) chan []byte {
	result := make(chan []byte)
	go func() {
		if keyLength == 1 {
			for char := start; char <= end; char++ {
				result <- []byte{char}
			}
			close(result)
			return
		}
		for char := start; char <= end; char++ {
			for rest := range keyGenerator(keyLength-1, start, end) {
				result <- append([]byte{char}, rest...)
			}
		}
		close(result)
	}()
	return result
}

func keyGenerator(keyLength int, start, end byte) chan string {
	result := make(chan string)
	go func() {
		defer close(result)
		if keyLength == 1 {
			for char := start; char <= end; char++ {
				result <- string(char)
			}
			return
		}
		for char := start; char <= end; char++ {
			for rest := range keyGenerator(keyLength-1, start, end) {
				result <- string(char) + rest
			}
		}
	}()
	return result
}
