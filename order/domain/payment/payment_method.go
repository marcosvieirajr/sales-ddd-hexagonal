package payment

// Method represents the payment method chosen by the customer.
type Method int

const (
	MethodCreditCard   Method = iota // MethodCreditCard represents payment by credit card.
	MethodDebitCard                  // MethodDebitCard represents payment by debit card.
	MethodCash                       // MethodCash represents payment in cash.
	MethodPix                        // MethodPix represents payment via Pix instant transfer.
	MethodBankTransfer               // MethodBankTransfer represents payment via bank transfer (TED/DOC).
	MethodBancSlip                   // MethodBancSlip represents payment via bank slip (boleto bancário).
)
