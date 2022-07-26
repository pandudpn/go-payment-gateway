package pg

import (
	"context"
	"reflect"
	"time"

	mds "github.com/pandudpn/go-payment-gateway/internal/midtrans"
	"github.com/pandudpn/go-payment-gateway/utils"
)

const (
	mdUriSandbox    string = "https://api.sandbox.midtrans.com"
	mdUriProduction string = "https://api.midtrans.com"
)

// midtrans configuration
type midtrans struct {
	// uri is base url of midtrans Core API
	uri string

	// credentials key for access midtrans Core API
	credentials *Credentials
}

// CreateEWalletCharge charge a payment e-wallet to payment gateway midtrans
func (m *midtrans) CreateEWalletCharge(e *mds.EWallet) (*mds.ChargeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return m.createEWalletCharge(ctx, e)
}

// CreateEWalletChargeWithContext charge a payment e-wallet with context
func (m *midtrans) CreateEWalletChargeWithContext(ctx context.Context, e *mds.EWallet) (*mds.ChargeResponse, error) {
	return m.createEWalletCharge(ctx, e)
}

// createEWalletCharge do a request to midtrans to charge payment e-wallet
func (m *midtrans) createEWalletCharge(ctx context.Context, e *mds.EWallet) (*mds.ChargeResponse, error) {
	// check general parameters required
	// if not exists, just given error parameters invalid
	if e == nil || e.TransactionDetails == nil || (e.ItemDetails == nil || len(e.ItemDetails) < 1) {
		utils.Log.Error("one or parameters midtrans.EWallet is nil")
		return nil, ErrInvalidParameter
	}

	// check required each payment
	switch e.PaymentType {
	case mds.EWalletShopeePay:
		if e.ShopeePay == nil {
			utils.Log.Error("one or parameters midtrans.EWalletDetail is nil")
			return nil, ErrInvalidParameter
		}
	}

	// create a instance e-wallet request
	ewallet := mds.NewChargeEWallet(e)
	ewallet.SetURI(m.uri + "/v2/charge")
	ewallet.SetUsername(m.credentials.ClientSecret)

	charge, err := ewallet.Do(ctx)
	return charge, err
}

// CreateBankTransferCharge charge a payment bank_transfer (va) to payment gateway midtrans
func (m *midtrans) CreateBankTransferCharge(b *mds.BankTransferCreateParams) (*mds.ChargeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return m.createBankTransferCharge(ctx, b)
}

// CreateBankTransferChargeWithContext charge a payment bank_transfer (va) with context
func (m *midtrans) CreateBankTransferChargeWithContext(ctx context.Context, b *mds.BankTransferCreateParams) (*mds.ChargeResponse, error) {
	return m.createBankTransferCharge(ctx, b)
}

// createBankTransferCharge do a request to midtrans to charge payment bank_transfer (va)
func (m *midtrans) createBankTransferCharge(ctx context.Context, b *mds.BankTransferCreateParams) (*mds.ChargeResponse, error) {
	// check general parameters required
	// if not exists, just given error parameters invalid
	if b == nil || b.TransactionDetails == nil || (b.ItemDetails == nil || len(b.ItemDetails) < 1) {
		utils.Log.Error("one or parameters midtrans.BankTransferCreateParams is nil")
		return nil, ErrInvalidParameter
	}

	// check required each payment
	switch b.PaymentType {
	case mds.BankTransferMandiri:
		if b.EChannel == nil || (reflect.ValueOf(b.EChannel.BillInfo1).IsZero() && reflect.ValueOf(b.EChannel.BillInfo2).IsZero()) {
			utils.Log.Error("one or parameters midtrans.EChannel is nil")
			return nil, ErrInvalidParameter
		}
	default:
		if b.BankTransfer == nil || reflect.ValueOf(b.BankTransfer.Bank).IsZero() {
			utils.Log.Error("one or parameters midtrans.BankTransfer is nil")
			return nil, ErrInvalidParameter
		}
	}

	// create a instance e-wallet request
	ewallet := mds.NewChargeBankTransfer(b)
	ewallet.SetURI(m.uri + "/v2/charge")
	ewallet.SetUsername(m.credentials.ClientSecret)

	charge, err := ewallet.Do(ctx)
	return charge, err
}
