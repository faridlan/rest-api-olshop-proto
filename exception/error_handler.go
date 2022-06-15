package exception

import (
	"net/http"

	"github.com/faridlan/rest-api-olshop-proto/helper"
	"github.com/faridlan/rest-api-olshop-proto/model/web"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if notFoundError(writer, request, err) {
		return
	}
	if validationErrors(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}

func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		writer.Header().Add("content-type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func validationErrors(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(ValidationError)
	if ok {
		writer.Header().Add("content-type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   exception.Error,
		}

		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

// func validationErrors(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
// 	exception, ok := err.(validator.ValidationErrors)
// 	defer log.Print(exception)
// 	if ok {
// 		writer.Header().Add("content-type", "application/json")
// 		writer.WriteHeader(http.StatusBadRequest)

// 		webResponse := web.WebResponse{
// 			Code:   http.StatusBadRequest,
// 			Status: "BAD REQUEST",
// 			Data:   exception.Error(),
// 		}

// 		helper.WriteToResponseBody(writer, webResponse)
// 		return true
// 	} else {
// 		return false
// 	}
// }

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	webResponse := web.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data:   err,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
