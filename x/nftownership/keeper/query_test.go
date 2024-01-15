package keeper_test

import (
	"math/big"

	utiltx "github.com/settlus/chain/testutil/tx"
	"github.com/settlus/chain/x/nftownership/types"
)

func (suite *NftOwnershipTestSuite) TestKeeper_Params() {
	params := types.DefaultParams()
	params.AllowedChainIds = []string{"chain1", "chain2"}
	suite.app.NftOwnershipKeeper.SetParams(suite.ctx, params)

	response, err := suite.app.NftOwnershipKeeper.Params(suite.ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(params, response.Params)
}

func (suite *NftOwnershipTestSuite) TestKeeper_GetNftOwner() {
	addr, err := s.DeployContract("Bored Ape Yatch Club", "BAYC")
	suite.Require().NoError(err)
	contractAddress := addr
	ownerAddress := utiltx.GenerateAddress()

	err = s.MintNFT(contractAddress, ownerAddress)
	suite.Require().NoError(err)

	exists, err := s.CheckNFTExists(contractAddress, big.NewInt(0))
	suite.Require().NoError(err)
	suite.Require().True(exists)

	response, err := suite.app.NftOwnershipKeeper.GetNftOwner(suite.ctx, &types.QueryGetNftOwnerRequest{
		ChainId:         suite.ctx.ChainID(),
		ContractAddress: contractAddress.Hex(),
		TokenIdHex:      "0x0",
	})

	suite.Require().NoError(err)
	suite.Require().Equal(ownerAddress.Hex(), response.OwnerAddress)
}
