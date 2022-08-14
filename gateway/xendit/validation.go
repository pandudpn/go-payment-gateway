package xendit

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	
	"github.com/pandudpn/go-payment-gateway"
)

func validationPhoneNumberCustomer(country, phone string) bool {
	if len(phone) < 2 || len(phone) > 13 {
		return false
	}
	
	// default regex is Indonesia
	var phoneRegexString = "^[62]+[0-9]"
	
	// when country is from PH (philippines)
	if country == "PHP" {
		phoneRegexString = "^[63]+[0-9]"
	}
	phoneRegex := regexp.MustCompile(phoneRegexString)
	
	return phoneRegex.MatchString(phone)
}

// ValidationParams required for any payment method
func ValidationParams(params interface{}) error {
	switch param := params.(type) {
	case *EWallet:
		return validationEWallet(param)
	case *CreateVirtualAccountParam:
		return validationCreateVirtualAccount(param)
	case *UpdateVirtualAccountParam:
		return validationUpdateVirtualAccount(param)
	default:
		return pg.ErrUnimplemented
	}
}

// validationEWallet params required for e-wallet payment
func validationEWallet(e *EWallet) error {
	var err = pg.ErrMissingParameter
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
		return fmt.Errorf("invalid channel_code. possible values are 'ChannelCodeDANA', 'ChannelCodeOVO', 'ChannelCodeLinkAja' or 'ChannelCodeShopeePay'")
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
		if e.ChannelCode == ChannelCodeShopeePay && e.CheckoutMethod == TokenizedPayment && reflect.ValueOf(e.ChannelProperties.RedeemPoints).IsZero() {
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
			if reflect.ValueOf(e.ChannelProperties.SuccessRedirectURL).IsZero() ||
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

// validationCreateVirtualAccount params required for Creating VirtualAccount Payment
func validationCreateVirtualAccount(va *CreateVirtualAccountParam) error {
	var err = pg.ErrMissingParameter
	
	// check required params ExternalID
	if reflect.ValueOf(va.ExternalID).IsZero() {
		return fmt.Errorf("%s. ExternalID is required", err)
	}
	
	// check required params BankCode
	if reflect.ValueOf(va.BankCode).IsZero() {
		return fmt.Errorf("%s. BankCode is required", err)
	}
	
	// check required params Name
	if reflect.ValueOf(va.Name).IsZero() {
		return fmt.Errorf("%s. Customer Name is required", err)
	}
	
	// validation virtual account
	if !reflect.ValueOf(va.VirtualAccountNumber).IsZero() {
		an, err := strconv.Atoi(va.VirtualAccountNumber)
		if err != nil {
			return fmt.Errorf("%s. VirtualAccountNumber must be a number", pg.ErrInvalidParameter)
		}
		
		if an < 9999000001 || an > 9999999999 {
			return fmt.Errorf("%s. VirtualAccountNumber must be in range 9999000001 - 9999999999", pg.ErrInvalidParameter)
		}
	}
	
	// check required params ExpectedAmount when IsClosed == true
	if va.IsClosed {
		if reflect.ValueOf(va.ExpectedAmount).IsZero() || va.ExpectedAmount < 1 {
			return fmt.Errorf("%s. ExpectedAmount must be greater than IDR 1 or USD 1", pg.ErrInvalidParameter)
		}
		
		switch va.BankCode {
		case BankPermata:
			if va.ExpectedAmount < 1 && va.ExpectedAmount > 9999999999 {
				return fmt.Errorf("%s. Expected amount should be around Rp 1 - Rp 9.999.999.999", pg.ErrInvalidParameter)
			}
		case BankBCA:
			if va.ExpectedAmount < 10000 && va.ExpectedAmount > 999999999999 {
				return fmt.Errorf("%s. Expected amount should be around Rp 10.000 - Rp 999.999.999.999", pg.ErrInvalidParameter)
			}
		default:
			if va.ExpectedAmount < 1 && va.ExpectedAmount > 50000000000 {
				return fmt.Errorf("%s. Expected amount should be around Rp 1 - Rp 50.000.000.000", pg.ErrInvalidParameter)
			}
		}
	}
	
	return nil
}

// checkBankCode eligible for bank code VirtualAccount
func checkBankCode(b BankCode) bool {
	var lists = []BankCode{BankBCA, BankBNI, BankBRI, BankMandiri, BankBSI, BankBJB, BankCIMB, BankDBS, BankPermata, BankSahabatSampoerna}
	
	for l := range lists {
		if lists[l] == b {
			return true
		}
	}
	
	return false
}

// validationUpdateVirtualAccount params required for Update their VirtualAccount created before
func validationUpdateVirtualAccount(va *UpdateVirtualAccountParam) error {
	var err = pg.ErrMissingParameter
	
	// check required params ExternalID
	if reflect.ValueOf(va.ID).IsZero() {
		return fmt.Errorf("%s. ID of Virtual Account is required", err)
	}
	
	return nil
}
