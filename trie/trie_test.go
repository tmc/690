package trie_test

import (
	trie "."
	"testing"
)

func TestTrie(t *testing.T) {
	tr := trie.NewTrie()
	tr.Add("foobar")
	tr.Add("foobaz")
	tr.Add("bazbar")
	if tr.Total() != 3 {
		t.Error(tr.Total(), "!=", 3)
	}
	if tr.NumPrefixed("foo") != 2 {
		t.Error(tr.NumPrefixed("foo"), "!=", 2)
	}
	if tr.NumPrefixedOfLength("foo", 6) != 2 {
		t.Error(tr.NumPrefixedOfLength("foo", 6), "!=", 2)
	}
}
