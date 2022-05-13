package dto

type CalendarDto struct {
	ExpiresIn int64 `json:"expires_in"`
}

type AdditionalInfoDto struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateCobDto struct {
	Value           float64             `json:"value"`
	Cal             CalendarDto         `json:"calendar"`
	AdditionalInfos []AdditionalInfoDto `json:"additional_info"`
}
