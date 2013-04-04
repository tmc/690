// package trie
package trie

// Trie implements an append-only trie for ascii strings
type Trie struct {
	children map[byte]*Trie
	count    int
}

// Create a new trie
func NewTrie() *Trie {
	return &Trie{make(map[byte]*Trie, 0), 0}
}

// Add a string to the trie, returning the resulting occurances of the string
func (t *Trie) Add(s string) int {
	if len(s) == 0 {
		t.count += 1
		return t.count
	}
	var child *Trie
	// if the child doesn't exist, create it
	if child = t.children[s[0]]; child == nil {
		child = NewTrie()
		t.children[s[0]] = child
	}
	return child.Add(s[1:])
}

// Get the number of occurances of string s in the trie
func (t *Trie) Count(s string) int {
	if len(s) == 0 {
		return t.count
	}
	child := t.children[s[0]]
	if child == nil {
		return 0
	}
	return child.Count(s[1:])
}

// Get the number of strings in the trie
func (t *Trie) Total() (total int) {
	total = t.count
	for _, child := range t.children {
		total += child.Total()
	}
	return total
}

// Get the number of strings in the trie of length n
func (t *Trie) TotalOfLength(n int) (total int) {
	if n == 0 {
		total = t.count
	}
	for _, child := range t.children {
		total += child.TotalOfLength(n - 1)
	}
	return total
}

// Get the number of strings in the trie prefixed with s
func (t *Trie) NumPrefixed(s string) int {
	if len(s) == 0 {
		return t.Total()
	}
	if child, ok := t.children[s[0]]; ok {
		return child.NumPrefixed(s[1:])
	}
	return 0
}

// Get the number of strings in the trie prefixed with s with length N
func (t *Trie) NumPrefixedOfLength(s string, n int) int {
	if len(s) == 0 {
		return t.TotalOfLength(n)
	}
	if child, ok := t.children[s[0]]; ok {
		return child.NumPrefixedOfLength(s[1:], n-1)
	}
	return 0
}
