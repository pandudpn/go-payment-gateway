package xendit

import (
	"time"

	"github.com/pandudpn/go-payment-gateway"
)

// Mapper handles conversion between unified and Xendit-specific types
type Mapper struct{}

// mapPaymentType maps unified payment type to Xendit payment method
func (m *Mapper) mapPaymentType(pt pg.PaymentType) (string, string) {
	switch pt {
	case pg.PaymentTypeGoPay:
		return "EWALLET", string(EWalletGoPay)
	case pg.PaymentTypeOVO:
		return "EWALLET", string(EWalletOVO)
	case pg.PaymentTypeDANA:
		return "EWALLET", string(EWalletDANA)
	case pg.PaymentTypeLinkAja:
		return "EWALLET", string(EWalletLinkAja)
	case pg.PaymentTypeShopeePay:
		return "EWALLET", string(EWalletShopeePay)
	case pg.PaymentTypeQRIS:
		return "QR_CODE", string(QRCodeDynamic)
	case pg.PaymentTypeVABCA:
		return "VIRTUAL_ACCOUNT", string(BankBCA)
	case pg.PaymentTypeVABNI:
		return "VIRTUAL_ACCOUNT", string(BankBNI)
	case pg.PaymentTypeVABRI:
		return "VIRTUAL_ACCOUNT", string(BankBRI)
	case pg.PaymentTypeVAMandiri:
		return "VIRTUAL_ACCOUNT", string(BankMANDIRI)
	case pg.PaymentTypeVAPermata:
		return "VIRTUAL_ACCOUNT", string(BankPERMATA)
	case pg.PaymentTypeVACIMB:
		return "VIRTUAL_ACCOUNT", string(BankCIMB)
	case pg.PaymentTypeAlfamart:
		return "RETAIL_OUTLET", string(RetailAlfamart)
	case pg.PaymentTypeIndomaret:
		return "RETAIL_OUTLET", string(RetailIndomaret)
	default:
		return "", ""
	}
}

// mapStatus maps Xendit status to unified status
func (m *Mapper) mapStatus(status PaymentStatus) pg.Status {
	switch status {
	case StatusPaid:
		return pg.StatusSuccess
	case StatusPending:
		return pg.StatusPending
	case StatusFailed:
		return pg.StatusFailed
	default:
		return pg.StatusPending
	}
}

// mapToEWalletRequest maps unified ChargeParams to Xendit EWallet request
func (m *Mapper) mapToEWalletRequest(params pg.ChargeParams) *CreateEWalletRequest {
	_, code := m.mapPaymentType(params.PaymentType)

	req := &CreateEWalletRequest{
		ExternalID:  params.OrderID,
		Amount:      float64(params.Amount),
		EWalletCode: EWalletCode(code),
		CallbackURL: params.CallbackURL,
		Currency:    "IDR",
	}

	if params.ReturnURL != "" {
		req.ChannelProperties = &ChannelProperties{
			SuccessRedirectURL: params.ReturnURL,
			FailureRedirectURL: params.ReturnURL,
			PendingRedirectURL: params.ReturnURL,
		}
	}

	if params.Customer.Phone != "" {
		req.Phone = params.Customer.Phone
	}

	return req
}

// mapToVARequest maps unified ChargeParams to Xendit VA request
func (m *Mapper) mapToVARequest(params pg.ChargeParams) *CreateVAResquest {
	_, code := m.mapPaymentType(params.PaymentType)

	req := &CreateVAResquest{
		ExternalID:     params.OrderID,
		BankCode:       BankCode(code),
		Name:           params.Customer.Name,
		ExpectedAmount: float64(params.Amount),
		IsClosed:       true,
		Currency:       "IDR",
	}

	// Set custom VA number if provided
	if vaNumber, ok := params.Custom["va_number"].(string); ok {
		req.VANumber = vaNumber
	}

	return req
}

// mapToInvoiceRequest maps unified ChargeParams to Xendit Invoice request
func (m *Mapper) mapToInvoiceRequest(params pg.ChargeParams) *CreateInvoiceRequest {
	methodType, _ := m.mapPaymentType(params.PaymentType)

	var paymentMethods []PaymentMethod
	if methodType != "" {
		paymentMethods = []PaymentMethod{
			{Type: methodType},
		}
	}

	req := &CreateInvoiceRequest{
		ExternalID:   params.OrderID,
		Amount:       float64(params.Amount),
		Currency:     "IDR",
		PaymentMethod: paymentMethods,
		Description:  params.Description,
	}

	// Set customer details
	if params.Customer.Name != "" || params.Customer.Email != "" {
		req.Customer = &CustomerDetail{
			GivenNames:   params.Customer.Name,
			Email:        params.Customer.Email,
			MobileNumber: params.Customer.Phone,
		}
	}

	// Set items
	if len(params.Items) > 0 {
		req.Items = make([]*Item, len(params.Items))
		for i, item := range params.Items {
			req.Items[i] = &Item{
				Name:     item.Name,
				Price:    float64(item.Price),
				Quantity: item.Quantity,
				Category: item.Category,
			}
		}
	}

	return req
}

