package c_product

import (
	"io/ioutil"
	"net/http"

	"github.com/faridlan/rest-api-olshop-proto/helper"
	"github.com/faridlan/rest-api-olshop-proto/model/web"
	"github.com/faridlan/rest-api-olshop-proto/model/web/w_product"
	"github.com/faridlan/rest-api-olshop-proto/service/s_product"
	"github.com/julienschmidt/httprouter"
)

type ControllerImpl struct {
	ProductService s_product.Service
}

func NewProductController(productService s_product.Service) Controller {
	return &ControllerImpl{
		ProductService: productService,
	}
}

func (controller *ControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	createRequest := w_product.CreateRequest{}
	helper.ReadFromRequestBody(request, &createRequest)

	productResponse := controller.ProductService.Create(request.Context(), createRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   productResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idProduct := params.ByName("productId")
	updateRequest := w_product.UpdateRequest{}
	helper.ReadFromRequestBody(request, &updateRequest)

	updateRequest.IdProduct = idProduct

	productResponse := controller.ProductService.Update(request.Context(), updateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   productResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idProduct := params.ByName("productId")

	controller.ProductService.Delete(request.Context(), idProduct)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	idProduct := params.ByName("productId")

	productResponse := controller.ProductService.FindById(request.Context(), idProduct)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   productResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pagination := helper.NewPagination(request)

	productResponses := controller.ProductService.Findall(request.Context(), pagination)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   productResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ControllerImpl) UploadImg(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	err := request.ParseMultipartForm(10 << 20)
	helper.PanicIfError(err)

	file, _, err := request.FormFile("productImage")
	helper.PanicIfError(err)
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	helper.PanicIfError(err)

	image := w_product.CreateRequest{
		ImageUrl: string(fileBytes),
	}

	productResponse := controller.ProductService.CreateImg(request.Context(), image)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   productResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
