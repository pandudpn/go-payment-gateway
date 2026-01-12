<!-- markdownlint-disable MD014 MD024 MD026 MD033 MD036 MD041 -->

# GO Payment Gateway

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/pandudpn/go-payment-gateway)
[![Coverage Status](https://coveralls.io/repos/github/pandudpn/go-payment-gateway/badge.svg?branch=master&kill_cache=1)](https://coveralls.io/github/pandudpn/go-payment-gateway?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pandudpn/go-payment-gateway)](https://goreportcard.com/report/github.com/pandudpn/go-payment-gateway)

Unified Go SDK for Indonesian payment gateways. Supports [Midtrans](https://api-docs.midtrans.com/), [Xendit](https://developers.xendit.co/api-reference), and [Doku](https://developers.doku.com).

## Features

- **Unified API** - Single interface for multiple payment providers
- **Multiple Payment Channels** - E-Wallets, Virtual Accounts, Credit Cards, Retail Outlets
- **Type Safe** - Full type definitions with enums for payment types and statuses
- **Context Support** - Built-in context for timeout and cancellation
- **Webhook Handling** - Verify and parse webhook notifications
- **Sandbox/Production** - Easy environment switching

## Installation

```bash
go get github.com/pandudpn/go-payment-gateway
```

## Quick Start

### Midtrans

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/pandudpn/go-payment-gateway"
    "github.com/pandudpn/go-payment-gateway/pg"
)

func main() {
    client, err := pg.NewClient(
        pg.WithProvider("midtrans"),
        pg.WithServerKey("SB-Mid-server-xxx"),
        pg.WithClientKey("SB-Mid-client-xxx"),
        pg.WithEnvironment("sandbox"),
    )
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.CreateCharge(context.Background(), pg.ChargeParams{
        OrderID:     "ORDER-001",
        Amount:      50000,
        PaymentType: pg.PaymentTypeGoPay,
        Customer: pg.Customer{
            ID:    "CUST-001",
            Name:  "John Doe",
            Email: "john@example.com",
            Phone: "+628123456789",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Payment URL:", resp.PaymentURL)
}
```

### Xendit

```go
client, err := pg.NewClient(
    pg.WithProvider("xendit"),
    pg.WithServerKey("xnd_development_xxx"),
    pg.WithClientKey("xnd_development_xxx"),
    pg.WithEnvironment("sandbox"),
)
```

### Doku

```go
client, err := pg.NewClient(
    pg.WithProvider("doku"),
    pg.WithServerKey("SB-MID-server-xxx"),
    pg.WithClientKey("BRN-02201-xxx"),
    pg.WithEnvironment("sandbox"),
)
```

### Environment Variables

```bash
export PG_PROVIDER=midtrans
export PG_SERVER_KEY=SB-Mid-server-xxx
export PG_CLIENT_KEY=SB-Mid-client-xxx
export PG_ENV=sandbox
```

```go
client, err := pg.NewClient()  // Loads from environment
```

### Check Payment Status

```go
status, err := client.GetStatus(context.Background(), "ORDER-001")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Status:", status.Status)  // PENDING, SUCCESS, FAILED, etc.
```

### Cancel Transaction

```go
err := client.Cancel(context.Background(), "ORDER-001")
```

### Handle Webhook

```go
func handleWebhook(w http.ResponseWriter, r *http.Request) {
    event, err := client.ParseWebhook(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Println("Order ID:", event.OrderID)
    fmt.Println("Status:", event.Status)
    fmt.Println("Amount:", event.Amount)

    w.WriteHeader(http.StatusOK)
}
```

---

<details>
<summary><b>View table of contents</b></summary>

- [Payment Channels Supported](#payment-channels-supported)
  - [E-Wallet](#e-wallets)
  - [Credit Card](#credit-card)
  - [Virtual Account](#virtual-account-or-bank-transfer)
  - [Retail Outlets](#retail-outlets)
  - [Paylater](#cardless-credit-or-paylater)

</details>

---

## Payment Channels Supported

We're supporting the payments:

### E-Wallets:

| Payment Channels                              | Midtrans | Xendit | Doku |
|-----------------------------------------------|:--------:|:------:|:----:|
| GoPay                                         |   :white_check_mark:   |   :x:    | :white_check_mark: |
| OVO                                           |     :x:     | :white_check_mark: | :white_check_mark: |
| ShopeePay                                     |   :white_check_mark:  | :white_check_mark: | :white_check_mark: |
| DANA                                          |     :x:     | :white_check_mark: | :white_check_mark: |
| LinkAja                                       |     :x:     | :white_check_mark: | :white_check_mark: |
| QRIS                                          |   :white_check_mark:  | :white_check_mark: | :white_check_mark: |

### Credit Card:

| Payment Channels    | Midtrans | Xendit | Doku |
|---------------------|:--------:|:------:|:----:|
| Credit/Debit Card   | :white_check_mark:  |    :hourglass:     | :white_check_mark: |

### Virtual Account:

| Payment Channels      | Midtrans | Xendit | Doku |
|-----------------------|:--------:|:------:|:----:|
| BCA VA                | :white_check_mark:  | :white_check_mark: | :white_check_mark: |
| BNI VA                | :white_check_mark:  | :white_check_mark: | :white_check_mark: |
| BRI VA                | :white_check_mark:  | :white_check_mark: | :white_check_mark: |
| Mandiri VA            | :white_check_mark:  | :white_check_mark: | :white_check_mark: |
| Permata VA            | :white_check_mark:  | :white_check_mark: | :white_check_mark: |
| CIMB VA               | :x:     | :white_check_mark: | :white_check_mark: |
| BSI VA                | :x:     | :white_check_mark: | :x: |

### Retail Outlets:

<details>

  <summary>Notes</summary>

> Payment channel Alfamart, customer can complete their payment at various stores under ***Alfa Group***.
>
> The supported stores are: **Alfamart**, **Alfamidi** and **DAN+DAN**

</details>

| Payment Channels | Midtrans | Xendit | Doku |
|------------------|:--------:|:------:|:----:|
| Alfamart         | :white_check_mark:     |    :hourglass:    | :white_check_mark: |
| Indomaret        | :white_check_mark:     |    :hourglass:    | :x: |

### Paylater:

| Payment Channels | Midtrans | Xendit | Doku |
|------------------|:--------:|:------:|:----:|
| Akulaku          | :white_check_mark:     |    :hourglass:    | :x: |
| Kredivo          | :white_check_mark:     |    :hourglass:    | :x: |
| UangMe           | :x:         |    :hourglass:    | :x: |

## License

MIT. Copyright 2022 by [pandudpn](LICENSE)
