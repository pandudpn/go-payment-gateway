package midtrans

import (
	"context"
	"net/http"
	"time"

	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/utils"
)

// CreateEWalletCharge will create a new instance of Charge Payment EWallet
func CreateEWalletCharge(e *EWallet, opts *pg.Options) (*ChargeResponse, error) {
	m, err := createChargeMidtrans(e, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return m.createEWalletCharge(ctx)
}

// CreateEWalletChargeWithContext will create a new instance of Charge Payment EWallet with context
func CreateEWalletChargeWithContext(ctx context.Context, e *EWallet, opts *pg.Options) (*ChargeResponse, error) {
	m, err := createChargeMidtrans(e, opts)
	if err != nil {
		return nil, err
	}

	return m.createEWalletCharge(ctx)
}

func (m *midtrans) createEWalletCharge(ctx context.Context) (*ChargeResponse, error) {
	var (
		chargeRes ChargeResponse
		err       error
		header    = make(http.Header)
	)

	// set basic auth
	header.Set("Authorization", utils.SetBasicAuthorization(m.opts.ServerKey, ""))

	err = m.opts.ApiCall.Call(ctx, http.MethodPost, m.uri, header, m.params, &chargeRes)
	if err != nil {
		return nil, err
	}

	return &chargeRes, nil
}
