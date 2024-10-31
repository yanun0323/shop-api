package payload

import (
	"context"
)

type (
	tokenStruct    struct{}
	userIDStruct   struct{}
	deviceIDStruct struct{}
)

var (
	tokenKey    = tokenStruct{}
	userIDKey   = userIDStruct{}
	deviceIDKey = deviceIDStruct{}
)

func SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func GetToken(ctx context.Context) string {
	if token, ok := ctx.Value(tokenKey).(string); ok {
		return token
	}

	return ""
}

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) int64 {
	if userID, ok := ctx.Value(userIDKey).(int64); ok {
		return userID
	}

	return 0
}

func SetDeviceID(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, deviceIDKey, deviceID)
}
