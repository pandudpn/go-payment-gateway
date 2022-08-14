package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/xendit"
)

const sandBoxServerKey string = "xnd_development_spnIw1pkbX4akj5wfZ9DBAmYIRFlikR1fmiROpTk51IWgek4JwfW8YcKSJ1rUMU"
const sandBoxClientKey string = "xnd_public_development_NLo6NyaSZ5dj8RmD6zJEmFkYnxMa8ZkcGnWHj1FZbuRPzPQsC5N4F5caHxDYC2z"

func ewallet(opts *pg.Options) {
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

	log.Println("response e-wallet", *res)
}

func createVirtualAccount(opts *pg.Options) *xendit.VirtualAccount {
	expired := time.Now().Add(time.Duration(72) * time.Hour)
	va := &xendit.CreateVirtualAccountParam{
		ExternalID:           uuid.NewString(),
		ExpectedAmount:       10000,
		BankCode:             xendit.BankBCA,
		ExpirationDate:       expired,
		IsClosed:             true,
		IsSingleUse:          true,
		Name:                 "Pandu dwi Putra",
		VirtualAccountNumber: "9999345678",
	}

	res, err := xendit.CreateVirtualAccount(va, opts)
	if err != nil {
		log.Fatalln("failed to create virtual account with error", err)
	}

	log.Println("response create virtual-account", *res)
	return res
}

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

	// ewallet example
	// ewallet(opts)

	// create virtual account
	_ = createVirtualAccount(opts)
}
