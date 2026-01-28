package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand/v2"
)

// HmacSHA256 calc hmac-sha256('key', 'content')
func HmacSHA256[T string | []byte](content string, key ...T) string {
	var k []byte
	if len(key) > 0 {
		k = []byte(key[0])
	}

	hasher := hmac.New(sha256.New, k) // k default nil is ok
	hasher.Write([]byte(content))
	bytes := hasher.Sum(nil)

	return hex.EncodeToString(bytes)
}

const CharactersLibrary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const useBits = 6 // 6个bit位可以表示全部字符库中的字符

// GenerateRandomBytes generate random readable Bytes
func GenerateRandomBytes[T string | []byte](length int) T {
	b := make([]byte, length)

	randomNum, remainBits := rand.Int64(), 64
	for i := 0; i < len(b); {
		if remainBits < useBits {
			randomNum, remainBits = rand.Int64(), 64
		}

		index := int(randomNum & (1<<useBits - 1)) // 0b0011 1111
		if index < len(CharactersLibrary) {
			randomNum >>= useBits
			remainBits -= useBits

			b[i] = CharactersLibrary[index]
			i++
		} else {
			randomNum >>= 1
			remainBits -= 1
		}
	}

	return T(b)
}
