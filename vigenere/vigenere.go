// package vigenere implements the Vigenere cipher.
package vigenere

import (
	"strings"
	"unicode"
)

func vigenere(text, key string, encrypt bool) string {
	result := make([]byte, len(text))
	for i, c := range []byte(text) {
		c = c - 'A'
		k := key[i%len(key)] - 'A'

		var r byte
		if encrypt {
			r = (c + k + 26) % 26
		} else {
			r = (c - k + 26) % 26 // go negative modulo is a bit funny
		}
		result[i] = r + 'A'
	}
	return string(result)
}

// Encrypt plaintext with given key.
func Encrypt(plaintext, key string) string {
	return vigenere(plaintext, key, true)
}

// Decrypt plaintext with given key.
func Decrypt(ciphertext, key string) string {
	return vigenere(ciphertext, key, false)
}

// Cleans a string for use as plaintext of a key.
// Removes non-letters and converts to upper case
func CleanString(s string) string {
	s = strings.Map(func(i rune) rune {
		if !unicode.IsLetter(i) {
			return -1
		}
		return unicode.ToUpper(i)

	}, s)
	return s
}
