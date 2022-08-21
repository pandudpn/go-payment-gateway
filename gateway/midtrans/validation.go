package midtrans

import (
	"fmt"
	"reflect"
	"strconv"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/retgits/creditcard"
)

// ValidationParams required for any payment method
func ValidationParams(params interface{}) error {
	switch param := params.(type) {
	case *EWallet:
		return validationEWallet(param)
	case *BankTransferCreateParams:
		return validationBankTransfer(param)
	case *CardToken:
		return validationCardToken(param)
	case *CardRegister:
		return validationCardRegister(param)
	case *CardPayment:
		return validationCardPayment(param)
	case *LinkAccountPay:
		return validationLinkAccountPay(param)
	case *PaylaterCreateParams:
		return validationPaylater(param)
	default:
		return pg.ErrUnimplemented
	}
}

// validationEWallet params required for e-wallet payment
func validationEWallet(e *EWallet) error {
	var err = pg.ErrInvalidParameter
	// transactionDetails not exists
	if e.TransactionDetails == nil || reflect.ValueOf(e.TransactionDetails.OrderID).IsZero() || e.TransactionDetails.GrossAmount < 100 {
		return fmt.Errorf("%s. one of midtrans.TransactionDetail is missing or invalid", err)
	}
	// items not exists
	if e.ItemDetails == nil || len(e.ItemDetails) < 1 {
		return fmt.Errorf("%s. []midtrans.ItemDetail is required", err)
	}

	// make sure only gopay or shopeepay payment type
	if !checkPaymentTypeEWallet(e.PaymentType) {
		return fmt.Errorf("invalid payment_type. possible values are 'PaymentTypeGopay' or 'PaymentTypeShopeePay'")
	}

	switch e.PaymentType {
	case PaymentTypeShopeePay:
		if e.ShopeePay == nil {
			return fmt.Errorf("%s. midtrans.ShopeePay is required", err)
		}
	}

	return nil
}

// checkPaymentTypeEWallet eligible for payment type EWallet
func checkPaymentTypeEWallet(pt PaymentType) bool {
	var lists = []PaymentType{PaymentTypeGopay, PaymentTypeShopeePay}

	for l := range lists {
		if lists[l] == pt {
			return true
		}
	}

	return false
}

// validationBankTransfer params required for bank transfer (virtual account) payment
func validationBankTransfer(bt *BankTransferCreateParams) error {
	var err = pg.ErrInvalidParameter
	// transactionDetails not exists
	if bt.TransactionDetails == nil || reflect.ValueOf(bt.TransactionDetails.OrderID).IsZero() || bt.TransactionDetails.GrossAmount < 100 {
		return fmt.Errorf("%s. one of midtrans.TransactionDetail is missing or invalid", err)
	}
	// items not exists
	if bt.ItemDetails == nil || len(bt.ItemDetails) < 1 {
		return fmt.Errorf("%s. []midtrans.ItemDetail is required", err)
	}

	// make sure only gopay or shopeepay payment type
	if !checkPaymentTypeBankTransfer(bt.PaymentType) {
		return fmt.Errorf("invalid payment_type. possible values are 'PaymentTypeBCA', 'PaymentTypeBRI', 'PaymentTypeBNI', 'PaymentTypeMandiri', or 'PaymentTypePermata'")
	}

	switch bt.PaymentType {
	case PaymentTypeMandiri:
		if bt.EChannel == nil || (reflect.ValueOf(bt.EChannel.BillInfo1).IsZero() && reflect.ValueOf(bt.EChannel.BillInfo2).IsZero()) {
			return fmt.Errorf("%s. midtrans.EChannel is required", err)
		}
	default:
		if bt.BankTransfer == nil || reflect.ValueOf(bt.BankTransfer.Bank).IsZero() {
			return fmt.Errorf("%s. midtrans.BankTransfer is required", err)
		}
	}

	return nil
}

// checkPaymentTypeBankTransfer eligible for payment type BankTransfer
func checkPaymentTypeBankTransfer(pt PaymentType) bool {
	var lists = []PaymentType{PaymentTypeBCA, PaymentTypeBRI, PaymentTypeBNI, PaymentTypePermata, PaymentTypeMandiri}

	for l := range lists {
		if lists[l] == pt {
			return true
		}
	}

	return false
}

// validationCardToken params required for create card token
func validationCardToken(ct *CardToken) error {
	var err = pg.ErrMissingParameter

	// cvv not exists
	if reflect.ValueOf(ct.CardCvv).IsZero() {
		return fmt.Errorf("%s. CardCvv is required", err)
	}

	// skip any process if token_id is exists
	if !reflect.ValueOf(ct.TokenID).IsZero() {
		return nil
	}

	return validationCard(ct.CardNumber, ct.CardExpMonth, ct.CardExpYear, ct.CardCvv)
}

// validationCardRegister params required for register card token
func validationCardRegister(cr *CardRegister) error {
	var err = pg.ErrMissingParameter

	// cvv not exists
	if reflect.ValueOf(cr.CardCvv).IsZero() {
		return fmt.Errorf("%s. CardCvv is required", err)
	}

	return validationCard(cr.CardNumber, cr.CardExpMonth, cr.CardExpYear, cr.CardCvv)
}

