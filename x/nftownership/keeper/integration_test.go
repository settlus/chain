package keeper_test

import (
	"math/big"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"

	utiltx "github.com/settlus/chain/testutil/tx"
)

var _ = Describe("Nft license module integration tests", Ordered, func() {
	BeforeEach(func() {
		s.SetupTest()
	})

	var contractAddress common.Address
	ownerAddress := utiltx.GenerateAddress()

	Context("with deployed NFT contracts", func() {
		BeforeEach(func() {
			// Deploy NFT contracts
			addr, err := s.DeployContract("Bored Ape Yatch Club", "BAYC")
			Expect(err).To(BeNil())
			contractAddress = addr

			err = s.MintNFT(contractAddress, ownerAddress)
			Expect(err).To(BeNil())
		})
		It("can get correct owner of NFT", func() {
			exists, err := s.CheckNFTExists(contractAddress, big.NewInt(0))
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())

			owner, err := s.app.NftOwnershipKeeper.FindInternalOwner(
				s.ctx,
				contractAddress.Hex(),
				"0x0",
			)
			Expect(err).To(BeNil())
			Expect(owner.Hex()).To(Equal(ownerAddress.Hex()))
		})
		It("fails to get owner of non-existent NFT", func() {
			randomContractAddress := utiltx.GenerateAddress()
			exists, err := s.CheckNFTExists(randomContractAddress, big.NewInt(0))
			Expect(err).To(BeNil())
			Expect(exists).To(BeFalse())

			owner, err := s.app.NftOwnershipKeeper.FindInternalOwner(
				s.ctx,
				randomContractAddress.Hex(),
				"0x0",
			)
			Expect(err).ToNot(BeNil())
			Expect(owner).To(BeNil())
		})
	})
})
