package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

const AccAddressByteSize = 20

type HexAddressString string

func NewHexAddrFromAccAddr(addr sdk.AccAddress) HexAddressString {
	return HexAddressString(common.BytesToAddress(addr.Bytes()).Hex())
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

func (a HexAddressString) Bytes() []byte {
	aStr := strings.TrimPrefix(string(a), "0x")
	if len(aStr)%2 == 1 {
		aStr = "0" + aStr
	}

	return common.Hex2BytesFixed(aStr, AccAddressByteSize)
}