// mapToChargeResponse maps Xendit Invoice response to unified ChargeResponse
func (m *Mapper) mapToChargeResponse(resp *InvoiceResponse, _ pg.PaymentType) *pg.ChargeResponse {
	if resp == nil {
		return nil
	}

	unified := &pg.ChargeResponse{
		TransactionID: resp.ID,
		OrderID:       resp.ExternalID,
		Amount:        int64(resp.Amount),
		Status:        m.mapStatus(resp.Status),
		PaymentURL:    resp.PaymentURL,
	}

	// Set created time
	if resp.Created != nil {
		unified.CreatedAt = *resp.Created
	}

	// Set VA specific fields
	if resp.PaymentDetails != nil && resp.PaymentDetails.Destination != "" {
		unified.VANumber = resp.PaymentDetails.Destination
		unified.VABank = resp.PaymentChannel
	}

	return unified
}

// mapToChargeResponseFromVA maps Xendit VA response to unified ChargeResponse
func (m *Mapper) mapToChargeResponseFromVA(resp *VAResponse, _ pg.PaymentType) *pg.ChargeResponse {
	if resp == nil {
		return nil
	}

	unified := &pg.ChargeResponse{
		TransactionID: resp.ID,
		OrderID:       resp.ExternalID,
		Amount:        int64(resp.ExpectedAmount),
		Status:        m.mapStatus(resp.Status),
		VANumber:      resp.AccountNumber,
		VABank:        string(resp.BankCode),
	}

	if resp.VANumber != "" {
		unified.VANumber = resp.VANumber
	}

	return unified
}

// mapToChargeResponseFromEWallet maps Xendit EWallet response to unified ChargeResponse
func (m *Mapper) mapToChargeResponseFromEWallet(resp *EWalletResponse, _ pg.PaymentType) *pg.ChargeResponse {
	if resp == nil {
		return nil
	}

	unified := &pg.ChargeResponse{
		TransactionID: resp.ID,
		OrderID:       resp.ExternalID,
		Amount:        int64(resp.Amount),
		Status:        m.mapStatus(resp.Status),
		PaymentURL:    resp.PaymentURL,
	}

	if resp.RedirectURL != "" {
		unified.PaymentURL = resp.RedirectURL
	}

	if resp.Created != nil {
		unified.CreatedAt = *resp.Created
	}

	return unified
}

// mapToPaymentStatus maps Xendit response to unified PaymentStatus
func (m *Mapper) mapToPaymentStatus(orderID string, resp interface{}) *pg.PaymentStatus {
	switch r := resp.(type) {
	case *InvoiceResponse:
		var paidAt *time.Time
		if r.Status == StatusPaid {
			paidAt = r.Updated
		}

		return &pg.PaymentStatus{
			TransactionID: r.ID,
			OrderID:       orderID,
			Status:        m.mapStatus(r.Status),
			Amount:        int64(r.Amount),
			PaidAmount:    int64(r.Amount),
			PaidAt:        paidAt,
		}
	case *VAResponse:
		var paidAt *time.Time
		if r.Payment != nil {
			now := time.Now()
			paidAt = &now
		}

		return &pg.PaymentStatus{
			TransactionID: r.ID,
			OrderID:       orderID,
			Status:        m.mapStatus(r.Status),
			Amount:        int64(r.ExpectedAmount),
			PaidAmount:    int64(r.ExpectedAmount),
			PaidAt:        paidAt,
		}
	}

	return nil
}

// unifiedPaymentType maps Xendit payment type to unified payment type
func (m *Mapper) unifiedPaymentType(methodType, methodCode string) pg.PaymentType {
	switch methodType {
	case "EWALLET":
		switch EWalletCode(methodCode) {
		case EWalletGoPay:
			return pg.PaymentTypeGoPay
		case EWalletOVO:
			return pg.PaymentTypeOVO
		case EWalletDANA:
			return pg.PaymentTypeDANA
		case EWalletLinkAja:
			return pg.PaymentTypeLinkAja
		case EWalletShopeePay:
			return pg.PaymentTypeShopeePay
		}
	case "VIRTUAL_ACCOUNT":
		switch BankCode(methodCode) {
		case BankBCA:
			return pg.PaymentTypeVABCA
		case BankBNI:
			return pg.PaymentTypeVABNI
		case BankBRI:
			return pg.PaymentTypeVABRI
		case BankMANDIRI:
			return pg.PaymentTypeVAMandiri
		case BankPERMATA:
			return pg.PaymentTypeVAPermata
		case BankCIMB:
			return pg.PaymentTypeVACIMB
		}
	case "QR_CODE":
		return pg.PaymentTypeQRIS
	case "RETAIL_OUTLET":
		switch RetailCode(methodCode) {
		case RetailAlfamart:
			return pg.PaymentTypeAlfamart
		case RetailIndomaret:
			return pg.PaymentTypeIndomaret
		}
	}
	return pg.PaymentType(methodType)
}
