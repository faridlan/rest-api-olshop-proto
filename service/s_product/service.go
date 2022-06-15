package s_product

import (
	"context"

	"github.com/faridlan/rest-api-olshop-proto/model/domain"
	"github.com/faridlan/rest-api-olshop-proto/model/web/w_product"
)

type Service interface {
	Create(ctx context.Context, request w_product.CreateRequest) w_product.Response
	CreateMultiple(ctx context.Context, request []w_product.CreateRequest) []w_product.Response
	CreateImg(ctx context.Context, request w_product.CreateRequest) w_product.Response
	Update(ctx context.Context, request w_product.UpdateRequest) w_product.Response
	Delete(ctx context.Context, idProduct string)
	FindById(ctx context.Context, idProduct string) w_product.Response
	Findall(ctx context.Context, pagination domain.Pagination) []w_product.Response
	DeleteTable(ctx context.Context)
}
