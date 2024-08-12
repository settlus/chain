package types

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const AccAddressByteSize = 20

type HexAddressString string

func NewHexAddrFromBytes(addr []byte) HexAddressString {
	return HexAddressString(common.BytesToAddress(addr).Hex())
}

func NormalizeHexAddress(addr string) HexAddressString {
	return HexAddressString(common.HexToAddress(addr).Hex())
}

func (a HexAddressString) Marshal() ([]byte, error) {
	return a.Bytes(), nil
}

func (a HexAddressString) MarshalTo(data []byte) (int, error) {
	bz, err := a.Marshal()
	if err != nil {
		return 0, err
	}
	copy(data, bz)
	return len(bz), nil
}

func (a *HexAddressString) Unmarshal(data []byte) error {
	*a = HexAddressString(common.BytesToAddress(data).Hex())

	return nil
}

func (a HexAddressString) Size() int {
	return AccAddressByteSize
}

func (a HexAddressString) String() string {
	return string(a)
}

func (a HexAddressString) IsNull() bool {
	for _, b := range a.Bytes() {
		if b != 0 {
			return false
		}
	}

	return true
}

func (a HexAddressString) Bytes() []byte {
	aStr := strings.TrimPrefix(string(a), "0x")
	if len(aStr)%2 == 1 {
		aStr = "0" + aStr
	}

	return common.Hex2BytesFixed(aStr, AccAddressByteSize)
}
