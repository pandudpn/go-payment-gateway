package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/xendit"
)

const sandBoxServerKey string = "xnd_development_spnIw1pkbX4akj5wfZ9DBAmYIRFlikR1fmiROpTk51IWgek4JwfW8YcKSJ1rUMU"
const sandBoxClientKey string = "xnd_public_development_NLo6NyaSZ5dj8RmD6zJEmFkYnxMa8ZkcGnWHj1FZbuRPzPQsC5N4F5caHxDYC2z"

func main() {
	var err error

	opts := &pg.Options{
		ServerKey: sandBoxServerKey,
		ClientId:  sandBoxClientKey,
	}

	opts, err = pg.NewOption(opts)
	if err != nil {
		log.Fatalln("create payment gateway options failed with error:", err)
	}

	e := &xendit.EWallet{
		ChannelCode: xendit.ChannelCodeShopeePay,
		Amount:      10000,
		ReferenceID: uuid.NewString(),
		ChannelProperties: &xendit.EWalletChannelProperties{
			SuccessRedirectURL: "https://www.google.com",
		},
	}

	res, err := xendit.CreateEWalletCharge(e, opts)
	if err != nil {
		log.Fatalln("failed to create e-wallet charge with error:", err)
	}

	log.Println("response", *res)
}
