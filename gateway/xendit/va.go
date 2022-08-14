package xendit

import (
	"context"
	"net/http"
	"time"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/utils"
)

// CreateVirtualAccount will create a new Virtual Account Number for your customer
func CreateVirtualAccount(cva *CreateVirtualAccountParam, opts *pg.Options) (*VirtualAccount, error) {
	x, err := createChargeXendit(cva, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return x.createVirtualAccount(ctx)
}

// CreateVirtualAccountWithContext will create a new Virtual Account Number for your customer with context params
func CreateVirtualAccountWithContext(ctx context.Context, cva *CreateVirtualAccountParam, opts *pg.Options) (*VirtualAccount, error) {
	x, err := createChargeXendit(cva, opts)
	if err != nil {
		return nil, err
	}

	return x.createVirtualAccount(ctx)
}

func (x *xendit) createVirtualAccount(ctx context.Context) (*VirtualAccount, error) {
	var (
		va     VirtualAccount
		err    error
		header = make(http.Header)
	)
	// set uri target
	x.uri = uri + vaUri

	// set basic auth
	header.Set("Authorization", utils.SetBasicAuthorization(x.opts.ServerKey, ""))

	err = x.opts.ApiCall.Call(ctx, http.MethodPost, x.uri, header, x.params, &va)
	if err != nil {
		return nil, err
	}

	return &va, nil
}
