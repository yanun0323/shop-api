package repository

import (
	"main/internal/domain/entity"
	"main/internal/domain/repository"

	"github.com/redis/go-redis/v9"
	"github.com/yanun0323/pkg/logs"
)

type notificationRepository struct {
	rdb *redis.Client
}

func NewNotificationRepository(rdb *redis.Client) repository.NotificationRepository {
	return &notificationRepository{
		rdb: rdb,
	}
}

func (repo *notificationRepository) SendEmail(email string, content entity.Notification) error {
	// - check rate limit from redis
	// - send email
	logs.Debugf("send email:\n<receiver> %s\n<subject> %s\n<body> %s", email, content.Subject, content.Body)
	// - update rate limit
	return nil
}

func (repo *notificationRepository) SendSMS(countryCode, phone string, content entity.Notification) error {
	// - check rate limit from redis
	// - send SMS
	logs.Debugf("send sms: <receiver> (%s)%s\n<subject> %s\n<body> %s", countryCode, phone, content.Subject, content.Body)
	// - update rate limit
	return nil
}
