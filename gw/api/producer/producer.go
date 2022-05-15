package producer

import (
	"encoding/json"
)

type Merchant struct {
	SubAcquirerId   string `json:"sub_acquirer_id,omitempty"`
	SubAcquirerName string `json:"sub_acquirer_name,omitempty"`
	MerchantId      string `json:"merchant_id,omitempty"`
	MerchantName    string `json:"merchant_name,omitempty"`
	MerchantAddress string `json:"merchant_address,omitempty"`
}

func (m *Merchant) Marshal() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Producer interface {
	Produce(merchant Merchant) error
	Close() error
}
