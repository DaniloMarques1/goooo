package util

type ApiError struct {
	Msg    string
	Status int
}

func NewApiError(msg string, status int) error {
	return &ApiError{Msg: msg, Status: status}
}

func (ae *ApiError) Error() string {
	return ae.Msg
}
