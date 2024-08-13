package e2e

import (
	"crypto/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) TestBankTokenTransfer() {
	s.Run("send_setl_between_accounts", func() {
		var err error
		sender := bobAddr
		recipient := makeRandomAccAddress()

		var beforeSenderASetlBalance sdk.Coin
		beforeRecipientASetlBalance := sdk.NewCoin(asetlDenom, sdk.NewInt(0))

		s.Require().Eventually(
			func() bool {
				beforeSenderASetlBalance, err = getSpecificBalance(chainAPIEndpoint, sender, asetlDenom)
				s.Require().NoError(err)

				return true
			},
			10*time.Second,
			2*time.Second,
		)

		tokenAmount := sdk.NewCoin(asetlDenom, sdk.NewInt(1000000000))
		s.execBankSend(sender, recipient, tokenAmount.String(), standardFees.String())

		s.Require().Eventually(
			func() bool {
				afterSenderASetlBalance, err := getSpecificBalance(chainAPIEndpoint, sender, asetlDenom)
				s.Require().NoError(err)

				afterRecipientASetlBalance, err := getSpecificBalance(chainAPIEndpoint, recipient, asetlDenom)
				s.Require().NoError(err)

				decremented := beforeSenderASetlBalance.Sub(tokenAmount).IsGTE(afterSenderASetlBalance)
				incremented := beforeRecipientASetlBalance.Add(tokenAmount).IsEqual(afterRecipientASetlBalance)

				return decremented && incremented
			},
			time.Minute,
			2*time.Second,
		)
	})
}

func makeRandomAccAddress() string {
	recipientBytes := make([]byte, 20)
	_, _ = rand.Read(recipientBytes)
	return sdk.AccAddress(recipientBytes).String()
}
