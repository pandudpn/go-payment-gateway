package pg

import "testing"

func TestPaymentType_IsEWallet(t *testing.T) {
	tests := []struct {
		name     string
		payment  PaymentType
		expected bool
	}{
		{"GoPay", PaymentTypeGoPay, true},
		{"OVO", PaymentTypeOVO, true},
		{"DANA", PaymentTypeDANA, true},
		{"ShopeePay", PaymentTypeShopeePay, true},
		{"LinkAja", PaymentTypeLinkAja, true},
		{"QRIS", PaymentTypeQRIS, false},
		{"VA BCA", PaymentTypeVABCA, false},
		{"Credit Card", PaymentTypeCC, false},
		{"Alfamart", PaymentTypeAlfamart, false},
		{"Unknown", PaymentType("UNKNOWN"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.payment.IsEWallet(); got != tt.expected {
				t.Errorf("IsEWallet() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPaymentType_IsVirtualAccount(t *testing.T) {
	tests := []struct {
		name     string
		payment  PaymentType
		expected bool
	}{
		{"VA BCA", PaymentTypeVABCA, true},
		{"VA BNI", PaymentTypeVABNI, true},
		{"VA BRI", PaymentTypeVABRI, true},
		{"VA Mandiri", PaymentTypeVAMandiri, true},
		{"VA Permata", PaymentTypeVAPermata, true},
		{"VA CIMB", PaymentTypeVACIMB, true},
		{"GoPay", PaymentTypeGoPay, false},
		{"QRIS", PaymentTypeQRIS, false},
		{"Credit Card", PaymentTypeCC, false},
		{"Alfamart", PaymentTypeAlfamart, false},
		{"Unknown", PaymentType("UNKNOWN"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.payment.IsVirtualAccount(); got != tt.expected {
				t.Errorf("IsVirtualAccount() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPaymentType_IsQRIS(t *testing.T) {
	tests := []struct {
		name     string
		payment  PaymentType
		expected bool
	}{
		{"QRIS", PaymentTypeQRIS, true},
		{"GoPay", PaymentTypeGoPay, false},
		{"VA BCA", PaymentTypeVABCA, false},
		{"Credit Card", PaymentTypeCC, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.payment.IsQRIS(); got != tt.expected {
				t.Errorf("IsQRIS() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPaymentType_IsCreditCard(t *testing.T) {
	tests := []struct {
		name     string
		payment  PaymentType
		expected bool
	}{
		{"Credit Card", PaymentTypeCC, true},
		{"GoPay", PaymentTypeGoPay, false},
		{"QRIS", PaymentTypeQRIS, false},
		{"VA BCA", PaymentTypeVABCA, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.payment.IsCreditCard(); got != tt.expected {
				t.Errorf("IsCreditCard() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPaymentType_IsRetail(t *testing.T) {
	tests := []struct {
		name     string
		payment  PaymentType
		expected bool
	}{
		{"Alfamart", PaymentTypeAlfamart, true},
		{"Indomaret", PaymentTypeIndomaret, true},
		{"GoPay", PaymentTypeGoPay, false},
		{"VA BCA", PaymentTypeVABCA, false},
		{"QRIS", PaymentTypeQRIS, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.payment.IsRetail(); got != tt.expected {
				t.Errorf("IsRetail() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPaymentType_String(t *testing.T) {
	tests := []struct {
		name     string
		payment  PaymentType
		expected string
	}{
		{"GoPay", PaymentTypeGoPay, "GOPAY"},
		{"OVO", PaymentTypeOVO, "OVO"},
		{"VA BCA", PaymentTypeVABCA, "VA_BCA"},
		{"QRIS", PaymentTypeQRIS, "QRIS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.payment.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{"Pending", StatusPending, "PENDING"},
		{"Processing", StatusProcessing, "PROCESSING"},
		{"Success", StatusSuccess, "SUCCESS"},
		{"Failed", StatusFailed, "FAILED"},
		{"Cancelled", StatusCancelled, "CANCELLED"},
		{"Expired", StatusExpired, "EXPIRED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatus_IsFinal(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"Pending", StatusPending, false},
		{"Processing", StatusProcessing, false},
		{"Success", StatusSuccess, true},
		{"Failed", StatusFailed, true},
		{"Cancelled", StatusCancelled, true},
		{"Expired", StatusExpired, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsFinal(); got != tt.expected {
				t.Errorf("IsFinal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEnvironmentType_String(t *testing.T) {
	tests := []struct {
		name     string
		env      EnvironmentType
		expected string
	}{
		{"Production", Production, "production"},
		{"Sandbox", SandBox, "sandbox"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.env.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEnvironmentType_IsProduction(t *testing.T) {
	tests := []struct {
		name     string
		env      EnvironmentType
		expected bool
	}{
		{"Production", Production, true},
		{"Sandbox", SandBox, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.env.IsProduction(); got != tt.expected {
				t.Errorf("IsProduction() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBooleanVariables(t *testing.T) {
	// Test that True and False are addressable
	truePtr := &True
	falsePtr := &False

	if *truePtr != true {
		t.Error("True variable should be true")
	}

	if *falsePtr != false {
		t.Error("False variable should be false")
	}

	if truePtr == falsePtr {
		t.Error("True and False should have different addresses")
	}
}
