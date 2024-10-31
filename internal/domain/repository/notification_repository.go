package repository

import "main/internal/domain/entity"

//go:generate domaingen -destination=../../repository/notification_repository.go -package=repository -constructor
type NotificationRepository interface {
	SendEmail(email string, content entity.Notification) error
	SendSMS(countryCode, phone string, content entity.Notification) error
}
