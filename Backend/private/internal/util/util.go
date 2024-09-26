package util

type placeholder struct {
	IsPlaceholder bool
	Endpoint      string
}

func Placeholder(endpoint string) *placeholder {
	return &placeholder{true, endpoint}
}
