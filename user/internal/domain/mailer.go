package domain

type SendAuthCodeInfo struct {
	To      []string
	Subject string
	Code    int
}
