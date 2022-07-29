package midtrans

import (
	"context"
	"net/http"
	"time"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/utils"
)

// CreateBankTransferCharge will create a new instance of Charge Payment Bank Transfer (VA)
func CreateBankTransferCharge(bt *BankTransferCreateParams, opts *pg.Options) (*ChargeResponse, error) {
	m, err := createChargeMidtrans(bt, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return m.createBankTransferCharge(ctx)
}

// CreateBankTransferChargeWithContext will create a new instance of Charge Payment Bank Transfer (VA) with context
func CreateBankTransferChargeWithContext(ctx context.Context, bt *BankTransferCreateParams, opts *pg.Options) (*ChargeResponse, error) {
	m, err := createChargeMidtrans(bt, opts)
	if err != nil {
		return nil, err
	}

	return m.createBankTransferCharge(ctx)
}

func (m *midtrans) createBankTransferCharge(ctx context.Context) (*ChargeResponse, error) {
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
