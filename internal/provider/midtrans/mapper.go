package midtrans

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/pandudpn/go-payment-gateway"
)

// Mapper handles conversion between unified and Midtrans-specific types
type Mapper struct{}

// mapPaymentType maps unified payment type to Midtrans payment type
func (m *Mapper) mapPaymentType(pt pg.PaymentType) string {
	switch pt {
	case pg.PaymentTypeGoPay:
		return "gopay"
	case pg.PaymentTypeShopeePay:
		return "shopeepay"
	case pg.PaymentTypeOVO:
		return "ovo"
	case pg.PaymentTypeDANA:
		return "dana"
	case pg.PaymentTypeLinkAja:
		return "linkaja"
	case pg.PaymentTypeQRIS:
		return "qris"
	case pg.PaymentTypeVABCA, pg.PaymentTypeVABNI, pg.PaymentTypeVABRI, pg.PaymentTypeVAPermata, pg.PaymentTypeVACIMB:
		return "bank_transfer"
	case pg.PaymentTypeVAMandiri:
		return "echannel"
	case pg.PaymentTypeCC:
		return "credit_card"
	default:
		return string(pt)
	}
}

// mapStatus maps Midtrans status to unified status
func (m *Mapper) mapStatus(status TransactionStatus) pg.Status {
	switch status {
	case Settlement, Capture:
		return pg.StatusSuccess
	case Pending:
		return pg.StatusPending
	case Deny, Failure:
		return pg.StatusFailed
	case Cancel:
		return pg.StatusCancelled
	case Expire:
		return pg.StatusExpired
	case Authorize:
		return pg.StatusProcessing
	default:
		return pg.StatusPending
	}
}

// mapPaymentTypeToBank maps unified payment type to Midtrans bank code
func (m *Mapper) mapPaymentTypeToBank(pt pg.PaymentType) BankCode {
	switch pt {
	case pg.PaymentTypeVABCA:
		return BankBCA
	case pg.PaymentTypeVABNI:
		return BankBNI
	case pg.PaymentTypeVABRI:
		return BankBRI
	case pg.PaymentTypeVAMandiri:
		return BankMandiri
	case pg.PaymentTypeVAPermata:
		return BankPermata
	case pg.PaymentTypeVACIMB:
		return BankCIMB
	default:
		return ""
	}
}

// mapToEWalletParams maps unified ChargeParams to Midtrans EWallet params
func (m *Mapper) mapToEWalletParams(params pg.ChargeParams) *EWallet {
	e := &EWallet{
		PaymentType:       PaymentType(m.mapPaymentType(params.PaymentType)),
		TransactionDetails: &TransactionDetail{
			OrderID:     params.OrderID,
			GrossAmount: params.Amount,
		},
		ItemDetails: make([]*ItemDetail, len(params.Items)),
	}

	// Map customer details
	if len(params.Customer.Name) > 0 {
		// Split name into first and last name
		firstName := params.Customer.Name
		lastName := ""
		if len(params.Customer.Name) > 20 {
			firstName = params.Customer.Name[:20]
			lastName = params.Customer.Name[20:]
		}
		e.CustomerDetails = &CustomerDetail{
			FirstName: firstName,
			LastName:  lastName,
			Email:     params.Customer.Email,
			Phone:     params.Customer.Phone,
		}
	}

	// Map items
	for i, item := range params.Items {
		e.ItemDetails[i] = &ItemDetail{
			ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
			Category: item.Category,
		}
	}

	// Set e-wallet specific details
	var ewalletDetail *EWalletDetail
	switch params.PaymentType {
	case pg.PaymentTypeGoPay:
		ewalletDetail = &EWalletDetail{
			CallbackURL: params.CallbackURL,
		}
		// Check for custom GoPay parameters
		if accountID, ok := params.Custom["gopay_account_id"].(string); ok {
			ewalletDetail.AccountID = accountID
		}
	case pg.PaymentTypeShopeePay:
		ewalletDetail = &EWalletDetail{
			CallbackURL: params.CallbackURL,
		}
	}

	if ewalletDetail != nil {
		switch params.PaymentType {
		case pg.PaymentTypeGoPay:
			e.Gopay = ewalletDetail
		case pg.PaymentTypeShopeePay:
			e.ShopeePay = ewalletDetail
		}
	}

	return e
}

