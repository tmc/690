// package vigenere implements the Vigenere cipher.
package vigenere

func vigenereBytes(text, key, destination []byte, encrypt bool) {
	for i, c := range text {
		c = c - 'A'
		k := key[i%len(key)] - 'A'

		var r byte
		if encrypt {
			r = (c + k + 26) % 26
		} else {
			r = (c - k + 26) % 26 // go negative modulo is a bit funny
		}
		destination[i] = r + 'A'
	}
}

// Encrypt plaintext with given key.
// len(destination) must be greater than or equal to len(plaintext)
func EncryptBytes(plaintext, key, destination []byte) {
	vigenereBytes(plaintext, key, destination, true)
}

// Decrypt plaintext with given key.
// len(destination) must be greater than or equal to len(plaintext)
func DecryptBytes(ciphertext, key, destination []byte) {
	vigenereBytes(ciphertext, key, destination, false)
}
