package xendit_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/xendit"
	"github.com/pandudpn/go-payment-gateway/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEWalletCharge_SuccessSandBox(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	e := getMockParamsEWallet()
	e.CheckoutMethod = xendit.TokenizedPayment
	e.ChannelCode = xendit.ChannelCodeOVO
	e.ChannelProperties.FailureRedirectURL = "https://www.google.com"
	e.PaymentMethodID = "abc"
	e.CustomerID = "abc"
	e.ChannelProperties.RedeemPoints = xendit.RedeemNone

	mockParamEWallet, _ := json.Marshal(e)
	mockUrl := getMockUrlEWallet()
	mockHeader := getMockHeaderSandBox()
	expectedResult := xendit.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamEWallet, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateEWalletCharge_SuccessProduction(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	e := getMockParamsEWallet()
	e.CheckoutMethod = xendit.TokenizedPayment
	e.ChannelCode = xendit.ChannelCodeShopeePay
	e.ChannelProperties.FailureRedirectURL = "https://www.google.com"
	e.PaymentMethodID = "abc"
	e.CustomerID = "abc"

	mockUrl := getMockUrlEWallet()
	mockHeader := getMockHeaderProduction()
	expectedResult := xendit.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mock.Anything, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsProduction()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateEWalletChargeWithContext_SuccessSandBox(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamEWallet := getMockParamsEWalletBytes()
	mockUrl := getMockUrlEWallet()
	mockHeader := getMockHeaderSandBox()
	expectedResult := xendit.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", ctx, http.MethodPost, mockUrl, mockHeader, mockParamEWallet, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	result, err := xendit.CreateEWalletChargeWithContext(ctx, e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateEWalletChargeWithContext_SuccessProduction(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamEWallet := getMockParamsEWalletBytes()
	mockUrl := getMockUrlEWallet()
	mockHeader := getMockHeaderProduction()
	expectedResult := xendit.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", ctx, http.MethodPost, mockUrl, mockHeader, mockParamEWallet, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsProduction()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	result, err := xendit.CreateEWalletChargeWithContext(ctx, e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateEWalletCharge_ErrorRequest(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	e := getMockParamsEWallet()
	e.CheckoutMethod = xendit.TokenizedPayment
	e.ChannelCode = xendit.ChannelCodeShopeePay
	e.ChannelProperties.FailureRedirectURL = "https://www.google.com"
	e.PaymentMethodID = "abc"
	e.CustomerID = "abc"

	mockUrl := getMockUrlEWallet()
	mockHeader := getMockHeaderProduction()
	expectedResult := xendit.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mock.Anything, &expectedResult).Return(errors.New("error request"))

	// mock options
	mockOptions := getMockOptionsProduction()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateEWalletCharge_ErrorParamsIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateEWalletCharge(nil, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletChargeWithContext_ErrorParamsIsNil(t *testing.T) {
	ctx := context.Background()
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateEWalletChargeWithContext(ctx, nil, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsReferenceID(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ReferenceID = ""
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsChannelCode(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelCode = ""
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsChannelProperties(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelProperties = nil
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_DefaultCheckoutMethod(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientSecret

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.CheckoutMethod = ""
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_DefaultCurrency(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientSecret

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.Currency = ""
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsEligibleChannelCodeEWallet(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelCode = "ABC"
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsMinAmount(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.Amount = 50
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsRequiredPaymentMethodID(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.CheckoutMethod = xendit.TokenizedPayment
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsNilSuccessRedirectURL(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelProperties = &xendit.EWalletChannelProperties{}
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_DefaultRedeemPoints(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelCode = xendit.ChannelCodeShopeePay
	e.CheckoutMethod = xendit.TokenizedPayment
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsNilMobileNumberOVO(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelCode = xendit.ChannelCodeOVO
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsInvalidPhoneNumberOVO(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelCode = xendit.ChannelCodeOVO
	e.ChannelProperties.MobileNumber = "081"
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsInvalidLenghtPhoneNumberOVO(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelCode = xendit.ChannelCodeOVO
	e.ChannelProperties.MobileNumber = "0"
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_InvalidPhilippinesPhoneNumber(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientSecret

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelCode = xendit.ChannelCodeOVO
	e.Currency = xendit.PHP
	e.ChannelProperties.MobileNumber = "083123"
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_ErrorValidationParamsNilURL(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelCode = xendit.ChannelCodeOVO
	e.CheckoutMethod = xendit.TokenizedPayment
	e.CustomerID = "abc"
	e.PaymentMethodID = "abc"
	e.ChannelProperties.FailureRedirectURL = ""
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateEWalletCharge_DefaultRedeemPointsOVO(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientSecret

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ChannelCode = xendit.ChannelCodeOVO
	e.CheckoutMethod = xendit.TokenizedPayment
	e.CustomerID = "abc"
	e.PaymentMethodID = "abc"
	e.ChannelProperties.FailureRedirectURL = "https://www.google.com"
	result, err := xendit.CreateEWalletCharge(e, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}
