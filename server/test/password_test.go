package test

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand/v2"
	"testing"
)

func TestPasswordWithSalt(t *testing.T) {
	salt := string(generateRandomBytes(10))
	password := "123456"
	password = CalcSHA256(password)
	password = CalcSHA256(password + salt)

	t.Log("pwd :", password)
	t.Log("salt:", salt)
}

func CalcSHA256(text string) string {
	hash := sha256.New()
	hash.Reset()
	hash.Write([]byte(text))
	bytes := hash.Sum(nil)

	return hex.EncodeToString(bytes)
}

const CharactersLibrary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const useBits = 6 // 6个bit位可以表示全部字符库中的字符

// GenerateRandomBytes generate random 'length' readable Bytes
func generateRandomBytes(length int) []byte {
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

	return b
}
