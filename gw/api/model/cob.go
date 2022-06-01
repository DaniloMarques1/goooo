package model

import (
	"time"
)

const (
	ACTIVE          = "ATIVA"
	CONCLUDED       = "CONCLUIDA"
	REMOVED_BY_USER = "REMOVIDA_PELO_USUARIO_RECEBEDOR"
	REMOVED_BY_PSP  = "REMOVIDA_PELO_USUARIO_PSP"
)

type Calendar struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	ExpiresIn int64     `json:"expires_in" bson:"expires_in"`
}

type AdditionalInfo struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

type Cob struct {
	TxId            string           `json:"txid" bson:"_id"`
	Value           float64          `json:"value" bson:"value"`
	Status          string           `json:"status" bson:"status"`
	KeyType         string           `json:"key_type" bson:"key_type"`
	Key             string           `json:"key" bson:"key"`
	Cal             Calendar         `json:"calendar" bson:"calendar"`
	AdditionalInfos []AdditionalInfo `json:"additional_info" bson:"additional_info"`
}

type CobRepository interface {
	Save(cob *Cob) error
	FindById(txid string) (*Cob, error)
	Update(cob *Cob) error
}
