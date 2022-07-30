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

func TestCreateBankTransferCharge_SuccessSandBox(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamBankTransfer := getMockParamsBankTransferBytes()
	mockUrl := getMockUrlSandBox()
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamBankTransfer, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	result, err := midtrans.CreateBankTransferCharge(e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateBankTransferCharge_SuccessProduction(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamBankTransfer := getMockParamsBankTransferBytes()
	mockUrl := getMockUrlProduction()
	mockHeader := getMockHeaderProduction()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamBankTransfer, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsProduction()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	result, err := midtrans.CreateBankTransferCharge(e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateBankTransferCharge_ErrorParamsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	result, err := midtrans.CreateBankTransferCharge(nil, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateBankTransferCharge_ErrorParamsTransactionDetailIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	e.TransactionDetails = nil
	result, err := midtrans.CreateBankTransferCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateBankTransferCharge_ErrorParamsItemDetailsIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	e.ItemDetails = nil
	result, err := midtrans.CreateBankTransferCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateBankTransferCharge_ErrorParamsInvalidPaymentTypeForBankTransfer(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	e.PaymentType = midtrans.PaymentTypeGopay
	result, err := midtrans.CreateBankTransferCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateBankTransferCharge_ErrorParamsEChannelIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	e.PaymentType = midtrans.PaymentTypeMandiri
	result, err := midtrans.CreateBankTransferCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateBankTransferCharge_ErrorParamsBankTransferIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	e.BankTransfer = nil
	result, err := midtrans.CreateBankTransferCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateBankTransferCharge_ErrorCredentials(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientSecret

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	result, err := midtrans.CreateBankTransferCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateBankTransferCharge_ErrorRequest(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamBankTransfer := getMockParamsBankTransferBytes()
	mockUrl := getMockUrlSandBox()
	mockHeader := getMockHeaderSandBox()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mockHeader, mockParamBankTransfer, &expectedResult).Return(errors.New("error request"))

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	result, err := midtrans.CreateBankTransferCharge(e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateBankTransferChargeWithContext_SuccessProduction(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockParamBankTransfer := getMockParamsBankTransferBytes()
	mockUrl := getMockUrlProduction()
	mockHeader := getMockHeaderProduction()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", ctx, http.MethodPost, mockUrl, mockHeader, mockParamBankTransfer, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsProduction()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	result, err := midtrans.CreateBankTransferChargeWithContext(ctx, e, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateBankTransferChargeWithContext_ErrorCredentials(t *testing.T) {
	ctx := context.Background()
	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ServerKey = clientSecret

	opts, _ := pg.NewOption(mockOptions)

	e := getMockParamsBankTransfer()
	result, err := midtrans.CreateBankTransferChargeWithContext(ctx, e, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}
