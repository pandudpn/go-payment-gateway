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

func TestCreateLinkPayAccount_SuccessSandBox(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamLinkAccountPay := getMockParamsLinkAccountPayBytes()
	mockUrl := getMockUrlCreatePayAccountSandBox()
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.LinkAccountPayResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamLinkAccountPay, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	lap := getMockParamsLinkAccountPay()
	result, err := midtrans.CreateLinkPayAccount(lap, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateLinkPayAccountWithContext_SuccessSandBox(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamLinkAccountPay := getMockParamsLinkAccountPayBytes()
	mockUrl := getMockUrlCreatePayAccountSandBox()
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.LinkAccountPayResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamLinkAccountPay, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	lap := getMockParamsLinkAccountPay()
	result, err := midtrans.CreateLinkPayAccountWithContext(ctx, lap, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateLinkPayAccount_ErrorRequest(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamLinkAccountPay := getMockParamsLinkAccountPayBytes()
	mockUrl := getMockUrlCreatePayAccountSandBox()
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.LinkAccountPayResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamLinkAccountPay, &expectedResult).Return(errors.New("error request"))

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	lap := getMockParamsLinkAccountPay()
	result, err := midtrans.CreateLinkPayAccount(lap, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateLinkPayAccount_DefaultPaymentType(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamLinkAccountPay := getMockParamsLinkAccountPayBytes()
	mockUrl := getMockUrlCreatePayAccountSandBox()
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.LinkAccountPayResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamLinkAccountPay, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	lap := getMockParamsLinkAccountPay()
	lap.PaymentType = ""
	result, err := midtrans.CreateLinkPayAccount(lap, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateLinkPayAccount_ErrorValidationPaymentTypeNotMatch(t *testing.T) {
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	lap := getMockParamsLinkAccountPay()
	lap.PaymentType = midtrans.PaymentTypeBCA
	result, err := midtrans.CreateLinkPayAccount(lap, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateLinkPayAccountWithContext_ErrorValidationPaymentTypeNotMatch(t *testing.T) {
	ctx := context.Background()
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	lap := getMockParamsLinkAccountPay()
	lap.PaymentType = midtrans.PaymentTypeBCA
	result, err := midtrans.CreateLinkPayAccountWithContext(ctx, lap, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateLinkPayAccount_ErrorValidationParamsGopayPartnerIsNil(t *testing.T) {
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	lap := getMockParamsLinkAccountPay()
	lap.GopayPartner = nil
	result, err := midtrans.CreateLinkPayAccount(lap, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestGetLinkPayAccountStatus_SuccessSandBox(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockAccountId := getMockAccountId()
	mockUrl := getMockUrlCreatePayAccountSandBox() + "/" + mockAccountId
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.LinkAccountPayResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mockHeader, nil, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := midtrans.GetLinkPayAccountStatus(mockAccountId, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestGetLinkPayAccountStatusWithContext_SuccessSandBox(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockAccountId := getMockAccountId()
	mockUrl := getMockUrlCreatePayAccountSandBox() + "/" + mockAccountId
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.LinkAccountPayResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mockHeader, nil, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := midtrans.GetLinkPayAccountStatusWithContext(ctx, mockAccountId, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestGetLinkPayAccountStatus_ErrorRequest(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockAccountId := getMockAccountId()
	mockUrl := getMockUrlCreatePayAccountSandBox() + "/" + mockAccountId
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.LinkAccountPayResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mockHeader, nil, &expectedResult).Return(errors.New("error request"))

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := midtrans.GetLinkPayAccountStatus(mockAccountId, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestGetLinkPayAccountStatusWithContext_ErrorRequest(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockAccountId := getMockAccountId()
	mockUrl := getMockUrlCreatePayAccountSandBox() + "/" + mockAccountId
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.LinkAccountPayResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mockHeader, nil, &expectedResult).Return(errors.New("error request"))

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := midtrans.GetLinkPayAccountStatusWithContext(ctx, mockAccountId, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestGetLinkPayAccountStatus_ErrorCredentials(t *testing.T) {
	// mock options
	mockAccountId := getMockAccountId()
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientId

	opts, _ := pg.NewOption(mockOptions)

	result, err := midtrans.GetLinkPayAccountStatus(mockAccountId, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestGetLinkPayAccountStatusWithContext_ErrorCredentials(t *testing.T) {
	ctx := context.Background()
	// mock options
	mockAccountId := getMockAccountId()
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientId

	opts, _ := pg.NewOption(mockOptions)

	result, err := midtrans.GetLinkPayAccountStatusWithContext(ctx, mockAccountId, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}
