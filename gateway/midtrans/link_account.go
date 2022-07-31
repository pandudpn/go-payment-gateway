package midtrans

import (
	"context"
	"fmt"
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

// GetLinkPayAccountStatus is to link the customer's account for payments
func GetLinkPayAccountStatus(accountId string, opts *pg.Options) (*LinkAccountPayResponse, error) {
	m, err := createChargeMidtrans(accountId, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return m.getLinkPayAccountStatus(ctx)
}

// GetLinkPayAccountStatusWithContext is to link the customer's account for payments
func GetLinkPayAccountStatusWithContext(ctx context.Context, accountId string, opts *pg.Options) (*LinkAccountPayResponse, error) {
	m, err := createChargeMidtrans(accountId, opts)
	if err != nil {
		return nil, err
	}

	return m.getLinkPayAccountStatus(ctx)
}

func (m *midtrans) getLinkPayAccountStatus(ctx context.Context) (*LinkAccountPayResponse, error) {
	var (
		res    LinkAccountPayResponse
		err    error
		header = make(http.Header)
	)

	// set authorization basic header
	header.Set("Authorization", utils.SetBasicAuthorization(m.opts.ServerKey, ""))

	m.uri += fmt.Sprintf("%s/%s", createPayAccountUri, string(m.params))
	err = m.opts.ApiCall.Call(ctx, http.MethodGet, m.uri, header, nil, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
