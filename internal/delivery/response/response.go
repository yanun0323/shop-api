package response

import (
	"fmt"
	"main/internal/domain/usecase"
	"main/internal/helper/pager"
	"strings"

	"github.com/pkg/errors"

	"github.com/yanun0323/pkg/logs"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DataResponse struct {
	Response

	Data any `json:"data"`
}

type PageDataResponse struct {
	Response

	Data any             `json:"data"`
	Page *pager.Response `json:"page"`
}

func Msg(message string) Response {
	return Response{
		Message: message,
	}
}

func Data(data any, message string) *DataResponse {
	return &DataResponse{
		Response: Msg(message),
		Data:     data,
	}
}

func PagedData(data any, page *pager.Response, message string) *PageDataResponse {
	return &PageDataResponse{
		Response: Msg(message),
		Data:     data,
		Page:     page,
	}
}

func MsgErr(message string, logMsgArgs ...any) Response {
	return Err(errors.New("internal error"), message, logMsgArgs...)
}

func Err(err error, message string, logMsgArgs ...any) Response {
	var (
		code       int = usecase.ErrInternal.Code
		usecaseErr *usecase.UsecaseError
	)

	if errors.As(err, &usecaseErr) {
		code = usecaseErr.Code
	}

	logs.Error(logMsg(err, message, logMsgArgs...))
	return Response{
		Code:    code,
		Message: message,
	}
}

func logMsg(err error, message string, logMsgArgs ...any) string {
	switch len(logMsgArgs) {
	case 0:
		return fmt.Sprintf("%s, err: %v", message, err)
	case 1:
		return fmt.Sprintf("%s, %s, err: %v", message, logMsgArgs[0], err)
	default:
		format, ok := logMsgArgs[0].(string)
		if ok {
			logMsgArgs = logMsgArgs[1:]
		} else {
			format = strings.Repeat("%v ", len(logMsgArgs))
		}

		logMsg := fmt.Sprintf(format, logMsgArgs...)
		return fmt.Sprintf("%s, %s, err: %v", message, logMsg, err)
	}
}
