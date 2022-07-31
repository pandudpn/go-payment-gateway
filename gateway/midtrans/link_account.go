package midtrans

import (
	"context"
	"net/http"
	"time"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/utils"
)

// CreateLinkPayAccount is to link the customer's account for payments
func CreateLinkPayAccount(lap *LinkAccountPay, opts *pg.Options) (*LinkAccountPayResponse, error) {
	m, err := createChargeMidtrans(lap, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return m.createLinkPayAccount(ctx)
}

// CreateLinkPayAccountWithContext is to link the customer's account for payments
func CreateLinkPayAccountWithContext(ctx context.Context, lap *LinkAccountPay, opts *pg.Options) (*LinkAccountPayResponse, error) {
	m, err := createChargeMidtrans(lap, opts)
	if err != nil {
		return nil, err
	}

	return m.createLinkPayAccount(ctx)
}

func (m *midtrans) createLinkPayAccount(ctx context.Context) (*LinkAccountPayResponse, error) {
	var (
		res    LinkAccountPayResponse
		err    error
		header = make(http.Header)
	)

	// set authorization basic header
	header.Set("Authorization", utils.SetBasicAuthorization(m.opts.ServerKey, ""))

	m.uri += createPayAccountUri
	err = m.opts.ApiCall.Call(ctx, http.MethodPost, m.uri, header, m.params, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
