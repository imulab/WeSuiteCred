package stringx

import "math/rand"

const alphaNumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandAlphaNumeric generates a random string of length n from the charset of [a-zA-Z0-9].
func RandAlphaNumeric(n int) string {
	return RandStringWithCharset(n, alphaNumeric)
}

// RandStringWithCharset generates a random string of length n from the given charset.
func RandStringWithCharset(n int, charset string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
