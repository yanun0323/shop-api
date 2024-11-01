package payload

import (
	"main/internal/helper/pager"
	"net/http"
	"strconv"
)

const (
	_queryPageKey     = "page"
	_queryPageSizeKey = "page_size"
)

func GetPage(r *http.Request) pager.Request {
	pageStr := r.URL.Query().Get(_queryPageKey)
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := r.URL.Query().Get(_queryPageSizeKey)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	return pager.New(page, pageSize)
}
