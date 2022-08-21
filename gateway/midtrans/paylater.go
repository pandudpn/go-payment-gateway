package midtrans

import (
	"context"
	"net/http"
	"time"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/utils"
)

// CreatePaylater will create a new instance of Charge Payment Cardless Credit (Paylater)
func CreatePaylater(pcp *PaylaterCreateParams, opts *pg.Options) (*ChargeResponse, error) {
	m, err := createChargeMidtrans(pcp, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return m.createPaylater(ctx)
}

// CreatePaylater will create a new instance of Charge Payment Cardless Credit (Paylater) with Context
func CreatePaylaterWithContext(ctx context.Context, pcp *PaylaterCreateParams, opts *pg.Options) (*ChargeResponse, error) {
	m, err := createChargeMidtrans(pcp, opts)
	if err != nil {
		return nil, err
	}

	return m.createPaylater(ctx)
}

func (m *midtrans) createPaylater(ctx context.Context) (*ChargeResponse, error) {
	var (
		chargeRes ChargeResponse
		err       error
		header    = make(http.Header)
	)

	// set basic auth
	header.Set("Authorization", utils.SetBasicAuthorization(m.opts.ServerKey, ""))

	m.uri += chargeUri
	err = m.opts.ApiCall.Call(ctx, http.MethodPost, m.uri, header, m.params, &chargeRes)
	if err != nil {
		return nil, err
	}

	return &chargeRes, nil
}
