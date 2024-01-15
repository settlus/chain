package feeder

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

const SaltLength = 4

// GenerateSalt generates a random salt string of length SaltLength by slicing a uuid string
func GenerateSalt() (string, error) {
	b := make([]byte, SaltLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", b), nil
}

// GeneratePrevoteHash generates a prevote hash from a block data string and a salt
func GeneratePrevoteHash(bd, salt string) string {
	sourceStr := fmt.Sprintf("%s%s", salt, bd)
	h := sha256.New()
	h.Write([]byte(sourceStr))
	bs := h.Sum(nil)
	return fmt.Sprintf("%X", bs)
}