// mapToBankTransferParams maps unified ChargeParams to Midtrans BankTransfer params
func (m *Mapper) mapToBankTransferParams(params pg.ChargeParams) *BankTransferCreateParams {
	bt := &BankTransferCreateParams{
		PaymentType:       PaymentType(m.mapPaymentType(params.PaymentType)),
		TransactionDetails: &TransactionDetail{
			OrderID:     params.OrderID,
			GrossAmount: params.Amount,
		},
		ItemDetails: make([]*ItemDetail, len(params.Items)),
	}

	// Map customer details
	if len(params.Customer.Name) > 0 {
		firstName := params.Customer.Name
		lastName := ""
		if len(params.Customer.Name) > 20 {
			firstName = params.Customer.Name[:20]
			lastName = params.Customer.Name[20:]
		}
		bt.CustomerDetails = &CustomerDetail{
			FirstName: firstName,
			LastName:  lastName,
			Email:     params.Customer.Email,
			Phone:     params.Customer.Phone,
		}
	}

	// Map items
	for i, item := range params.Items {
		bt.ItemDetails[i] = &ItemDetail{
			ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
			Category: item.Category,
		}
	}

	// Set bank transfer details
	bank := m.mapPaymentTypeToBank(params.PaymentType)
	bankTransfer := &BankTransfer{
		Bank: bank,
	}

	// Get custom VA number if provided
	if vaNumber, ok := params.Custom["va_number"].(string); ok {
		bankTransfer.VANumber = vaNumber
	}

	// Handle Mandiri e-channel
	if params.PaymentType == pg.PaymentTypeVAMandiri {
		eChannel := &EChannel{
			BillInfo1: "Payment",
			BillInfo2: params.OrderID,
		}
		bt.EChannel = eChannel
	} else {
		bt.BankTransfer = bankTransfer
	}

	return bt
}

// mapToChargeResponse maps Midtrans ChargeResponse to unified ChargeResponse
func (m *Mapper) mapToChargeResponse(resp *ChargeResponse) *pg.ChargeResponse {
	if resp == nil {
		return nil
	}

	// Parse gross amount
	var amount int64
	if resp.GrossAmount != "" {
		amount, _ = strconv.ParseInt(resp.GrossAmount, 10, 64)
	}

	unified := &pg.ChargeResponse{
		TransactionID: resp.TransactionID,
		OrderID:       resp.OrderID,
		Amount:        amount,
		Status:        m.mapStatus(resp.TransactionStatus),
		PaymentURL:    resp.RedirectURL,
		CreatedAt:     resp.TransactionTime,
		VABank:        string(resp.Bank),
	}

	// Extract payment URL from actions if available
	if len(resp.Actions) > 0 {
		for _, action := range resp.Actions {
			if action.URL != "" {
				unified.PaymentURL = action.URL
				break
			}
		}
	}

	// Set VA specific fields
	if resp.PermataVANumber != "" {
		unified.VANumber = resp.PermataVANumber
		unified.VABank = "permata"
	}
	if len(resp.VANumbers) > 0 {
		unified.VANumber = resp.VANumbers[0].VANumber
		unified.VABank = string(resp.VANumbers[0].Bank)
	}
	if resp.BillKey != "" {
		unified.VANumber = resp.BillerCode + "-" + resp.BillKey
		unified.VABank = "mandiri"
	}

	// Store raw response
	raw := make(map[string]interface{})
	rawBytes, _ := json.Marshal(resp)
	json.Unmarshal(rawBytes, &raw)
	unified.Raw = raw

	return unified
}

// mapToPaymentStatus maps Midtrans response to unified PaymentStatus
func (m *Mapper) mapToPaymentStatus(orderID string, resp *ChargeResponse) *pg.PaymentStatus {
	if resp == nil {
		return nil
	}

	// Parse gross amount
	var amount int64
	if resp.GrossAmount != "" {
		amount, _ = strconv.ParseInt(resp.GrossAmount, 10, 64)
	}

	paidAt := resp.TransactionTime
	if resp.TransactionStatus != Settlement && resp.TransactionStatus != Capture {
		paidAt = time.Time{}
	}

	return &pg.PaymentStatus{
		TransactionID: resp.TransactionID,
		OrderID:       orderID,
		Status:        m.mapStatus(resp.TransactionStatus),
		Amount:        amount,
		PaidAmount:    amount,
		PaidAt:        &paidAt,
		PaymentType:   m.unifiedPaymentType(string(resp.PaymentType)),
	}
}

// unifiedPaymentType maps Midtrans payment type to unified payment type
func (m *Mapper) unifiedPaymentType(pt string) pg.PaymentType {
	switch pt {
	case "gopay":
		return pg.PaymentTypeGoPay
	case "shopeepay":
		return pg.PaymentTypeShopeePay
	case "ovo":
		return pg.PaymentTypeOVO
	case "dana":
		return pg.PaymentTypeDANA
	case "linkaja":
		return pg.PaymentTypeLinkAja
	case "qris":
		return pg.PaymentTypeQRIS
	case "credit_card":
		return pg.PaymentTypeCC
	default:
		return pg.PaymentType(pt)
	}
}

// mapEventType maps Midtrans status to unified event type
func (m *Mapper) mapEventType(status TransactionStatus) string {
	switch status {
	case Settlement, Capture:
		return pg.EventPaymentCompleted
	case Deny, Failure:
		return pg.EventPaymentFailed
	case Pending:
		return pg.EventPaymentPending
	case Expire:
		return pg.EventPaymentExpired
	case Cancel:
		return pg.EventPaymentCancelled
	default:
		return pg.EventPaymentPending
	}
}
