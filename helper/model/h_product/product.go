package h_product

import (
	"github.com/faridlan/rest-api-olshop-proto/model/domain"
	"github.com/faridlan/rest-api-olshop-proto/model/web/w_product"
)

func ToProductResponse(domain domain.Product) w_product.Response {
	return w_product.Response{
		IdProduct: domain.IdProduct,
		Name:      domain.Name,
		Price:     domain.Price,
		Quantity:  domain.Quantity,
		ImageUrl:  domain.ImageUrl,
		CreatedAt: domain.CreatedAt,
	}
}

func ToProductResponses(domain []domain.Product) []w_product.Response {
	productResponses := []w_product.Response{}
	for _, product := range domain {
		productResponses = append(productResponses, ToProductResponse(product))
	}

	return productResponses
}
