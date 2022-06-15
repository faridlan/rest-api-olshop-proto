package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/faridlan/rest-api-olshop-proto/app"
	"github.com/faridlan/rest-api-olshop-proto/controller/c_product"
	"github.com/faridlan/rest-api-olshop-proto/controller/c_seeder"
	"github.com/faridlan/rest-api-olshop-proto/helper"
	"github.com/faridlan/rest-api-olshop-proto/model/domain"
	"github.com/faridlan/rest-api-olshop-proto/repository/r_product"
	"github.com/faridlan/rest-api-olshop-proto/repository/r_uuid"
	"github.com/faridlan/rest-api-olshop-proto/service/s_product"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func setupDBTest() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/olshop_proto_test?parseTime=true")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {

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

	return router

}

func truncateProducts(db *sql.DB) {
	db.Exec("TRUNCATE products")
}

func TestSeederCreate(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/seeder", nil)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestSeederDelete(t *testing.T) {
	db := setupDBTest()
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/seeder", nil)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestCreateProductSuccess(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`
	{"name": "productA",
	"price": 99999,
	"quantity": 99,
	"image_url": "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png"}
	`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/products", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestCreateProductFailed(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`
	{"name": "",
	"price": 99999,
	"quantity": 99,
	"image_url": "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png"}
	`)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/products", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestUpdateProductSuccess(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	epoch := helper.EpochTime()

	tx, _ := db.Begin()
	uuidRepository := r_uuid.NewUuidRepository()
	uuid, _ := uuidRepository.CreteUui(context.Background(), tx)
	productRepository := r_product.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, domain.Product{
		IdProduct: uuid.Uuid,
		Name:      "productB",
		Price:     99999,
		Quantity:  99,
		ImageUrl:  "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png",
		CreatedAt: epoch,
	})
	tx.Commit()

	requestBody := strings.NewReader(`
	{"name": "productC",
	"price": 99999,
	"quantity": 100,
	"image_url": "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png"}
	`)

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/products/"+product.IdProduct, requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, product.IdProduct, responseBody["data"].(map[string]interface{})["id_product"])
	assert.Equal(t, "productC", responseBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, 99999, int(responseBody["data"].(map[string]interface{})["price"].(float64)))
	assert.Equal(t, 100, int(responseBody["data"].(map[string]interface{})["quantity"].(float64)))
	assert.Equal(t, "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png", responseBody["data"].(map[string]interface{})["image_url"])
	assert.Equal(t, product.CreatedAt, int64(responseBody["data"].(map[string]interface{})["created_at"].(float64)))
}

func TestUpdateProductFailed(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	epoch := helper.EpochTime()

	tx, _ := db.Begin()
	uuidRepository := r_uuid.NewUuidRepository()
	uuid, _ := uuidRepository.CreteUui(context.Background(), tx)
	productRepository := r_product.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, domain.Product{
		IdProduct: uuid.Uuid,
		Name:      "productB",
		Price:     99999,
		Quantity:  99,
		ImageUrl:  "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png",
		CreatedAt: epoch,
	})
	tx.Commit()

	requestBody := strings.NewReader(`
	{"name": "",
	"price": 99999,
	"quantity": 100,
	"image_url": "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png"}
	`)

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/products/"+product.IdProduct, requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestDeleteProductSuccess(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	epoch := helper.EpochTime()

	tx, _ := db.Begin()
	uuidRepository := r_uuid.NewUuidRepository()
	uuid, _ := uuidRepository.CreteUui(context.Background(), tx)
	productRepository := r_product.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, domain.Product{
		IdProduct: uuid.Uuid,
		Name:      "productB",
		Price:     99999,
		Quantity:  99,
		ImageUrl:  "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png",
		CreatedAt: epoch,
	})
	tx.Commit()

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/products/"+product.IdProduct, nil)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestDeleteProductFailed(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/products/404", nil)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestGetProductSuccess(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	epoch := helper.EpochTime()

	tx, _ := db.Begin()
	uuidRepository := r_uuid.NewUuidRepository()
	uuid, _ := uuidRepository.CreteUui(context.Background(), tx)
	productRepository := r_product.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, domain.Product{
		IdProduct: uuid.Uuid,
		Name:      "productB",
		Price:     99999,
		Quantity:  99,
		ImageUrl:  "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png",
		CreatedAt: epoch,
	})
	tx.Commit()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/products/"+product.IdProduct, nil)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, product.IdProduct, responseBody["data"].(map[string]interface{})["id_product"])
	assert.Equal(t, "productB", responseBody["data"].(map[string]interface{})["name"])
	assert.Equal(t, 99999, int(responseBody["data"].(map[string]interface{})["price"].(float64)))
	assert.Equal(t, 99, int(responseBody["data"].(map[string]interface{})["quantity"].(float64)))
	assert.Equal(t, "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png", responseBody["data"].(map[string]interface{})["image_url"])
	assert.Equal(t, product.CreatedAt, int64(responseBody["data"].(map[string]interface{})["created_at"].(float64)))
}

func TestGetProductFailed(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/products/404", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestListProducts(t *testing.T) {
	db := setupDBTest()
	truncateProducts(db)
	router := setupRouter(db)

	epoch := helper.EpochTime()

	tx, _ := db.Begin()
	uuidRepository := r_uuid.NewUuidRepository()
	uuid, _ := uuidRepository.CreteUui(context.Background(), tx)
	uuid2, _ := uuidRepository.CreteUui(context.Background(), tx)
	productRepository := r_product.NewProductRepository()
	product1 := productRepository.Save(context.Background(), tx, domain.Product{
		IdProduct: uuid.Uuid,
		Name:      "productB",
		Price:     99999,
		Quantity:  99,
		ImageUrl:  "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png",
		CreatedAt: epoch,
	})
	product2 := productRepository.Save(context.Background(), tx, domain.Product{
		IdProduct: uuid2.Uuid,
		Name:      "productC",
		Price:     10000,
		Quantity:  100,
		ImageUrl:  "http://sg1.products.digital.com/zzzztttyyyyy.png",
		CreatedAt: epoch,
	})
	tx.Commit()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/products", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	var products = responseBody["data"].([]interface{})

	productResponse1 := products[0].(map[string]interface{})
	productResponse2 := products[1].(map[string]interface{})

	assert.Equal(t, product1.IdProduct, productResponse1["id_product"])
	assert.Equal(t, "productB", productResponse1["name"])
	assert.Equal(t, 99999, int(productResponse1["price"].(float64)))
	assert.Equal(t, 99, int(productResponse1["quantity"].(float64)))
	assert.Equal(t, "http://sg1.products.digital.com/yyetruetgdshghfsgdfjs.png", productResponse1["image_url"])
	assert.Equal(t, product1.CreatedAt, int64(productResponse1["created_at"].(float64)))

	assert.Equal(t, product2.IdProduct, productResponse2["id_product"])
	assert.Equal(t, "productC", productResponse2["name"])
	assert.Equal(t, 10000, int(productResponse2["price"].(float64)))
	assert.Equal(t, 100, int(productResponse2["quantity"].(float64)))
	assert.Equal(t, "http://sg1.products.digital.com/zzzztttyyyyy.png", productResponse2["image_url"])
	assert.Equal(t, product2.CreatedAt, int64(productResponse2["created_at"].(float64)))
}
