package pg

import (
	"context"
	"time"

	mds "github.com/pandudpn/go-payment-gateway/internal/midtrans"
	"github.com/pandudpn/go-payment-gateway/utils"
)

const (
	mdUriSandbox    string = "https://api.sandbox.midtrans.com"
	mdUriProduction string = "https://api.midtrans.com"
)

// Midtrans configuration
type Midtrans struct {
	// uri is base url of Midtrans Core API
	uri string

	// credentials key for access Midtrans Core API
	credentials *Credentials
}

// CreateEWalletCharge charge a payment e-wallet to payment gateway midtrans
func (p *PG) CreateEWalletCharge(e *mds.EWallet) (*mds.ChargeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return p.createEWalletCharge(ctx, e)
}

// CreateEWalletChargeWithContext charge a payment e-wallet with context
func (p *PG) CreateEWalletChargeWithContext(ctx context.Context, e *mds.EWallet) (*mds.ChargeResponse, error) {
	return p.createEWalletCharge(ctx, e)
}

// createEWalletCharge do a request to midtrans to charge payment e-wallet
func (p *PG) createEWalletCharge(ctx context.Context, e *mds.EWallet) (*mds.ChargeResponse, error) {
	// check general parameters required
	// if not exists, just given error parameters invalid
	if e == nil || e.TransactionDetails == nil || (e.ItemDetails == nil || len(e.ItemDetails) < 1) {
		utils.Log.Error("one or parameters midtrans.EWallet is nil")
		return nil, ErrInvalidParameter
	}

	// if credentials midtrans not exists, send error credentials nil
	if p.midtrans == nil {
		utils.Log.Errorf("credentials midtrans is nil")
		return nil, ErrNilCredentials
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
	req := e.CreateRequest()
	req.SetURI(p.midtrans.uri + "/v2/charge")
	req.SetUsername(p.midtrans.credentials.ClientSecret)

	charge, err := req.Do(ctx)
	return charge, err
}
