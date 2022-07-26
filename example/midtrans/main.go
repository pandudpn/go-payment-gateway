package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/internal/midtrans"
)

func main() {
	var err error
	payment := pg.New()
	err = payment.SetMidtransCredentials("SB-Mid-server-ZKoMj1ghHnJqKQy7kNyQhUOu")
	if err != nil {
		panic(err)
	}

	err = payment.SetXenditCredentials("xnd_development_spnIw1pkbX4akj5wfZ9DBAmYIRFlikR1fmiROpTk51IWgek4JwfW8YcKSJ1rUMU")
	if err != nil {
		panic(err)
	}

	e := &midtrans.EWallet{
		PaymentType: midtrans.EWalletShopeePay,
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

	mds, err := payment.Midtrans()
	if err != nil {
		panic(err)
	}

	charge, err := mds.CreateEWalletCharge(e)
	if err != nil {
		panic(err)
	}

	fmt.Println("midtrans", *charge)
}
