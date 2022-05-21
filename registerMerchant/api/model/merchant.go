package model

type Merchant struct {
	MerchantId      string `json:"merchant_id,omitempty"`
	MerchantName    string `json:"merchant_name,omitempty"`
	MerchantAddress string `json:"merchant_address,omitempty"`
	SubAcquirerId   string `json:"sub_acquirer_id,omitempty"`
	SubAcquirerName string `json:"sub_acquirer_name,omitempty"`
}

type MerchantRepository interface {
	Save(merchant *Merchant) error
	FindById(merchantId string) (*Merchant, error)
	FindAll() ([]Merchant, error)
}
