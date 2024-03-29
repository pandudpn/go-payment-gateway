<!-- markdownlint-disable MD014 MD024 MD026 MD033 MD036 MD041 -->

# GO Payment Gateway

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/pandudpn/go-payment-gateway)
[![Coverage Status](https://coveralls.io/repos/github/pandudpn/go-payment-gateway/badge.svg?branch=master&kill_cache=1)](https://coveralls.io/github/pandudpn/go-payment-gateway?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pandudpn/go-payment-gateway)](https://goreportcard.com/report/github.com/pandudpn/go-payment-gateway)

SDK GO for Payment Gateway in Indonesia. Currently only supports [Midtrans Core API](https://api-docs.midtrans.com/), [Xendit API](https://developers.xendit.co/api-reference), and [OY! Indonesia](https://api-docs.oyindonesia.com/).

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

| Payment Channels                              | Midtrans (Core API) | Xendit (Core API)  | OY! Indonesia (Core API) |
|-----------------------------------------------|:-------------------:|:------------------:|:------------------------:|
| Gopay (non-tokenization)                      | :white_check_mark:  |        :x:         |           :x:            |
| Gopay (tokenization)                          | :white_check_mark:  |        :x:         |           :x:            |
| OVO (non-tokenization)                        |         :x:         | :white_check_mark: |       :hourglass:        |
| OVO (tokenization)                            |         :x:         | :white_check_mark: |           :x:            |
| ShopeePay (non-tokenization)                  | :white_check_mark:  | :white_check_mark: |       :hourglass:        |
| ShopeePay (tokenization)                      |         :x:         | :white_check_mark: |           :x:            |
| DANA (non-tokenization)                       |         :x:         | :white_check_mark: |       :hourglass:        |
| DANA (tokenization)                           |         :x:         |        :x:         |           :x:            |
| LinkAja (non-tokenization)                    |         :x:         | :white_check_mark: |       :hourglass:        |
| LinkAja (tokenization)                        |         :x:         |        :x:         |           :x:            |

### Credit Card:

| Payment Channels                                  | Midtrans (Core API) | Xendit (Core API)  | OY! Indonesia (Core API) |
|---------------------------------------------------|:-------------------:|:------------------:|:------------------------:|
| Credit or Debit Card                              | :white_check_mark:  |    :hourglass:     |           :x:            |

### Virtual Account or Bank Transfer:

| Payment Channels                                  | Midtrans (Core API) | Xendit (Core API)  | OY! Indonesia (Core API) |
|---------------------------------------------------|:-------------------:|:------------------:|:------------------------:|
| BCA Virtual Account (Open Amount)                 | :white_check_mark:  | :white_check_mark: |       :hourglass:        |
| BCA Virtual Account (Closed Amount)               |         :x:         | :white_check_mark: |       :hourglass:        |
| BNI Virtual Account (Open Amount)                 | :white_check_mark:  | :white_check_mark: |       :hourglass:        |
| BNI Virtual Account (Closed Amount)               |         :x:         | :white_check_mark: |       :hourglass:        |
| BRI Virtual Account (Open Amount)                 | :white_check_mark:  | :white_check_mark: |       :hourglass:        |
| BRI Virtual Account (Closed Amount)               |         :x:         | :white_check_mark: |       :hourglass:        |
| Mandiri Virtual Account (Open Amount)             | :white_check_mark:  | :white_check_mark: |       :hourglass:        |
| Mandiri Virtual Account (Closed Amount)           |         :x:         | :white_check_mark: |       :hourglass:        |
| Permata Virtual Account (Open Amount)             | :white_check_mark:  | :white_check_mark: |       :hourglass:        |
| Permata Virtual Account (Closed Amount)           |         :x:         | :white_check_mark: |       :hourglass:        |
| BJB Virtual Account (Open Amount)                 |         :x:         | :white_check_mark: |           :x:            |
| BJB Virtual Account (Closed Amount)               |         :x:         | :white_check_mark: |           :x:            |
| BSI Virtual Account (Open Amount)                 |         :x:         | :white_check_mark: |           :x:            |
| BSI Virtual Account (Closed Amount)               |         :x:         | :white_check_mark: |           :x:            |
| CIMB Virtual Account (Open Amount)                |         :x:         | :white_check_mark: |       :hourglass:        |
| CIMB Virtual Account (Closed Amount)              |         :x:         | :white_check_mark: |       :hourglass:        |
| DBS Virtual Account (Open Amount)                 |         :x:         | :white_check_mark: |           :x:            |
| DBS Virtual Account (Closed Amount)               |         :x:         | :white_check_mark: |           :x:            |
| Sahabat Sampoerna Virtual Account (Open Amount)   |         :x:         | :white_check_mark: |           :x:            |
| Sahabat Sampoerna Virtual Account (Closed Amount) |         :x:         | :white_check_mark: |           :x:            |
| BTPN Virtual Account (Open Amount)                |         :x:         |        :x:         |       :hourglass:        |
| BTPN Virtual Account (Closed Amount)              |         :x:         |        :x:         |       :hourglass:        |

### Retail Outlets

<details>

  <summary>Notes</summary>

> Payment channel Alfamart, customer can complete their payment at various stores under ***Alfa Group***.
>
> The supported stores are: **Alfamart**, **Alfamidi** and **DAN+DAN**

</details>

| Payment Channels | Midtrans (Core API) | Xendit (Core API) | OY! Indonesia (Core API) |
|------------------|:-------------------:|:-----------------:|:------------------------:|
| Alfamart         |     :hourglass:     |    :hourglass:    |           :x:            |
| Indomaret        |     :hourglass:     |    :hourglass:    |           :x:            |

### Cardless Credit or Paylater

| Payment Channels | Midtrans (Core API) | Xendit (Core API) | OY! Indonesia (Core API) |
|------------------|:-------------------:|:-----------------:|:------------------------:|
| Akulaku          |     :hourglass:     |    :hourglass:    |           :x:            |
| Kredivo          |     :hourglass:     |    :hourglass:    |           :x:            |
| UangMe           |         :x:         |    :hourglass:    |           :x:            |
| IndoDana         |         :x:         |    :hourglass:    |           :x:            |
| Atome            |         :x:         |    :hourglass:    |           :x:            |

## License

MIT. Copyright 2022 by [pandudpn](LICENSE)
