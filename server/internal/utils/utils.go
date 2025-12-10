package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// CalcSHA256 calc sha256('text'[+'extension'])
func CalcSHA256(text string, append ...string) string {
	hash := sha256.New()
	hash.Reset()
	hash.Write([]byte(text + strings.Join(append, "")))
	bytes := hash.Sum(nil)

	return hex.EncodeToString(bytes)
}
