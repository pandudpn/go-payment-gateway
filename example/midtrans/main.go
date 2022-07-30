package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/midtrans"
)

const sandBoxServerKey string = "SB-Mid-server-ZKoMj1ghHnJqKQy7kNyQhUOu"
const sandBoxClientKey string = "SB-Mid-client-B5YDy_W5MCk53L5U"

func ewalletCharge(opts *pg.Options) {
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
		log.Fatalln("failed to create e-wallet one_time_payment charge with error:", err)
	}

	log.Println("response e-wallet one_time_payment charge", *res)
}

func bankTransferCharge(opts *pg.Options) {
	id := uuid.NewString()
	bt := &midtrans.BankTransferCreateParams{
		PaymentType: midtrans.PaymentTypeMandiri,
		TransactionDetails: &midtrans.TransactionDetail{
			OrderID:     id,
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
		EChannel: &midtrans.EChannel{
			BillKey:   "123456789",
			BillInfo1: "Pembayaran:",
			BillInfo2: id,
		},
		// BankTransfer: &midtrans.BankTransfer{
		// 	Bank:     midtrans.BankPermata,
		// 	VANumber: "1234567890",
		// },
	}

	res, err := midtrans.CreateBankTransferCharge(bt, opts)
	if err != nil {
		log.Fatalln("failed to create bank_transfer charge with error:", err)
	}

	log.Println("response bank_transfer charge", *res)
	if len(res.VANumbers) > 0 {
		log.Println("virtual account", bt.BankTransfer.Bank, res.VANumbers[0].VANumber)
	}

	if res.PermataVANumber != "" {
		log.Println("virtual account permata", res.PermataVANumber)
	}

	if res.BillKey != "" {
		log.Println("virtual account mandiri", res.BillerCode+"-"+res.BillKey)
	}
}

func createCardToken(cr *midtrans.CardResponse, opts *pg.Options) {
	var ct = new(midtrans.CardToken)

	if cr == nil {
		ct.CardNumber = "5211111111111117"
		ct.CardExpMonth = "12"
		ct.CardExpYear = "2022"
		ct.CardCvv = "123"
	} else {
		ct.CardCvv = "123"
		ct.TokenID = cr.SavedTokenID
	}

	res, err := midtrans.CreateCardToken(ct, opts)
	if err != nil {
		log.Fatalln("failed to create card_token with error:", err)
	}

	log.Println("response card_token", *res)
}

func createRegisterCard(opts *pg.Options) *midtrans.CardResponse {
	cr := &midtrans.CardRegister{
		CardNumber:   "5211111111111117",
		CardExpMonth: "12",
		CardExpYear:  "2022",
		CardCvv:      "123",
	}

	res, err := midtrans.CreateCardRegister(cr, opts)
	if err != nil {
		log.Fatalln("failed to create register_card with error:", err)
		return nil
	}

	log.Println("response register_card", *res)
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

	// example e-wallet one_time_payment
	// ewalletCharge(opts)

	// example bank_transfer (virtual account)
	// bankTransferCharge(opts)

	// register card token for credit_card or debit_card
	cr := createRegisterCard(opts)

	// create card token for credit_card or debit_card
	createCardToken(cr, opts)
}
