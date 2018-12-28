package config

type TransactionConfig struct {
	Status StatusTransaction
}

type StatusTransaction struct {
	Unpaid         int
	Paid           int
	Sending        int
	Received       int
	Done           int
	Cancel         int
	Refund         int
	PaymentExpired int
}

func NewTransactionConfig() *TransactionConfig {
	return &TransactionConfig{
		Status: StatusTransaction{
			Unpaid:         0,
			Paid:           1,
			Sending:        2,
			Received:       3,
			Done:           4,
			Cancel:         5,
			Refund:         6,
			PaymentExpired: 7,
		},
	}
}
