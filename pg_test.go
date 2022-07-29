package pg_test

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/midtrans"
)

func getMockUrl() string {
	return "https://api.sandbox.midtrans.com/v2/charge"
}

func getMockOptionsFalse() *pg.Options {
	return &pg.Options{
		ServerKey:   "abc",
		ClientId:    "abc",
		Logging:     &pg.False,
		Environment: pg.SandBox,
	}
}

func getMockOptionsTrue() *pg.Options {
	return &pg.Options{
		ServerKey:   "abc",
		ClientId:    "abc",
		Logging:     &pg.True,
		Environment: pg.SandBox,
	}
}

func getMockEWalletMidtrans() *midtrans.EWallet {
	e := &midtrans.EWallet{
		PaymentType: midtrans.PaymentTypeShopeePay,
		TransactionDetails: &midtrans.TransactionDetail{
			OrderID:     uuid.New().String(),
			GrossAmount: 10000,
		},
		ItemDetails: []*midtrans.ItemDetail{
			{
				ID:       uuid.NewString(),
				Price:    10000,
				Name:     "abc",
				Quantity: 1,
			},
		},
		ShopeePay: &midtrans.EWalletDetail{
			EnableCallback: false,
			CallbackURL:    "https://playground.api.pandudpn.id/psj",
		},
	}

	return e
}

func getMockEWalletMidtransBytes() []byte {
	e := &midtrans.EWallet{
		PaymentType: midtrans.PaymentTypeShopeePay,
		TransactionDetails: &midtrans.TransactionDetail{
			OrderID:     uuid.New().String(),
			GrossAmount: 10000,
		},
		ItemDetails: []*midtrans.ItemDetail{
			{
				ID:       uuid.NewString(),
				Price:    10000,
				Name:     "abc",
				Quantity: 1,
			},
		},
		ShopeePay: &midtrans.EWalletDetail{
			EnableCallback: false,
			CallbackURL:    "https://playground.api.pandudpn.id/psj",
		},
	}

	b, _ := json.Marshal(e)

	return b
}
