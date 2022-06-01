package provider

import (
	"github.com/danilomarques1/godemo/gw/api/dto"
	"github.com/danilomarques1/godemo/gw/api/model"
)

type Provider interface {
	CreateCob(token string, cobDto dto.CreateCobDto) (*model.Cob, error)
	FindCob(token, txid string) (*model.Cob, error)
	// Cancel(token, txid string) error // TODO implement inside itau provider
}
