// package cracker
package cracker

import (
	"fmt"
	"github.com/traviscline/690/vigenere"
	"sort"
)

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

func (c *cracker) CrackVigenereUsingTrie(ciphertext string, keyLength, firstWordLength int) (chan Result, error) {
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
