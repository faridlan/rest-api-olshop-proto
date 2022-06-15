package helper

import (
	"net/http"
	"strconv"

	"github.com/faridlan/rest-api-olshop-proto/model/domain"
)

func NewPagination(request *http.Request) domain.Pagination {
	queryParam := request.URL.Query()
	var pagination domain.Pagination
	if len(queryParam) == 0 {
		pagination.Page = 0
		pagination.Limit = 10
		return pagination
	} else {
		page, err := strconv.Atoi(queryParam.Get("page"))
		PanicIfError(err)
		limit, err := strconv.Atoi(queryParam.Get("limit"))
		PanicIfError(err)

		pagination.Page = (page - 1) * limit
		pagination.Limit = limit
		return pagination
	}
}
