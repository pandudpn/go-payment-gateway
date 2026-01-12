package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pg "github.com/pandudpn/go-payment-gateway"
)

func main() {
	// Initialize the Midtrans client
	client, err := pg.NewClient(
		pg.WithProvider("midtrans"),
		pg.WithServerKey(os.Getenv("MIDTRANS_SERVER_KEY")),
		pg.WithClientKey(os.Getenv("MIDTRANS_CLIENT_KEY")),
		pg.WithEnvironment("sandbox"), // or "production"
		pg.WithSnap(),                  // Enable SNAP mode
	)
	if err != nil {
		log.Fatalf("Failed to create Midtrans client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Create a GoPay payment charge
	fmt.Println("=== Example 1: GoPay Payment ===")
	gopayCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("ORDER-GOPAY-%d", 12345),
		Amount:      50000,
		PaymentType: pg.PaymentTypeGoPay,
		Customer: pg.Customer{
			ID:    "CUST-001",
			Name:  "John Doe",
			Email: "john@example.com",
			Phone: "+628123456789",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-001",
				Name:     "Test Product",
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
		log.Printf("  Amount: %d", gopayCharge.Amount)
		log.Printf("  Status: %s", gopayCharge.Status)
		log.Printf("  Payment URL: %s", gopayCharge.PaymentURL)
	}

	// Example 2: Create a Virtual Account BCA payment charge
	fmt.Println("\n=== Example 2: Virtual Account BCA Payment ===")
	vaCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("ORDER-VA-BCA-%d", 12345),
		Amount:      100000,
		PaymentType: pg.PaymentTypeVABCA,
		Customer: pg.Customer{
			ID:    "CUST-002",
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Phone: "+628123456789",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-002",
				Name:     "Product B",
				Price:    100000,
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

	// Example 3: Create a QRIS payment charge
	fmt.Println("\n=== Example 3: QRIS Payment ===")
	qrisCharge, err := client.CreateCharge(ctx, pg.ChargeParams{
		OrderID:     fmt.Sprintf("ORDER-QRIS-%d", 12345),
		Amount:      75000,
		PaymentType: pg.PaymentTypeQRIS,
		Customer: pg.Customer{
			ID:    "CUST-003",
			Name:  "Bob Smith",
			Email: "bob@example.com",
			Phone: "+628123456789",
		},
		Items: []pg.Item{
			{
				ID:       "ITEM-003",
				Name:     "Product C",
				Price:    75000,
				Quantity: 1,
			},
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

	// Example 4: Get payment status
	fmt.Println("\n=== Example 4: Get Payment Status ===")
	orderID := "ORDER-GOPAY-12345" // Use a real order ID
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

	// Example 5: Cancel a transaction (if supported)
	fmt.Println("\n=== Example 5: Cancel Transaction ===")
	err = client.Cancel(ctx, orderID)
	if err != nil {
		log.Printf("Cancel failed (may not be supported): %v", err)
	} else {
		log.Printf("Transaction cancelled successfully")
	}
}
