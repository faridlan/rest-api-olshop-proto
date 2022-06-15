package app

import (
	"github.com/faridlan/rest-api-olshop-proto/controller/c_product"
	"github.com/faridlan/rest-api-olshop-proto/controller/c_seeder"
	"github.com/faridlan/rest-api-olshop-proto/exception"
	"github.com/julienschmidt/httprouter"
)

type Controller struct {
	ProductController c_product.Controller
	SeederController  c_seeder.SeederController
}

func NewRouter(controller Controller) *httprouter.Router {

	router := httprouter.New()

	//Product
	router.POST("/api/products", controller.ProductController.Create)
	router.POST("/api/image/products", controller.ProductController.UploadImg)
	router.PUT("/api/products/:productId", controller.ProductController.Update)
	router.DELETE("/api/products/:productId", controller.ProductController.Delete)
	router.GET("/api/products/:productId", controller.ProductController.FindById)
	router.GET("/api/products", controller.ProductController.FindAll)

	//Seeder
	router.POST("/api/seeder", controller.SeederController.Create)
	router.DELETE("/api/seeder", controller.SeederController.Delete)

	router.PanicHandler = exception.ErrorHandler
	return router
}
