package entity

type NotifyType int

const (
	NotifyTypeRegisterOtp NotifyType = iota + 1
	NotifyTypeLoginOtp
)

type Notification struct {
	Type NotifyType `json:"type"`

	Subject string `json:"subject"`
	Body    string `json:"body"`
}
