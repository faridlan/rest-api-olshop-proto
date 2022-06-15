package main

import (
	"net/http"

	"github.com/faridlan/rest-api-olshop-proto/app"
	"github.com/faridlan/rest-api-olshop-proto/controller/c_product"
	"github.com/faridlan/rest-api-olshop-proto/controller/c_seeder"
	"github.com/faridlan/rest-api-olshop-proto/helper"
	"github.com/faridlan/rest-api-olshop-proto/repository/r_product"
	"github.com/faridlan/rest-api-olshop-proto/repository/r_uuid"
	"github.com/faridlan/rest-api-olshop-proto/service/s_product"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := app.NewDatabase()
	validate := validator.New()
	uuidRepository := r_uuid.NewUuidRepository()

	//Product
	productRepo := r_product.NewProductRepository()
	productService := s_product.NewProductService(productRepo, uuidRepository, db, validate)
	productController := c_product.NewProductController(productService)

	seederController := c_seeder.NewSeederController(productService)

	controller := app.Controller{
		ProductController: productController,
		SeederController:  seederController,
	}

	router := app.NewRouter(controller)

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
