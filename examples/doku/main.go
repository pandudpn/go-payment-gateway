package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pg "github.com/pandudpn/go-payment-gateway"
)

func main() {
	// Initialize the Doku client
	client, err := pg.NewClient(
		pg.WithProvider("doku"),
		pg.WithServerKey(os.Getenv("DOKU_SHARED_KEY")),
		pg.WithClientKey(os.Getenv("DOKU_CLIENT_ID")),
		pg.WithEnvironment("sandbox"), // or "production"
	)
	if err != nil {
		log.Fatalf("Failed to create Doku client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Create a Virtual Account BCA payment
	fmt.Println("=== Example 1: Virtual Account BCA Payment ===")
	vaCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("DOKU-VA-BCA-%d", 12345),
		Amount:      100000,
		PaymentType: pg.PaymentTypeVABCA,
		Customer: pg.Customer{
			ID:    "CUST-001",
			Name:  "John Doe",
			Email: "john@example.com",
			Phone: "+628123456789",
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

	// Example 2: Create a Virtual Account Mandiri payment
	fmt.Println("\n=== Example 2: Virtual Account Mandiri Payment ===")
	vaMandiriCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("DOKU-VA-MANDIRI-%d", 12345),
		Amount:      150000,
		PaymentType: pg.PaymentTypeVAMandiri,
		Customer: pg.Customer{
			ID:    "CUST-002",
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Phone: "+628123456789",
		},
		CallbackURL: "https://example.com/callback",
	})
	if err != nil {
		log.Printf("Failed to create Mandiri VA charge: %v", err)
	} else {
		log.Printf("Mandiri VA Charge created:")
		log.Printf("  Transaction ID: %s", vaMandiriCharge.TransactionID)
		log.Printf("  VA Number: %s", vaMandiriCharge.VANumber)
		log.Printf("  VA Bank: %s", vaMandiriCharge.VABank)
		log.Printf("  Status: %s", vaMandiriCharge.Status)
	}

	// Example 3: Create a GoPay e-wallet payment
	fmt.Println("\n=== Example 3: GoPay E-Wallet Payment ===")
	gopayCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("DOKU-EW-GOPAY-%d", 12345),
		Amount:      50000,
		PaymentType: pg.PaymentTypeGoPay,
		Customer: pg.Customer{
			ID:    "CUST-003",
			Name:  "Bob Smith",
			Email: "bob@example.com",
			Phone: "+628123456789",
		},
		CallbackURL: "https://example.com/callback",
		ReturnURL:   "https://example.com/return",
	})
	if err != nil {
		log.Printf("Failed to create GoPay charge: %v", err)
	} else {
		log.Printf("GoPay Charge created:")
		log.Printf("  Transaction ID: %s", gopayCharge.TransactionID)
		log.Printf("  Payment URL: %s", gopayCharge.PaymentURL)
		log.Printf("  Status: %s", gopayCharge.Status)
	}

	// Example 4: Create an OVO e-wallet payment
	fmt.Println("\n=== Example 4: OVO E-Wallet Payment ===")
	ovoCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("DOKU-EW-OVO-%d", 12345),
		Amount:      75000,
		PaymentType: pg.PaymentTypeOVO,
		Customer: pg.Customer{
			ID:    "CUST-004",
			Name:  "Alice Johnson",
			Email: "alice@example.com",
			Phone: "+628123456789",
		},
		CallbackURL: "https://example.com/callback",
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
		OrderID:     fmt.Sprintf("DOKU-EW-DANA-%d", 12345),
		Amount:      80000,
		PaymentType: pg.PaymentTypeDANA,
		Customer: pg.Customer{
			ID:    "CUST-005",
			Name:  "Charlie Brown",
			Email: "charlie@example.com",
			Phone: "+628123456789",
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

	// Example 6: Create a LinkAja e-wallet payment
	fmt.Println("\n=== Example 6: LinkAja E-Wallet Payment ===")
	linkAjaCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("DOKU-EW-LINKAJA-%d", 12345),
		Amount:      85000,
		PaymentType: pg.PaymentTypeLinkAja,
		Customer: pg.Customer{
			ID:    "CUST-006",
			Name:  "Diana Prince",
			Email: "diana@example.com",
			Phone: "+628123456789",
		},
		CallbackURL: "https://example.com/callback",
	})
	if err != nil {
		log.Printf("Failed to create LinkAja charge: %v", err)
	} else {
		log.Printf("LinkAja Charge created:")
		log.Printf("  Transaction ID: %s", linkAjaCharge.TransactionID)
		log.Printf("  Payment URL: %s", linkAjaCharge.PaymentURL)
		log.Printf("  Status: %s", linkAjaCharge.Status)
	}

	// Example 7: Create a QRIS payment
	fmt.Println("\n=== Example 7: QRIS Payment ===")
	qrisCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("DOKU-QRIS-%d", 12345),
		Amount:      90000,
		PaymentType: pg.PaymentTypeQRIS,
		Customer: pg.Customer{
			ID:    "CUST-007",
			Name:  "Eve Anderson",
			Email: "eve@example.com",
			Phone: "+628123456789",
		},
		CallbackURL: "https://example.com/callback",
	})
	if err != nil {
		log.Printf("Failed to create QRIS charge: %v", err)
	} else {
		log.Printf("QRIS Charge created:")
		log.Printf("  Transaction ID: %s", qrisCharge.TransactionID)
		log.Printf("  Payment URL: %s", qrisCharge.PaymentURL)
		log.Printf("  Status: %s", qrisCharge.Status)
	}

	// Example 8: Get payment status
	fmt.Println("\n=== Example 8: Get Payment Status ===")
	orderID := "DOKU-VA-BCA-12345" // Use a real order ID
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

	// Example 9: Cancel a transaction
	// Note: Doku doesn't support direct cancellation, payments expire automatically
	fmt.Println("\n=== Example 9: Cancel Transaction ===")
	err = client.Cancel(ctx, orderID)
	if err != nil {
		log.Printf("Cancel not supported by Doku: %v", err)
	} else {
		log.Printf("Transaction cancelled")
	}
}
