package e2e

func (s *IntegrationTestSuite) TestLatestBLock() {
	s.Run("get latest block", func() {
		blockId, err := queryLatestBlockId(chainAPIEndpoint)

		s.Require().NoError(err)
		s.Require().NotEmpty(blockId)
	})
}
