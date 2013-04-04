// package trie
package trie

// Trie implements an append-only trie for ascii strings
type Trie struct {
	children map[byte]*Trie
	count    int
}

// Create a new trie
func NewTrie() *Trie {
	at := &Trie{}
	at.children = make(map[byte]*Trie, 0)
	return at
}

// Add a string to the trie, returning the resulting occurances of the string
func (at *Trie) Add(s string) int {
	if len(s) == 0 {
		at.count += 1
		return at.count
	}
	var child *Trie
	// if the child doesn't exist, create it
	if child = at.children[s[0]]; child == nil {
		child = NewTrie()
		at.children[s[0]] = child
	}
	return child.Add(s[1:])
}

// Get the number of occurances of string s in the trie
func (at *Trie) Count(s string) int {
	if len(s) == 0 {
		return at.count
	}
	child := at.children[s[0]]
	if child == nil {
		return 0
	}
	return child.Count(s[1:])
}

// Get the number of strings in the trie
func (at *Trie) Total() (total int) {
	total = at.count
	for _, child := range at.children {
		total += child.Total()
	}
	return total
}

// Get the number of strings in the trie of length n
func (at *Trie) TotalOfLength(n int) (total int) {
	if n == 0 {
		total = at.count
	}
	for _, child := range at.children {
		total += child.TotalOfLength(n - 1)
	}
	return total
}

// Get the number of strings in the trie prefixed with s
func (at *Trie) NumPrefixed(s string) int {
	if len(s) == 0 {
		return at.Total()
	}
	if child, ok := at.children[s[0]]; ok {
		return child.NumPrefixed(s[1:])
	}
	return 0
}

// Get the number of strings in the trie prefixed with s with length N
func (at *Trie) NumPrefixedOfLength(s string, n int) int {
	if len(s) == 0 {
		return at.TotalOfLength(n)
	}
	if child, ok := at.children[s[0]]; ok {
		return child.NumPrefixedOfLength(s[1:], n-1)
	}
	return 0
}
