// package cracker
package cracker

// asciiTrie implements an append-only trie for ascii strings
type asciiTrie struct {
	children map[byte]*asciiTrie
	count    int
}

// Create a new trie
func NewAsciiTrie() *asciiTrie {
	at := &asciiTrie{}
	at.children = make(map[byte]*asciiTrie, 0)
	return at
}

// Add a string to the trie, returning the resulting occurances of the string
func (at *asciiTrie) Add(s string) int {
	if len(s) == 0 {
		at.count += 1
		return at.count
	}
	var child *asciiTrie
	// if the child doesn't exist, create it
	if child = at.children[s[0]]; child == nil {
		child = NewAsciiTrie()
		at.children[s[0]] = child
	}
	return child.Add(s[1:])
}

// Get the number of occurances of string s in the trie
func (at *asciiTrie) Count(s string) int {
	if len(s) == 0 {
		return at.count
	}
	child := at.children[s[0]]
	if child == nil {
		return 0
	}
	return child.Count(s[1:])
}
