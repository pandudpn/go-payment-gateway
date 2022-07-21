package midtrans_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	mds "github.com/pandudpn/go-payment-gateway/internal/midtrans"
	"github.com/pandudpn/go-payment-gateway/internal/midtrans/mocks"
	"github.com/stretchr/testify/mock"
)

func TestEWallet_Success(t *testing.T) {
	mockPaymentInterface := new(mocks.PaymentInterface)
	mockPaymentInterface.
		On("CreateRequest").
		On("SetUsername", mock.Anything).
		On("SetURI", mock.Anything).
		On("Do", context.Background()).Return()

	e := &mds.EWallet{
		PaymentType: mds.EWalletShopeePay,
		TransactionDetails: &mds.TransactionDetail{
			OrderID:     uuid.New().String(),
			GrossAmount: 10000,
		},
		ItemDetails: []*mds.ItemDetail{
			{
				ID:       uuid.NewString(),
				Price:    10000,
				Name:     "abc",
				Quantity: 1,
			},
		},
		ShopeePay: &mds.EWalletDetail{
			EnableCallback: false,
			CallbackURL:    "https://playground.api.pandudpn.id/psj",
		},
	}

	req := e.CreateRequest()
}
