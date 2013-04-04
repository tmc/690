// package cracker
package cracker

import (
	"fmt"
	"github.com/traviscline/690/vigenere"
	"runtime"
	"sync"
)

func (c *cracker) CrackVigenereBruteForce(ciphertext string, keyLength, firstWordLength int) (chan Result, error) {
	results := make(chan Result)
	if c.dictionary == nil {
		return nil, fmt.Errorf("Dictionary not set.")
	}

	firstWordCiphered := ciphertext[:firstWordLength]

	N := runtime.NumCPU()
	N = 26

	var wg sync.WaitGroup
	wg.Add(N)
	done := make(chan struct{})
	for i := 'A'; i <= 'Z'; i++ {
		go c.checkKeys(byte(i), keyLength, firstWordCiphered, ciphertext, &wg, results, done)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results, nil
}

func (c *cracker) checkKeys(start byte, keyLength int, firstWord, ciphertext string, wg *sync.WaitGroup, results chan Result, done chan struct{}) {
	defer wg.Done()
	firstWordBytes := []byte(vigenere.CleanString(firstWord))
	dest := make([]byte, len(firstWord))
	keys := keyGenerator(start, keyLength-1, 'A', 'Z')
	for {
		select {
		case key, ok := <-keys:
			if !ok {
				return
			}
			vigenere.DecryptBytes(firstWordBytes, []byte(key), dest)
			if c.dictionary.Count(string(dest)) > 0 {
				p := vigenere.Decrypt(ciphertext, string(key))
				results <- Result{
					Key:       string(key),
					Plaintext: p,
				}
				close(done)
			}
		case <-done:
			return
		}
	}
}

func keyGenerator(start byte, length int, min, max byte) chan string {
	out := make(chan string, 1000)
	go func() {
		defer close(out)
		current := make([]byte, length+1)
		current[0] = start
		running := true
		for running {
			for i := len(current) - 1; i >= 1; i-- {
				if current[i] < min {
					current[i] = min
				} else if current[i] == max {
					current[i] = min
					if i == 0 {
						running = false
					}
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
