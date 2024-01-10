package entity

type SendAuthCodeInfo struct {
	To      []string
	Subject string
	Code    int
}
