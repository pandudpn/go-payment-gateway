package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pg "github.com/pandudpn/go-payment-gateway"
)

func main() {
	// Initialize the Xendit client
	client, err := pg.NewClient(
		pg.WithProvider("xendit"),
		pg.WithServerKey(os.Getenv("XENDIT_SECRET_KEY")),
		pg.WithClientKey(os.Getenv("XENDIT_CALLBACK_TOKEN")),
		pg.WithEnvironment("sandbox"), // Xendit uses same URL for both environments
	)
	if err != nil {
		log.Fatalf("Failed to create Xendit client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Create an Invoice payment
	fmt.Println("=== Example 1: Invoice Payment ===")
	invoiceCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("INV-%d", 12345),
		Amount:      100000,
		PaymentType: pg.PaymentTypeQRIS, // Use QRIS for invoice
		Customer: pg.Customer{
			ID:    "CUST-001",
			Name:  "John Doe",
			Email: "john@example.com",
			Phone: "+628123456789",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-001",
				Name:     "Product A",
				Price:    50000,
				Quantity: 1,
			},
			{
				ID:       "ITEM-002",
				Name:     "Product B",
				Price:    50000,
				Quantity: 1,
			},
		},
		Description: "Test invoice payment",
		CallbackURL: "https://example.com/callback",
		ReturnURL:   "https://example.com/return",
	})
	if err != nil {
		log.Printf("Failed to create invoice charge: %v", err)
	} else {
		log.Printf("Invoice Charge created:")
		log.Printf("  Transaction ID: %s", invoiceCharge.TransactionID)
		log.Printf("  Order ID: %s", invoiceCharge.OrderID)
		log.Printf("  Amount: %d", invoiceCharge.Amount)
		log.Printf("  Status: %s", invoiceCharge.Status)
		log.Printf("  Payment URL: %s", invoiceCharge.PaymentURL)
	}

	// Example 2: Create a Virtual Account BCA payment
	fmt.Println("\n=== Example 2: Virtual Account BCA Payment ===")
	vaCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("VA-BCA-%d", 12345),
		Amount:      150000,
		PaymentType: pg.PaymentTypeVABCA,
		Customer: pg.Customer{
			ID:    "CUST-002",
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Phone: "+628123456789",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-003",
				Name:     "Product C",
				Price:    150000,
				Quantity: 1,
			},
		},
		CallbackURL: "https://example.com/callback",
	})
	if err != nil {
		log.Printf("Failed to create VA charge: %v", err)
	} else {
		log.Printf("VA Charge created:")
		log.Printf("  Transaction ID: %s", vaCharge.TransactionID)
		log.Printf("  Order ID: %s", vaCharge.OrderID)
		log.Printf("  VA Number: %s", vaCharge.VANumber)
		log.Printf("  VA Bank: %s", vaCharge.VABank)
		log.Printf("  Status: %s", vaCharge.Status)
	}

	// Example 3: Create a GoPay e-wallet payment
	fmt.Println("\n=== Example 3: GoPay E-Wallet Payment ===")
	gopayCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("EW-GOPAY-%d", 12345),
		Amount:      50000,
		PaymentType: pg.PaymentTypeGoPay,
		Customer: pg.Customer{
			ID:    "CUST-003",
			Name:  "Bob Smith",
			Email: "bob@example.com",
			Phone: "+628123456789",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-004",
				Name:     "Product D",
				Price:    50000,
				Quantity: 1,
			},
		},
		CallbackURL: "https://example.com/callback",
		ReturnURL:   "https://example.com/return",
	})
	if err != nil {
		log.Printf("Failed to create GoPay charge: %v", err)
	} else {
		log.Printf("GoPay Charge created:")
		log.Printf("  Transaction ID: %s", gopayCharge.TransactionID)
		log.Printf("  Order ID: %s", gopayCharge.OrderID)
		log.Printf("  Payment URL: %s", gopayCharge.PaymentURL)
		log.Printf("  Status: %s", gopayCharge.Status)
	}

	// Example 4: Create an OVO e-wallet payment
	fmt.Println("\n=== Example 4: OVO E-Wallet Payment ===")
	ovoCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("EW-OVO-%d", 12345),
		Amount:      75000,
		PaymentType: pg.PaymentTypeOVO,
		Customer: pg.Customer{
			ID:    "CUST-004",
			Name:  "Alice Johnson",
			Email: "alice@example.com",
			Phone: "+628123456789",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-005",
				Name:     "Product E",
				Price:    75000,
				Quantity: 1,
			},
		},
		CallbackURL: "https://example.com/callback",
		ReturnURL:   "https://example.com/return",
	})
	if err != nil {
		log.Printf("Failed to create OVO charge: %v", err)
	} else {
		log.Printf("OVO Charge created:")
		log.Printf("  Transaction ID: %s", ovoCharge.TransactionID)
		log.Printf("  Payment URL: %s", ovoCharge.PaymentURL)
		log.Printf("  Status: %s", ovoCharge.Status)
	}

	// Example 5: Create a DANA e-wallet payment
	fmt.Println("\n=== Example 5: DANA E-Wallet Payment ===")
	danaCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("EW-DANA-%d", 12345),
		Amount:      80000,
		PaymentType: pg.PaymentTypeDANA,
		Customer: pg.Customer{
			ID:    "CUST-005",
			Name:  "Charlie Brown",
			Email: "charlie@example.com",
			Phone: "+628123456789",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-006",
				Name:     "Product F",
				Price:    80000,
				Quantity: 1,
			},
		},
		CallbackURL: "https://example.com/callback",
	})
	if err != nil {
		log.Printf("Failed to create DANA charge: %v", err)
	} else {
		log.Printf("DANA Charge created:")
		log.Printf("  Transaction ID: %s", danaCharge.TransactionID)
		log.Printf("  Payment URL: %s", danaCharge.PaymentURL)
		log.Printf("  Status: %s", danaCharge.Status)
	}

	// Example 6: Create a ShopeePay e-wallet payment
	fmt.Println("\n=== Example 6: ShopeePay E-Wallet Payment ===")
	shopeePayCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("EW-SHOPEEPAY-%d", 12345),
		Amount:      90000,
		PaymentType: pg.PaymentTypeShopeePay,
		Customer: pg.Customer{
			ID:    "CUST-006",
			Name:  "Diana Prince",
			Email: "diana@example.com",
			Phone: "+628123456789",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-007",
				Name:     "Product G",
				Price:    90000,
				Quantity: 1,
			},
		},
		CallbackURL: "https://example.com/callback",
	})
	if err != nil {
		log.Printf("Failed to create ShopeePay charge: %v", err)
	} else {
		log.Printf("ShopeePay Charge created:")
		log.Printf("  Transaction ID: %s", shopeePayCharge.TransactionID)
		log.Printf("  Payment URL: %s", shopeePayCharge.PaymentURL)
		log.Printf("  Status: %s", shopeePayCharge.Status)
	}

	// Example 7: Get payment status
	fmt.Println("\n=== Example 7: Get Payment Status ===")
	orderID := "INV-12345" // Use a real order ID from your Xendit dashboard
	status, err := client.GetStatus(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get status: %v", err)
	} else {
		log.Printf("Payment Status:")
		log.Printf("  Transaction ID: %s", status.TransactionID)
		log.Printf("  Order ID: %s", status.OrderID)
		log.Printf("  Status: %s", status.Status)
		log.Printf("  Amount: %d", status.Amount)
		log.Printf("  Paid Amount: %d", status.PaidAmount)
	}

	// Example 8: Cancel a transaction
	fmt.Println("\n=== Example 8: Cancel Transaction ===")
	err = client.Cancel(ctx, orderID)
	if err != nil {
		log.Printf("Cancel failed: %v", err)
	} else {
		log.Printf("Transaction cancelled successfully")
	}
}
