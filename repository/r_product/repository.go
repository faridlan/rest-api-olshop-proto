package r_product

import (
	"context"
	"database/sql"

	"github.com/faridlan/rest-api-olshop-proto/model/domain"
)

type Repository interface {
	Save(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product
	SaveMultiple(ctx context.Context, tx *sql.Tx, products []domain.Product) []domain.Product
	Update(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product
	Delete(ctx context.Context, tx *sql.Tx, product domain.Product)
	FindById(ctx context.Context, tx *sql.Tx, idProduct string) (domain.Product, error)
	Findall(ctx context.Context, tx *sql.Tx, pagination domain.Pagination) []domain.Product
	DeleteTable(ctx context.Context, tx *sql.Tx)
}
