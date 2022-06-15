package s_product

import (
	"context"
	"database/sql"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/faridlan/rest-api-olshop-proto/exception"
	"github.com/faridlan/rest-api-olshop-proto/helper"
	"github.com/faridlan/rest-api-olshop-proto/helper/model/h_product"
	"github.com/faridlan/rest-api-olshop-proto/model/domain"
	"github.com/faridlan/rest-api-olshop-proto/model/web/w_product"
	"github.com/faridlan/rest-api-olshop-proto/repository/r_product"
	"github.com/faridlan/rest-api-olshop-proto/repository/r_uuid"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type ServiceImpl struct {
	ProductRepo r_product.Repository
	UuidRepo    r_uuid.Repository
	DB          *sql.DB
	Validate    *validator.Validate
}

func NewProductService(productRepo r_product.Repository, uuidRepo r_uuid.Repository, db *sql.DB, validate *validator.Validate) Service {
	return ServiceImpl{
		ProductRepo: productRepo,
		UuidRepo:    uuidRepo,
		DB:          db,
		Validate:    validate,
	}
}

func (service ServiceImpl) Create(ctx context.Context, request w_product.CreateRequest) w_product.Response {

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(service.Validate, trans)

	err := service.Validate.Struct(request)
	if err != nil {
		errs := exception.TranslateError(err, trans)
		var x []string
		for _, v := range errs {
			x = append(x, v.Error())
		}
		panic(exception.NewValidationError(x))
	}

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	epoch := helper.EpochTime()

	defer helper.CommitOrRollback(tx)

	uuid, err := service.UuidRepo.CreteUui(ctx, tx)
	helper.PanicIfError(err)

	product := domain.Product{
		IdProduct: uuid.Uuid,
		Name:      request.Name,
		Price:     request.Price,
		Quantity:  request.Quantity,
		ImageUrl:  request.ImageUrl,
		CreatedAt: epoch,
	}

	product = service.ProductRepo.Save(ctx, tx, product)
	return h_product.ToProductResponse(product)
}

func (service ServiceImpl) CreateMultiple(ctx context.Context, request []w_product.CreateRequest) []w_product.Response {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	epoch := helper.EpochTime()
	defer helper.CommitOrRollback(tx)

	var products []domain.Product
	for _, createProduct := range request {
		uuid, err := service.UuidRepo.CreteUui(ctx, tx)
		helper.PanicIfError(err)
		product := domain.Product{
			IdProduct: uuid.Uuid,
			Name:      createProduct.Name,
			Price:     createProduct.Price,
			Quantity:  createProduct.Quantity,
			ImageUrl:  createProduct.ImageUrl,
			CreatedAt: epoch,
		}

		products = append(products, product)
	}

	products = service.ProductRepo.SaveMultiple(ctx, tx, products)

	return h_product.ToProductResponses(products)
}

func (service ServiceImpl) CreateImg(ctx context.Context, request w_product.CreateRequest) w_product.Response {
	random := helper.RandStringRunes(10)
	s3Client, enpoint := helper.S3Config()

	object := s3.PutObjectInput{
		Bucket: aws.String("olshop-proto"),
		Key:    aws.String("/products/" + random + ".png"),
		Body:   strings.NewReader(string(request.ImageUrl)),
		ACL:    aws.String("public-read"),
	}

	_, err := s3Client.PutObject(&object)
	helper.PanicIfError(err)

	image := w_product.Response{
		ImageUrl: "https://" + *object.Bucket + "." + enpoint + *object.Key,
	}

	return image
}

func (service ServiceImpl) Update(ctx context.Context, request w_product.UpdateRequest) w_product.Response {
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(service.Validate, trans)

	err := service.Validate.Struct(request)
	if err != nil {
		errs := exception.TranslateError(err, trans)
		var x []string
		for _, v := range errs {
			x = append(x, v.Error())
		}
		panic(exception.NewValidationError(x))
	}

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepo.FindById(ctx, tx, request.IdProduct)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	product.IdProduct = request.IdProduct
	product.Name = request.Name
	product.Price = request.Price
	product.Quantity = request.Quantity
	product.ImageUrl = request.ImageUrl

	product = service.ProductRepo.Update(ctx, tx, product)
	return h_product.ToProductResponse(product)
}

func (service ServiceImpl) Delete(ctx context.Context, idProduct string) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepo.FindById(ctx, tx, idProduct)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.ProductRepo.Delete(ctx, tx, product)
}

func (service ServiceImpl) FindById(ctx context.Context, idProduct string) w_product.Response {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepo.FindById(ctx, tx, idProduct)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return h_product.ToProductResponse(product)
}

func (service ServiceImpl) Findall(ctx context.Context, pagination domain.Pagination) []w_product.Response {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	products := service.ProductRepo.Findall(ctx, tx, pagination)

	return h_product.ToProductResponses(products)
}

func (service ServiceImpl) DeleteTable(ctx context.Context) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	service.ProductRepo.DeleteTable(ctx, tx)
}
