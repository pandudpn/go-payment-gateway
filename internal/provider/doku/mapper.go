package doku

import (
	"time"

	"github.com/pandudpn/go-payment-gateway"
)

// Mapper handles conversion between unified and Doku-specific types
type Mapper struct{}

// mapPaymentType maps unified payment type to Doku payment type
func (m *Mapper) mapPaymentType(pt pg.PaymentType) PaymentType {
	switch pt {
	case pg.PaymentTypeGoPay, pg.PaymentTypeOVO, pg.PaymentTypeDANA, pg.PaymentTypeLinkAja, pg.PaymentTypeShopeePay:
		return PaymentTypeEWallet
	case pg.PaymentTypeQRIS:
		return PaymentTypeQRCode
	case pg.PaymentTypeVABCA, pg.PaymentTypeVABNI, pg.PaymentTypeVABRI, pg.PaymentTypeVAMandiri, pg.PaymentTypeVAPermata, pg.PaymentTypeVACIMB:
		return PaymentTypeVirtualAccount
	default:
		return PaymentTypeVirtualAccount
	}
}

// mapStatus maps Doku status to unified status
func (m *Mapper) mapStatus(status PaymentStatus) pg.Status {
	switch status {
	case StatusSuccess:
		return pg.StatusSuccess
	case StatusPending:
		return pg.StatusPending
	case StatusFailed:
		return pg.StatusFailed
	case StatusCancelled:
		return pg.StatusCancelled
	default:
		return pg.StatusPending
	}
}

// mapToGenerateRequest maps unified ChargeParams to Doku GeneratePaymentRequest
func (m *Mapper) mapToGenerateRequest(params pg.ChargeParams) *GeneratePaymentRequest {
	paymentType := m.mapPaymentType(params.PaymentType)

	req := &GeneratePaymentRequest{
		OrderAmount:   params.Amount,
		TransactionID: params.OrderID,
		PaymentType:   paymentType,
		TransactionDate: time.Now(),
		Locale:        "en",
	}

	// Set customer
	if params.Customer.Name != "" || params.Customer.Email != "" {
		req.Customer = &Customer{
			Name:  params.Customer.Name,
			Email: params.Customer.Email,
			Phone: params.Customer.Phone,
			ID:    params.Customer.ID,
		}
	}

	// Set payment detail based on type
	req.PaymentDetail = &PaymentDetail{}

	if paymentType == PaymentTypeEWallet {
		req.PaymentDetail.EWallet = &EWalletComponent{
			Name:        string(params.PaymentType),
			EWalletType: string(params.PaymentType),
			Amount:      formatAmount(params.Amount),
		}
		if params.Customer.Phone != "" {
			req.PaymentDetail.EWallet.Phone = params.Customer.Phone
		}
	} else if paymentType == PaymentTypeQRCode {
		req.PaymentDetail.QRCode = &QRCodeComponent{
			Name:   "QRIS",
			Amount: formatAmount(params.Amount),
			QRType: "DYNAMIC",
		}
	} else if paymentType == PaymentTypeVirtualAccount {
		req.PaymentDetail.VirtualAccount = &VAComponent{
			Name:   "VIRTUAL_ACCOUNT",
			VaType: string(params.PaymentType),
			Amount: formatAmount(params.Amount),
		}
	}

	return req
}

// mapToChargeResponse maps Doku GeneratePaymentResponse to unified ChargeResponse
func (m *Mapper) mapToChargeResponse(resp *GeneratePaymentResponse, paymentType pg.PaymentType) *pg.ChargeResponse {
	if resp == nil {
		return nil
	}

	unified := &pg.ChargeResponse{
		TransactionID: resp.TransactionID,
		OrderID:       resp.TransactionID,
		Amount:        resp.OrderAmount,
		Status:        pg.StatusPending, // Doku returns pending on successful creation
		PaymentURL:    resp.PaymentURL,
		VANumber:      resp.VANumber,
		VABank:        resp.VABank,
		CreatedAt:     time.Now(),
	}

	return unified
}

// mapToPaymentStatus maps Doku TransactionStatusResponse to unified PaymentStatus
func (m *Mapper) mapToPaymentStatus(orderID string, resp *TransactionStatusResponse) *pg.PaymentStatus {
	if resp == nil {
		return nil
	}

	var paidAt *time.Time
	if resp.PaymentDate != nil && resp.TransactionStatus == StatusSuccess {
		paidAt = resp.PaymentDate
	}

	return &pg.PaymentStatus{
		TransactionID: resp.TransactionID,
		OrderID:       orderID,
		Status:        m.mapStatus(resp.TransactionStatus),
		Amount:        resp.OrderAmount,
		PaidAmount:    resp.OrderAmount,
		PaidAt:        paidAt,
	}
}

// mapEventType maps Doku status to event type
func (m *Mapper) mapEventType(status PaymentStatus) string {
	switch status {
	case StatusSuccess:
		return pg.EventPaymentCompleted
	case StatusFailed:
		return pg.EventPaymentFailed
	case StatusCancelled:
		return pg.EventPaymentCancelled
	case StatusPending:
		return pg.EventPaymentPending
	default:
		return pg.EventPaymentPending
	}
}

// formatAmount formats amount as string for Doku
func formatAmount(amount int64) string {
	return string(rune(amount))
}
