package midtrans

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pandudpn/go-payment-gateway"
)

// queryParam for createToken or registerToken
const (
	clientKey    string = "client_key"
	cardNumber   string = "card_number"
	cardExpMonth string = "card_exp_month"
	cardExpYear  string = "card_exp_year"
	cardCvv      string = "card_cvv"
	tokenId      string = "token_id"
)

// CreateCardToken will create a new instance of CardToken Created
func CreateCardToken(ct *CardToken, opts *pg.Options) (*CardResponse, error) {
	m, err := createChargeMidtrans(ct, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return m.createCardToken(ctx)
}

// CreateCardTokenWithContext will create a new instance of CardToken Created with context
func CreateCardTokenWithContext(ctx context.Context, ct *CardToken, opts *pg.Options) (*CardResponse, error) {
	m, err := createChargeMidtrans(ct, opts)
	if err != nil {
		return nil, err
	}

	return m.createCardToken(ctx)
}

func (m *midtrans) createCardToken(ctx context.Context) (*CardResponse, error) {
	var (
		chargeRes CardResponse
		err       error
		header    = make(http.Header)
	)

	m.uri += fmt.Sprintf("%s?%s", createCardTokenUri, string(m.params))
	err = m.opts.ApiCall.Call(ctx, http.MethodGet, m.uri, header, nil, &chargeRes)
	if err != nil {
		return nil, err
	}

	return &chargeRes, nil
}

// CreateCardRegister will create a new instance of CardRegister Register
func CreateCardRegister(cr *CardRegister, opts *pg.Options) (*CardResponse, error) {
	m, err := createChargeMidtrans(cr, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return m.createCardRegister(ctx)
}

// CreateCardRegisterWithContext will create a new instance of CardRegister Register with context
func CreateCardRegisterWithContext(ctx context.Context, cr *CardRegister, opts *pg.Options) (*CardResponse, error) {
	m, err := createChargeMidtrans(cr, opts)
	if err != nil {
		return nil, err
	}

	return m.createCardRegister(ctx)
}

func (m *midtrans) createCardRegister(ctx context.Context) (*CardResponse, error) {
	var (
		chargeRes CardResponse
		err       error
		header    = make(http.Header)
	)

	m.uri += fmt.Sprintf("%s?%s", createRegisterCardUri, string(m.params))
	err = m.opts.ApiCall.Call(ctx, http.MethodGet, m.uri, header, nil, &chargeRes)
	if err != nil {
		return nil, err
	}

	return &chargeRes, nil
}
