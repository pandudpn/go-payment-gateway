package midtrans

import (
	"fmt"
	"reflect"

	"github.com/pandudpn/go-payment-gateway"
)

// validationParams required for any payment method
func validationParams(params interface{}) error {
	switch params.(type) {
	case *EWallet:
		return validationEWallet(params.(*EWallet))
	case *BankTransferCreateParams:
		return validationBankTransfer(params.(*BankTransferCreateParams))
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
