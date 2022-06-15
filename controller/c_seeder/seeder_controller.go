package c_seeder

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SeederController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
