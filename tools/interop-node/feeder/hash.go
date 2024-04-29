package feeder

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/settlus/chain/tools/interop-node/types"
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
func GeneratePrevoteHash(vda types.VoteDataArr, salt string) string {
	var sb strings.Builder
	sb.WriteString(salt)
	for _, vd := range vda {
		for _, d := range vd.Data {
			sb.WriteString(d)
		}
	}
	h := sha256.New()
	h.Write([]byte(sb.String()))
	bs := h.Sum(nil)
	return fmt.Sprintf("%X", bs)
}
