package model

type Token struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int64  `gorm:"column:user_id"`
	DeviceID     string `gorm:"column:device_id"`
	RefreshToken string `gorm:"column:refresh_token"`
	ExpiredAt    int64  `gorm:"column:expired_at"`
	CreatedAt    int64  `gorm:"column:created_at;autoCreateTime"`
}

func (Token) TableName() string {
	return "token"
}
