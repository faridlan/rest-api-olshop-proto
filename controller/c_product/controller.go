package c_product

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Controller interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UploadImg(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
