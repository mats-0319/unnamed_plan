package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand/v2"
	"strings"

	"github.com/google/uuid"
)

// CalcSHA256 calc sha256('text'[+'extension'])
func CalcSHA256(text string, append ...string) string {
	hash := sha256.New()
	hash.Reset()
	hash.Write([]byte(text + strings.Join(append, "")))
	bytes := hash.Sum(nil)

	return hex.EncodeToString(bytes)
}

const CharactersLibrary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const useBits = 6 // 6个bit位可以表示全部字符库中的字符

// GenerateRandomStr generate random 'length' readable Bytes
func GenerateRandomStr(length int) string {
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

	return string(b)
}

// Uuid return uuid v4 string,
// with same 'data', it will return same string,
// without 'data', it will return random string.
func Uuid[T string | []byte](data ...T) string {
	if len(data) < 1 {
		return uuid.NewString()
	}

	return uuid.NewHash(sha256.New(), uuid.Nil, []byte(data[0]), 4).String()
}
