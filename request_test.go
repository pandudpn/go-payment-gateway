package pg_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	pg "github.com/pandudpn/go-payment-gateway"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApiRequest_Call_SuccessWithDataStruct(t *testing.T) {
	ctx := context.Background()

	mockOptions := getMockOptionsTrue()
	mockOptions, _ = pg.NewOption(mockOptions)

	mockData := getMockEWalletMidtrans()
	mockUrl := getMockUrl()

	var header = make(http.Header)
	var dest interface{}

	err := mockOptions.ApiCall.Call(ctx, http.MethodPost, mockUrl, header, mockData, &dest)
	assert.Nil(t, err, "error should be nil")
}

func TestApiRequest_Call_SuccessWithDataBytes(t *testing.T) {
	ctx := context.Background()

	mockOptions := getMockOptionsTrue()
	mockOptions, _ = pg.NewOption(mockOptions)

	mockData := getMockEWalletMidtransBytes()
	mockUrl := getMockUrl()

	var header = make(http.Header)
	var dest interface{}

	err := mockOptions.ApiCall.Call(ctx, http.MethodPost, mockUrl, header, mockData, &dest)
	assert.Nil(t, err, "error should be nil")
}

func TestApiRequest_Call_ErrorMarshal(t *testing.T) {
	ctx := context.Background()

	mockOptions := getMockOptionsTrue()
	mockOptions, _ = pg.NewOption(mockOptions)

	mockData := make(chan error, 1)
	mockUrl := getMockUrl()

	var header = make(http.Header)
	var dest interface{}

	err := mockOptions.ApiCall.Call(ctx, http.MethodPost, mockUrl, header, mockData, &dest)
	assert.NotNil(t, err, "error should not to be nil")
}

func TestApiRequest_Call_ErrorDoRequest(t *testing.T) {
	ctx := context.Background()

	mockOptions := getMockOptionsTrue()
	mockOptions, _ = pg.NewOption(mockOptions)

	mockData := getMockEWalletMidtrans()
	mockUrl := mock.Anything

	var header = make(http.Header)
	var dest interface{}

	err := mockOptions.ApiCall.Call(ctx, http.MethodPost, mockUrl, header, mockData, &dest)
	assert.NotNil(t, err, "error should not to be nil")
}

func TestApiRequest_Call_ErrorNewRequest(t *testing.T) {
	mockOptions := getMockOptionsTrue()
	mockOptions, _ = pg.NewOption(mockOptions)

	mockData := getMockEWalletMidtrans()
	mockUrl := getMockUrl()

	var header = make(http.Header)
	var dest interface{}

	err := mockOptions.ApiCall.Call(nil, http.MethodPost, mockUrl, header, mockData, &dest)
	assert.NotNil(t, err, "error should not to be nil")
}

func TestApiRequest_DoRequest_ErrorUnmarshalBodyResponse(t *testing.T) {
	mockOptions := getMockOptionsTrue()
	mockOptions, _ = pg.NewOption(mockOptions)

	mockData := getMockEWalletMidtransBytes()
	mockUrl := "https://www.google.com/api"

	req, _ := http.NewRequest(http.MethodPost, mockUrl, bytes.NewBuffer(mockData))

	var dest interface{}

	err := mockOptions.ApiCall.DoRequest(req, &dest)
	assert.NotNil(t, err, "error should not to be nil")
}
