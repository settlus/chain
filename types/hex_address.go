package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

const AccAddressByteSize = 20

type HexAddressString string

func NewHexAddressString(addr sdk.AccAddress) HexAddressString {
	return HexAddressString(common.BytesToAddress(addr.Bytes()).Hex())
}

func (a HexAddressString) Marshal() ([]byte, error) {
	return common.FromHex(string(a)), nil
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
