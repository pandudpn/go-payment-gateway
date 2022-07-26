package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/internal/xendit"
)

func main() {
	var err error
	payment := pg.New()

	err = payment.SetXenditCredentials("xnd_development_spnIw1pkbX4akj5wfZ9DBAmYIRFlikR1fmiROpTk51IWgek4JwfW8YcKSJ1rUMU")
	if err != nil {
		panic(err)
	}

	xndEWallet := &xendit.EWallet{
		ChannelCode: xendit.EWalletOVO,
		Amount:      10000,
		ReferenceID: uuid.NewString(),
		ChannelProperties: &xendit.EWalletChannelProperties{
			MobileNumber:       "6281234567890abc",
			SuccessRedirectURL: "https://www.google.com",
		},
	}

	xnd, err := payment.Xendit()
	if err != nil {
		panic(err)
	}

	charge, err := xnd.CreateEWalletCharge(xndEWallet)
	if err != nil {
		panic(err)
	}

	fmt.Println("charge", *charge)
}
