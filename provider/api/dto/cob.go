package dto

type CalendarDto struct {
	ExpiresIn int64 `validate:"gt=0" json:"expires_in" `
}

type AdditionalInfoDto struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TODO key wont be uuid always, fix based on key type
// key validation might be:
// uuid (RANDOMKEY)
// len=11 (MOBILEPHONE, NATIONALID)
// len=16 (NATIONALID)

type CreateCobDto struct {
	Value           float64             `json:"value" validate:"required,gt=0"`
	KeyType         string              `json:"key_type" validate:"oneof=MOBILEPHONE RANDOMKEY NATIONALID"`
	Key             string              `json:"key" validate:"uuid_if=KeyType=MOBILEPHONE"`
	Cal             CalendarDto         `json:"calendar" validate:"dive"`
	AdditionalInfos []AdditionalInfoDto `json:"additional_info" validate:"dive"`
}
