package midtrans_test

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	mds "github.com/pandudpn/go-payment-gateway/internal/midtrans"
)

func getMockChargeResponse() *mds.ChargeResponse {
	return &mds.ChargeResponse{
		ID:                uuid.NewString(),
		PaymentType:       "gopay",
		GrossAmount:       "10000",
		OrderID:           uuid.NewString(),
		StatusCode:        "201",
		TransactionStatus: "settlement",
		TransactionID:     uuid.NewString(),
		FraudStatus:       "accept",
		StatusMessage:     "created",
		TransactionTime:   time.Now().Format("2006-01-02 15:04:05"),
		Actions: []*mds.Action{
			{
				Name:   "deeplink",
				URL:    "https://gopay.com",
				Method: http.MethodGet,
			},
		},
	}
}
