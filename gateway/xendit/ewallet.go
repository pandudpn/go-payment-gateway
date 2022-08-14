package xendit

import (
	"context"
	"net/http"
	"time"

	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/utils"
)

// CreateEWalletCharge will create a new instance of Charge Payment EWallet
func CreateEWalletCharge(e *EWallet, opts *pg.Options) (*ChargeResponse, error) {
	x, err := createChargeXendit(e, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return x.createEWalletCharge(ctx)
}

// CreateEWalletChargeWithContext will create a new instance of Charge Payment EWallet with context
func CreateEWalletChargeWithContext(ctx context.Context, e *EWallet, opts *pg.Options) (*ChargeResponse, error) {
	x, err := createChargeXendit(e, opts)
	if err != nil {
		return nil, err
	}

	return x.createEWalletCharge(ctx)
}

func (x *xendit) createEWalletCharge(ctx context.Context) (*ChargeResponse, error) {
	var (
		chargeRes ChargeResponse
		err       error
		header    = make(http.Header)
	)
	// set uri target
	x.uri = uri + ewalletUri

	// set basic auth
	header.Set("Authorization", utils.SetBasicAuthorization(x.opts.ServerKey, ""))

	err = x.opts.ApiCall.Call(ctx, http.MethodPost, x.uri, header, x.params, &chargeRes)
	if err != nil {
		return nil, err
	}

	// check error code
	err = GetErrorCode(chargeRes)
	if err != nil {
		return nil, err
	}

	return &chargeRes, nil
}
