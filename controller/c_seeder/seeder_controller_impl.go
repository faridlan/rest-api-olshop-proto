package c_seeder

import (
	"embed"
	"encoding/json"
	"net/http"

	"github.com/faridlan/rest-api-olshop-proto/helper"
	"github.com/faridlan/rest-api-olshop-proto/model/web"
	"github.com/faridlan/rest-api-olshop-proto/model/web/w_product"
	"github.com/faridlan/rest-api-olshop-proto/service/s_product"
	"github.com/julienschmidt/httprouter"
)

type SeederControllerImpl struct {
	ProductService s_product.Service
}

func NewSeederController(productService s_product.Service) SeederController {
	return &SeederControllerImpl{
		ProductService: productService,
	}
}

//go:embed json/product.json

var Json embed.FS

func (controller *SeederControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	product, err := Json.ReadFile("json/product.json")
	helper.PanicIfError(err)
	productCreate := []w_product.CreateRequest{}
	err = json.Unmarshal(product, &productCreate)
	helper.PanicIfError(err)

	productResponse := controller.ProductService.CreateMultiple(request.Context(), productCreate)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   productResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *SeederControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	controller.ProductService.DeleteTable(request.Context())

	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}
