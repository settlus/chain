package e2e

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) TestBankTokenTransfer() {
	s.Run("send_setl_between_accounts", func() {
		var err error
		sender := treasuryAddr
		recipient := "settlus10z74aw2m660tuezej4w5zr35zye6684t5ejjmk"

		var (
			beforeSenderASetlBalance    sdk.Coin
			beforeRecipientASetlBalance sdk.Coin
		)

		s.Require().Eventually(
			func() bool {
				beforeSenderASetlBalance, err = getSpecificBalance(chainAPIEndpoint, sender, asetlDenom)
				s.Require().NoError(err)

				beforeRecipientASetlBalance, err = getSpecificBalance(chainAPIEndpoint, recipient, asetlDenom)
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
