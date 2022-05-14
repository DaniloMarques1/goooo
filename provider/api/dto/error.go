package dto

type ApiErrorDto struct {
	Message string  `json:"message"`
	Errors  []Error `json:"errors,omitempty"`
}

type Error struct {
	Field         string      `json:"field"`
	RejectedValue interface{} `json:"rejected_value"`
}
