package vigenere_test

import (
	vigenere "."
	"testing"
)

var (
	testKey string = "Lemon"
	testPT         = "Attack at Dawn"
)

func TestVigenere(t *testing.T) {

	k := vigenere.CleanString(testKey)
	p := vigenere.CleanString(testPT)

	ct := vigenere.Encrypt(p, k)
	pt := vigenere.Decrypt(ct, k)

	expected := vigenere.CleanString(p)
	if pt != expected {
		t.Errorf("'%s' != '%s'", pt, expected)
	}
}

func BenchmarkVigenereEncrypt(b *testing.B) {
	b.StopTimer()
	k := vigenere.CleanString(testKey)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = vigenere.Encrypt("Why hello there", k)
	}
}
