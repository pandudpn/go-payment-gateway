package xendit_test

import (
	"context"
	"net/http"
	"testing"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/pandudpn/go-payment-gateway/gateway/xendit"
	"github.com/pandudpn/go-payment-gateway/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateVirtualAccountWithContext_1(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccountWithContext(ctx, createVA, opts)

	assert.NotNil(t, result, "result should not to be nil")
	assert.Nil(t, err, "error should be nil")
}

func TestCreateVirtualAccountWithContext_2(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccountWithContext(ctx, nil, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_1(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.NotNil(t, result, "result should not to be nil")
	assert.Nil(t, err, "error should be nil")
}

func TestCreateVirtualAccount_2(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.BankCode = xendit.BankMandiri
	createVA.ExpectedAmount = 9999999999999999

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_3(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.BankCode = xendit.BankPermata
	createVA.ExpectedAmount = 99999999999999999

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_4(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.ExpectedAmount = 5000

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_5(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.ExpectedAmount = 0

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_6(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(http.ErrHandlerTimeout)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_7(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.ExternalID = ""

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_8(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.BankCode = ""

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_9(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.Name = ""

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_10(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.BankCode = "Mandiri"

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_11(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.VirtualAccountNumber = "abc"

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestCreateVirtualAccount_12(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	createVA := getMockCreateVirtualAccount()
	createVA.VirtualAccountNumber = "123123"

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.CreateVirtualAccount(createVA, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestUpdateVirtualAccount_1(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	updateVa := getMockUpdateVirtualAccount()

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.UpdateVirtualAccount(updateVa, opts)

	assert.NotNil(t, result, "result should not to be nil")
	assert.Nil(t, err, "error should be nil")
}

func TestUpdateVirtualAccount_2(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	updateVa := getMockUpdateVirtualAccount()

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(http.ErrHandlerTimeout)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.UpdateVirtualAccount(updateVa, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestUpdateVirtualAccount_3(t *testing.T) {
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	updateVa := getMockUpdateVirtualAccount()
	updateVa.ID = ""

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.UpdateVirtualAccount(updateVa, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}

func TestUpdateVirtualAccountWithContext_1(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock request
	updateVa := getMockUpdateVirtualAccount()

	// doing mock call
	mockApiRequest.
		On("Call", mock.Anything, http.MethodPost, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.UpdateVirtualAccountWithContext(ctx, updateVa, opts)

	assert.NotNil(t, result, "result should not to be nil")
	assert.Nil(t, err, "error should be nil")
}

func TestUpdateVirtualAccountWithContext_2(t *testing.T) {
	ctx := context.Background()
	// mockApiRequest interface
	// for mocking Call API
	mockApiRequest := mocks.NewApiRequestInterface(t)

	// mock options
	mockOptions := getMockOptionsSandBox()
	mockOptions.ApiCall = mockApiRequest

	opts, _ := pg.NewOption(mockOptions)

	result, err := xendit.UpdateVirtualAccountWithContext(ctx, nil, opts)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "error should not to be nil")
}
