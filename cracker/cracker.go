// package cracker
package cracker

import (
	"bufio"
	"fmt"
	"github.com/traviscline/690/trie"
	"github.com/traviscline/690/vigenere"
	"io"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type cracker struct {
	dictionary                 *trie.Trie
	cipherText                 string
	keyLength, firstWordLength int
}

type Result struct {
	Key, Plaintext string
}

func NewCracker() *cracker {
	return &cracker{}
}

func (c *cracker) SetDictionary(r io.Reader) error {
	c.dictionary = trie.NewTrie()
	br := bufio.NewReader(r)
	var (
		s   string
		err error
	)
	for s, err = br.ReadString('\n'); err == nil; s, err = br.ReadString('\n') {
		c.dictionary.Add(strings.TrimSpace(s))
	}
	if err != io.EOF {
		return err
	} else {
		c.dictionary.Add(strings.TrimSpace(s))
	}
	return nil
}

type candidate struct {
	key         string
	numPrefixes int
	numWords    int
	plaintext   string
}

type candidateList []candidate

func (cl candidateList) Len() int {
	return len(cl)
}

func (cl candidateList) Less(i, j int) bool {
	return cl[i].numPrefixes > cl[j].numPrefixes
}

func (cl candidateList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

type byNumWords struct {
	candidateList
}

func (cl byNumWords) Less(i, j int) bool {
	return cl.candidateList[i].numWords > cl.candidateList[j].numWords
}

func (cl candidateList) nextRound(judge func(s string) int) candidateList {
	result := make(candidateList, 0, len(cl))
	sort.Sort(cl)
	for _, c := range cl {
		for char := 'A'; char <= 'Z'; char++ {
			key := c.key + string(char)
			if n := judge(key); n > 0 {
				result = append(result, candidate{key, n, 0, ""})
			}
		}
	}
	sort.Sort(result)
	return result
}

func (c *cracker) CrackVigenere(ciphertext string, keyLength, firstWordLength int) (chan Result, error) {
	results := make(chan Result)
	if c.dictionary == nil {
		return nil, fmt.Errorf("Dictionary not set.")
	}
	candidates := candidateList{candidate{"", 0, 0, ""}}

	for i := 0; i < keyLength; i++ {
		candidates = candidates.nextRound(func(s string) int {
			p := vigenere.Decrypt(ciphertext[:len(s)], s)
			n := c.dictionary.NumPrefixedOfLength(p, firstWordLength)
			return n
		})
	}
	go func() {
		for _, candidate := range candidates {
			candidate.plaintext = vigenere.Decrypt(ciphertext, candidate.key)
			if c.dictionary.Count(candidate.plaintext[:firstWordLength]) > 0 {
				results <- Result{
					Key:       candidate.key,
					Plaintext: candidate.plaintext,
				}
			}
		}
		close(results)
	}()
	return results, nil
}

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
			//fmt.Println(string(start), key)
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