// validationCard params required for CreditCard
func validationCard(cn, cem, cey, cvv string) error {
	var err = pg.ErrMissingParameter
	// card number not exists
	if reflect.ValueOf(cn).IsZero() {
		return fmt.Errorf("%s. CardNumber is required", err)
	}

	// expired month of card not exists
	if reflect.ValueOf(cem).IsZero() {
		return fmt.Errorf("%s. CardExpMonth is required", err)
	}

	// expired year of card not exists
	if reflect.ValueOf(cey).IsZero() {
		return fmt.Errorf("%s. CardExpYear is required", err)
	}

	// convert to number
	expMonth, err := strconv.Atoi(cem)
	if err != nil {
		return fmt.Errorf("CardExpMonth must be number")
	}

	expYear, err := strconv.Atoi(cey)
	if err != nil {
		return fmt.Errorf("CardExpYear must be number")
	}

	c := creditcard.Card{CVV: cvv, Number: cn, ExpiryMonth: expMonth, ExpiryYear: expYear}

	if len(c.Validate().Errors) > 0 {
		return fmt.Errorf("%s", c.Validate().Errors[0])
	}

	return nil
}

// validationCardPayment params required for create card token
func validationCardPayment(cp *CardPayment) error {
	var err = pg.ErrInvalidParameter

	// transactionDetails not exists
	if cp.TransactionDetails == nil || reflect.ValueOf(cp.TransactionDetails.OrderID).IsZero() || cp.TransactionDetails.GrossAmount < 100 {
		return fmt.Errorf("%s. one of midtrans.TransactionDetail is missing or invalid", err)
	}

	// items not exists
	if cp.ItemDetails == nil || len(cp.ItemDetails) < 1 {
		return fmt.Errorf("%s. []midtrans.ItemDetail is required", err)
	}

	// make sure only gopay or shopeepay payment type
	if !checkPaymentTypeCard(cp.PaymentType) {
		return fmt.Errorf("invalid payment_type. possible value only 'PaymentTypeCard'")
	}

	// credit_card body not exists
	if cp.CreditCard == nil {
		return fmt.Errorf("midtrans.CreditCard is missing")
	}

	// required field CreditCard.TokenID
	if reflect.ValueOf(cp.CreditCard.TokenID).IsZero() {
		return fmt.Errorf("midtrans.CreditCard.TokenID is missing")
	}

	return nil
}

// checkPaymentTypeCard eligible for payment type CreditCard
func checkPaymentTypeCard(pt PaymentType) bool {
	return pt == PaymentTypeCard
}

// validationLinkAccountPay params required for create pay account
func validationLinkAccountPay(lap *LinkAccountPay) error {
	var err = pg.ErrInvalidParameter

	// make default when paymentType from params is not exists
	if reflect.ValueOf(lap.PaymentType).IsZero() {
		lap.PaymentType = PaymentTypeGopay
	}

	// make sure only gopay or shopeepay payment type
	if !checkPaymentTypeLinkAccountPay(lap.PaymentType) {
		return fmt.Errorf("invalid payment_type. possible value only 'PaymentTypeGopay'")
	}

	// if params gopay partner is not exists
	// return an error
	if lap.GopayPartner == nil || reflect.ValueOf(lap.GopayPartner.PhoneNumber).IsZero() {
		return fmt.Errorf("%s. one or more parameters in midtrans.LinkAccountPay.GopayPartner is missing", err)
	}

	return nil
}

// checkPaymentTypeLinkAccountPay eligible for payment type PayAccount
func checkPaymentTypeLinkAccountPay(pt PaymentType) bool {
	return pt == PaymentTypeGopay
}

// validationPaylater params required for cardless credit payment (paylater)
func validationPaylater(pcp *PaylaterCreateParams) error {
	var err = pg.ErrInvalidParameter
	// transactionDetails not exists
	if pcp.TransactionDetails == nil || reflect.ValueOf(pcp.TransactionDetails.OrderID).IsZero() || pcp.TransactionDetails.GrossAmount < 100 {
		return fmt.Errorf("%s. one of midtrans.TransactionDetail is missing or invalid", err)
	}
	// items not exists
	if pcp.ItemDetails == nil || len(pcp.ItemDetails) < 1 {
		return fmt.Errorf("%s. []midtrans.ItemDetail is required", err)
	}

	// make sure only gopay or shopeepay payment type
	if !checkPaymentTypePaylater(pcp.PaymentType) {
		return fmt.Errorf("invalid payment_type. possible values are 'PaymentTypeAkulaku' or 'PaymentTypeKredivo'")
	}

	return nil
}

// checkPaymentTypePaylater eligible for payment type Cardless Credit (Paylater)
func checkPaymentTypePaylater(pt PaymentType) bool {
	var lists = []PaymentType{PaymentTypeAkulaku, PaymentTypeKredivo}

	for l := range lists {
		if lists[l] == pt {
			return true
		}
	}

	return false
}
