package dto

type CalendarDto struct {
	ExpiresIn int64 `validate:"gt=0" json:"expires_in"`
}

type AdditionalInfoDto struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateCobDto struct {
	Value           float64             `json:"value" validate:"required,gt=0"`
	KeyType         string              `json:"key_type" validate:"oneof=MOBILEPHONE RANDOMKEY NATIONALID MERCHANTNATIONALID"`
	Key             string              `json:"key" validate:"pix-key"`
	Cal             CalendarDto         `json:"calendar" validate:"dive"`
	AdditionalInfos []AdditionalInfoDto `json:"additional_info" validate:"dive"`
}
