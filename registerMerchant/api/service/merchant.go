package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/danilomarques1/godemo/registerMerchant/api/consumer"
	"github.com/danilomarques1/godemo/registerMerchant/api/model"
)

type MerchantService interface {
	Save(merchant *model.Merchant) error
	FindAll() ([]model.Merchant, error)
	FindById(merchantId string) (*model.Merchant, error)
}

type MerchantServiceImpl struct {
	merchantRepository model.MerchantRepository
}

func NewMerchantServiceImpl(merchantRepository model.MerchantRepository) *MerchantServiceImpl {
	return &MerchantServiceImpl{merchantRepository: merchantRepository}
}

func (ms *MerchantServiceImpl) Save(merchant *model.Merchant) error {
	if _, err := ms.merchantRepository.FindById(merchant.MerchantId); err == nil {
		return errors.New("Merchant already registered")
	}

	if err := ms.merchantRepository.Save(merchant); err != nil {
		return err
	}
	return nil
}

func (ms *MerchantServiceImpl) ConsumeAndSave(c consumer.Consumer) {
	for {
		b := c.Consume()
		log.Printf("Mensagem consumida %v\n", string(b))

		if len(b) > 0 {
			merchant := &model.Merchant{}
			err := json.Unmarshal(b, merchant)
			if err == nil {
				log.Printf("Salvando merchant em banco %v\n", merchant)
				err := ms.Save(merchant)
				if err != nil {
					log.Printf("Error saving data %v\n", err)
				}

			} else {
				log.Printf("Error unmarshaling data = %v\n", err)
			}
		}
	}
}
