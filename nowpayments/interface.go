package nowpayments

import "github.com/kehiy/RoboPac/store"

type INowpayment interface {
	CreatePayment(party *store.TwitterParty) error
	UpdatePayment(party *store.TwitterParty) error
}
