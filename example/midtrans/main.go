package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/midtrans"
)

const sandBoxServerKey string = "SB-Mid-server-ZKoMj1ghHnJqKQy7kNyQhUOu"
const sandBoxClientKey string = "SB-Mid-client-B5YDy_W5MCk53L5U"

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

	res, err := midtrans.CreateEWalletCharge(e, opts)
	if err != nil {
		log.Fatalln("failed to create e-wallet charge with error:", err)
	}

	log.Println("response", *res)
}
