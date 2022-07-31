package midtrans_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/midtrans"
	"github.com/pandudpn/go-payment-gateway/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCardToken_Success(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := fmt.Sprintf("%s?%s", getMockUrlCardTokenSandBox(), getMockQueryParamsCardToken())
	expectedResult := midtrans.CardResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mock.Anything, mock.Anything, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateCardToken_SuccessWithTokenID(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := fmt.Sprintf("%s?%s", getMockUrlCardTokenSandBox(), getMockQueryParamsCardTokenBySavedTokenId())
	expectedResult := midtrans.CardResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mock.Anything, mock.Anything, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardTokenWithSavedTokenId()
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateCardTokenWithContext_Success(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := fmt.Sprintf("%s?%s", getMockUrlCardTokenSandBox(), getMockQueryParamsCardToken())
	expectedResult := midtrans.CardResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mock.Anything, mock.Anything, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	result, err := midtrans.CreateCardTokenWithContext(ctx, ct, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateCardRegister_Success(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := fmt.Sprintf("%s?%s", getMockUrlRegisterCardSandBox(), getMockQueryParamsCardToken())
	expectedResult := midtrans.CardResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mock.Anything, mock.Anything, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsRegisterCard()
	result, err := midtrans.CreateCardRegister(ct, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateCardRegisterWithContext_Success(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := fmt.Sprintf("%s?%s", getMockUrlRegisterCardSandBox(), getMockQueryParamsCardToken())
	expectedResult := midtrans.CardResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mock.Anything, mock.Anything, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsRegisterCard()
	result, err := midtrans.CreateCardRegisterWithContext(ctx, ct, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateCardToken_ErrorRequest(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := fmt.Sprintf("%s?%s", getMockUrlCardTokenSandBox(), getMockQueryParamsCardToken())
	expectedResult := midtrans.CardResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mock.Anything, mock.Anything, &expectedResult).Return(errors.New("error request"))

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardRegister_ErrorRequest(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := fmt.Sprintf("%s?%s", getMockUrlRegisterCardSandBox(), getMockQueryParamsCardToken())
	expectedResult := midtrans.CardResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodGet, mockUrl, mock.Anything, mock.Anything, &expectedResult).Return(errors.New("error request"))

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsRegisterCard()
	result, err := midtrans.CreateCardRegister(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardToken_ErrorValidationParamsCvvIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	ct.CardCvv = ""
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardTokenWithContext_ErrorValidationParamsCvvIsNil(t *testing.T) {
	ctx := context.Background()
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	ct.CardCvv = ""
	result, err := midtrans.CreateCardTokenWithContext(ctx, ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardRegister_ErrorValidationParamsCvvIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsRegisterCard()
	ct.CardCvv = ""
	result, err := midtrans.CreateCardRegister(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardRegisterWithContext_ErrorValidationParamsCvvIsNil(t *testing.T) {
	ctx := context.Background()
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsRegisterCard()
	ct.CardCvv = ""
	result, err := midtrans.CreateCardRegisterWithContext(ctx, ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardToken_ErrorValidationParamsCardNumberIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	ct.CardNumber = ""
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardToken_ErrorValidationParamsExpiredMonthIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	ct.CardExpMonth = ""
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardToken_ErrorValidationParamsExpiredYearIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	ct.CardExpYear = ""
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardToken_ErrorConvertExpiredMonth(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	ct.CardExpMonth = "abc"
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardToken_ErrorConvertExpiredYear(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	ct.CardExpYear = "abc"
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardToken_ErrorCardIsExpired(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	ct := getMockParamsCardToken()
	ct.CardExpYear = "2021"
	result, err := midtrans.CreateCardToken(ct, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardCharge_Success(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := getMockUrlSandBox()
	mockParams := getMockParamsCardChargeBytes()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mock.Anything, mockParams, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	cp := getMockParamsCardCharge()
	result, err := midtrans.CreateCardCharge(cp, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateCardChargeWithContext_Success(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := getMockUrlSandBox()
	mockParams := getMockParamsCardChargeBytes()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mock.Anything, mockParams, &expectedResult).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	cp := getMockParamsCardCharge()
	result, err := midtrans.CreateCardChargeWithContext(ctx, cp, opts)

	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, &expectedResult, result)
}

func TestCreateCardCharge_ErrorRequest(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	mockUrl := getMockUrlSandBox()
	mockParams := getMockParamsCardChargeBytes()
	expectedResult := midtrans.ChargeResponse{}

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mockUrl, mock.Anything, mockParams, &expectedResult).Return(errors.New("error request"))

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	cp := getMockParamsCardCharge()
	result, err := midtrans.CreateCardCharge(cp, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardCharge_ErrorValidationTransactionDetailsIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	cp := getMockParamsCardCharge()
	cp.TransactionDetails = nil
	result, err := midtrans.CreateCardCharge(cp, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardCharge_ErrorValidationItemDetailsIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	cp := getMockParamsCardCharge()
	cp.ItemDetails = nil
	result, err := midtrans.CreateCardCharge(cp, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardCharge_ErrorValidationPaymentTypeNotMatch(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	cp := getMockParamsCardCharge()
	cp.PaymentType = midtrans.PaymentTypeBCA
	result, err := midtrans.CreateCardCharge(cp, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardCharge_ErrorValidationParamsCreditCardIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	cp := getMockParamsCardCharge()
	cp.CreditCard = nil
	result, err := midtrans.CreateCardCharge(cp, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardCharge_ErrorValidationParamsTokenIdIsNil(t *testing.T) {
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	cp := getMockParamsCardCharge()
	cp.CreditCard.TokenID = ""
	result, err := midtrans.CreateCardCharge(cp, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestCreateCardChargeWithContext_ErrorValidationParamsTokenIdIsNil(t *testing.T) {
	ctx := context.Background()
	// mock options
	mockOptions := getMockOptionsSandBox()

	opts, _ := pg.NewOption(mockOptions)

	cp := getMockParamsCardCharge()
	cp.CreditCard.TokenID = ""
	result, err := midtrans.CreateCardChargeWithContext(ctx, cp, opts)

	assert.NotNil(t, err, "error should not to be nil")
	assert.Nil(t, result, "result should be nil")
}
