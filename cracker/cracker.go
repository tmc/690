// package cracker
package cracker

import (
	"bufio"
	"github.com/traviscline/690/trie"
	"io"
	"strings"
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
