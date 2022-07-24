package pg

import (
	"context"
	"reflect"
	"time"

	xnd "github.com/pandudpn/go-payment-gateway/internal/xendit"
	"github.com/pandudpn/go-payment-gateway/utils"
)

const xndUri string = "https://api.xendit.co"

type xendit struct {
	// uri is base url of Xendit API
	uri string

	// credentials key for access Xendit API
	credentials *Credentials
}

// CreateEWalletCharge charge a payment e-wallet to payment gateway xendit
func (x *xendit) CreateEWalletCharge(e *xnd.EWallet) (*xnd.ChargeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	return x.createEWalletCharge(ctx, e)
}

// CreateEWalletChargeWithContext charge a payment e-wallet with context
func (x *xendit) CreateEWalletChargeWithContext(ctx context.Context, e *xnd.EWallet) (*xnd.ChargeResponse, error) {
	return x.createEWalletCharge(ctx, e)
}

// createEWalletCharge do a request to xendit to charge payment e-wallet
func (x *xendit) createEWalletCharge(ctx context.Context, e *xnd.EWallet) (*xnd.ChargeResponse, error) {
	// check general parameters required
	// if not exists, just given error parameters invalid
	if e == nil || reflect.ValueOf(e.ReferenceID).IsZero() || reflect.ValueOf(e.ChannelCode).IsZero() ||
		(reflect.ValueOf(e.ChannelProperties).IsZero() || reflect.ValueOf(e.ChannelProperties).IsNil()) {
		utils.Log.Error("one or parameters in xendit.EWallet is required")
		return nil, ErrInvalidParameter
	}

	// give a default value checkoutMethod
	// default: one_time_payment
	if reflect.ValueOf(e.CheckoutMethod).IsZero() {
		e.CheckoutMethod = xnd.OneTimePayment
	}

	// give default value for currency
	// default: IDR (Indonesian Rupiah)
	if reflect.ValueOf(e.Currency).IsZero() {
		e.Currency = xnd.IDR
	}

	// check minimum amount
	if (e.Amount < 100 && e.Currency == xnd.IDR) || (e.Amount < 1 && e.Currency == xnd.PHP) {
		return nil, ErrMinAmount
	}

	// check required params for tokenized
	if e.CheckoutMethod == xnd.TokenizedPayment && (reflect.ValueOf(e.PaymentMethodID).IsZero() || reflect.ValueOf(e.CustomerID).IsZero()) {
		utils.Log.Error("parameters PaymentMethodID or CustomerID is required for CheckoutMethod Tokenized")
		return nil, ErrInvalidParameter
	}

	// check required parameter in channelProperties for specific ChannelCode
	switch e.ChannelCode {
	case xnd.EWalletDANA, xnd.EWalletLinkAja, xnd.EWalletShopeePay:
		if reflect.ValueOf(e.ChannelProperties.SuccessRedirectURL).IsZero() {
			utils.Log.Error("parameter SuccessRedirectURL is required")
			return nil, ErrInvalidParameter
		}

		// default value of redeem points when not exists
		if e.ChannelCode == xnd.EWalletShopeePay && reflect.ValueOf(e.ChannelProperties.RedeemPoints).IsZero() {
			e.ChannelProperties.RedeemPoints = xnd.RedeemNone
		}
	case xnd.EWalletOVO:
		// when checkoutMethod is OneTimePayment
		// required phoneNumber to sent notification payments to user apps
		if e.CheckoutMethod == xnd.OneTimePayment && reflect.ValueOf(e.ChannelProperties.MobileNumber).IsZero() {
			utils.Log.Error("parameter MobileNumber is required for ChannelCode OVO with CheckoutMethod one time payment")
			return nil, ErrInvalidParameter
		}

		// for checkoutMethod tokenized
		if e.CheckoutMethod == xnd.TokenizedPayment &&
			reflect.ValueOf(e.ChannelProperties.SuccessRedirectURL).IsZero() &&
			reflect.ValueOf(e.ChannelProperties.FailureRedirectURL).IsZero() {

			utils.Log.Error("one or more parameters is required")
			return nil, ErrInvalidParameter
		}

		// default value of redeem points when not exists
		if reflect.ValueOf(e.ChannelProperties.RedeemPoints).IsZero() {
			e.ChannelProperties.RedeemPoints = xnd.RedeemNone
		}
	}

	// create a instance e-wallet request
	ewallet := xnd.NewChargeEWallet(e)
	ewallet.SetURI(x.uri + "/ewallets/charges")
	ewallet.SetUsername(x.credentials.ClientSecret)

	charge, err := ewallet.Do(ctx)
	return charge, err
}
