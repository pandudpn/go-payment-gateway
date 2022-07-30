package midtrans_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/midtrans"
	"github.com/pandudpn/go-payment-gateway/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEWalletCharge_SuccessSandBox(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamEWallet := getMockParamsEWalletBytes()
	mockUrl := getMockUrlSandBox()
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamEWallet, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	result, err := midtrans.CreateEWalletCharge(e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateEWalletCharge_SuccessProduction(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamEWallet := getMockParamsEWalletBytes()
	mockUrl := getMockUrlProduction()
	mockHeader := getMockHeaderProduction()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamEWallet, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsProduction()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	result, err := midtrans.CreateEWalletCharge(e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateEWalletCharge_ErrorParamsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	result, err := midtrans.CreateEWalletCharge(nil, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateEWalletCharge_ErrorParamsTransactionDetailIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.TransactionDetails = nil
	result, err := midtrans.CreateEWalletCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateEWalletCharge_ErrorParamsItemDetailsIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.ItemDetails = nil
	result, err := midtrans.CreateEWalletCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateEWalletCharge_ErrorParamsInvalidPaymentTypeForEWallet(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.PaymentType = midtrans.PaymentTypeBCA
	result, err := midtrans.CreateEWalletCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateEWalletCharge_ErrorParamsShopeeIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	e.PaymentType = midtrans.PaymentTypeShopeePay
	result, err := midtrans.CreateEWalletCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateEWalletCharge_ErrorCredentials(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientSecret

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	result, err := midtrans.CreateEWalletCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateEWalletCharge_ErrorRequest(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamEWallet := getMockParamsEWalletBytes()
	mockUrl := getMockUrlSandBox()
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamEWallet, &expectedResult).Return(errors.New("error request"))

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	result, err := midtrans.CreateEWalletCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateEWalletChargeWithContext_SuccessProduction(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamEWallet := getMockParamsEWalletBytes()
	mockUrl := getMockUrlProduction()
	mockHeader := getMockHeaderProduction()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", ctx, http.MethodPost, mockUrl, mockHeader, mockParamEWallet, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsProduction()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	result, err := midtrans.CreateEWalletChargeWithContext(ctx, e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateEWalletChargeWithContext_ErrorCredentials(t *testing.T) {
	ctx := context.Background()
	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientSecret

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsEWallet()
	result, err := midtrans.CreateEWalletChargeWithContext(ctx, e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}
