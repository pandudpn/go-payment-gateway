package xendit

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/pandudpn/go-payment-gateway"
)

func validationPhoneNumberCustomer(country, phone string) bool {
	if len(phone) < 2 || len(phone) > 13 {
		return false
	}

	// default regex is Indonesia
	var phoneRegexString = "^[62]+[0-9]"

	// when country is from PH (philippines)
	if country == "PH" {
		phoneRegexString = "^[63]+[0-9]"
	}
	phoneRegex := regexp.MustCompile(phoneRegexString)

	return phoneRegex.MatchString(phone)
}

// validationParams required for any payment method
func validationParams(params interface{}) error {
	switch params.(type) {
	case *EWallet:
		return validationEWallet(params.(*EWallet))
	default:
		return pg.ErrUnimplemented
	}
}

// validationEWallet params required for e-wallet payment
func validationEWallet(e *EWallet) error {
	var err = pg.ErrInvalidParameter
	// ReferenceID not exists
	if reflect.ValueOf(e.ReferenceID).IsZero() {
		return fmt.Errorf("%s. ReferenceID is required", err)
	}

	// ChannelCode not exists
	if reflect.ValueOf(e.ChannelCode).IsZero() {
		return fmt.Errorf("%s. ChannelCode is required", err)
	}

	// ChannelProperties not exists
	if reflect.ValueOf(e.ChannelProperties).IsZero() || reflect.ValueOf(e.ChannelProperties).IsNil() {
		return fmt.Errorf("%s. ChannelProperties is required", err)
	}

	// give a default value checkoutMethod
	// default: OneTimePayment
	if reflect.ValueOf(e.CheckoutMethod).IsZero() {
		e.CheckoutMethod = OneTimePayment
	}

	// give a default value for currency
	// default: IDR (Indonesian Rupiah)
	if reflect.ValueOf(e.Currency).IsZero() {
		e.Currency = IDR
	}

	// make sure only ovo, linkaja, shopee, and dana channel code
	if !checkChannelCodeEWallet(e.ChannelCode) {
		return fmt.Errorf("invalid payment_type. possible values are 'PaymentTypeGopay' or 'PaymentTypeShopeePay'")
	}

	// check minimum amount
	if (e.Amount < 100 && e.Currency == IDR) || (e.Amount < 1 && e.Currency == PHP) {
		return pg.ErrMinAmount
	}

	// check required params for tokenized
	if e.CheckoutMethod == TokenizedPayment && (reflect.ValueOf(e.PaymentMethodID).IsZero() || reflect.ValueOf(e.CustomerID).IsZero()) {
		return fmt.Errorf("%s. PaymentMethodID or CustomerID is required for CheckoutMethod Tokenized", err)
	}

	// check required parameter in channelProperties for specific ChannelCode
	switch e.ChannelCode {
	case ChannelCodeDANA, ChannelCodeLinkAja, ChannelCodeShopeePay:
		if reflect.ValueOf(e.ChannelProperties.SuccessRedirectURL).IsZero() {
			return fmt.Errorf("%s. ChannelProperties.SuccessRedirectURL is required", err)
		}

		// default value of redeem points when not exists
		if e.ChannelCode == ChannelCodeShopeePay && reflect.ValueOf(e.ChannelProperties.RedeemPoints).IsZero() {
			e.ChannelProperties.RedeemPoints = RedeemNone
		}
	case ChannelCodeOVO:
		// when checkoutMethod is OneTimePayment
		// required phoneNumber to sent notification payments to user apps
		if e.CheckoutMethod == OneTimePayment {
			// phoneNumber is exists
			if reflect.ValueOf(e.ChannelProperties.MobileNumber).IsZero() {
				return fmt.Errorf("%s. ChannelProperties.MobileNumber is required", err)
			}

			// validation indonesian phoneNumber
			if !validationPhoneNumberCustomer(string(e.Currency), e.ChannelProperties.MobileNumber) {
				return fmt.Errorf("invalid phone number. %s", pg.ErrInvalidPhoneNumber)
			}
		}

		// for checkoutMethod tokenized
		if e.CheckoutMethod == TokenizedPayment {
			// check success redirect and failure redirect is not exists
			if reflect.ValueOf(e.ChannelProperties.SuccessRedirectURL).IsZero() &&
				reflect.ValueOf(e.ChannelProperties.FailureRedirectURL).IsZero() {

				return fmt.Errorf("%s. ChannelProperties.SuccessRedirectURL and ChannelProperties.FailureRedirectURL is required", err)
			}

			// default value of redeem points when not exists
			if reflect.ValueOf(e.ChannelProperties.RedeemPoints).IsZero() {
				e.ChannelProperties.RedeemPoints = RedeemNone
			}
		}
	}

	return nil
}

// checkChannelCodeEWallet eligible for channel code EWallet
func checkChannelCodeEWallet(c ChannelCode) bool {
	var lists = []ChannelCode{ChannelCodeDANA, ChannelCodeOVO, ChannelCodeLinkAja, ChannelCodeShopeePay}

	for l := range lists {
		if lists[l] == c {
			return true
		}
	}

	return false
}
