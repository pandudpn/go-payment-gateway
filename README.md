<!-- markdownlint-disable MD014 MD024 MD026 MD033 MD036 MD041 -->

# GO Payment Gateway

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/pandudpn/go-payment-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/pandudpn/go-payment-gateway)](https://goreportcard.com/report/github.com/pandudpn/go-payment-gateway)

SDK GO for Payment Gateway in Indonesia. Currently only supports [Midtrans Core API](https://api-docs.midtrans.com/), [Xendit API](https://developers.xendit.co/api-reference), and [OY! Indonesia](https://api-docs.oyindonesia.com/).

---

<details>
<summary><b>View table of contents</b></summary>

- [Payment Channels Supported](#payment-channels-supported)

</details>

---

## Payment Channels Supported

We're supporting the payments:

- Bank Transfer via Virtual Account (BCA, BNI, BRI, Mandiri, Permata, and other's Bank) with closed or open amount.
- EWallet (GOPAY, OVO, DANA, LinkAja, ShopeePay) with OneTime Payment or Linked Account.
- Credit or Debit Card with/without installment.
- Recurring payment with Credit Card, GOPAY.
- Retail Outlet (Alfamart, Indomaret, Alfamidi, and other's RO).
- Cardless Credit (Kredivo, Akulaku, Indodana, and other's credit).

## License

MIT. Copyright 2022 by [pandudpn](LICENSE)
