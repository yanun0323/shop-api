package payload

import "net/http"

const (
	_headerDeviceIDKey = "X-Device-ID"
)

func GetDeviceID(r *http.Request) string {
	return r.Header.Get(_headerDeviceIDKey)
}
