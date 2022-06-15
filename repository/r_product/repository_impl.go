package r_product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/faridlan/rest-api-olshop-proto/helper"
	"github.com/faridlan/rest-api-olshop-proto/model/domain"
)

type RepositoryImpl struct {
}

func NewProductRepository() Repository {
	return RepositoryImpl{}
}

func (repository RepositoryImpl) Save(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product {
	SQL := "insert into products (id_product, name, price, quantity, image_url, created_at) values (?,?,?,?,?,?)"
	_, err := tx.ExecContext(ctx, SQL, product.IdProduct, product.Name, product.Price, product.Quantity, product.ImageUrl, product.CreatedAt)
	helper.PanicIfError(err)

	return product
}

func (repository RepositoryImpl) SaveMultiple(ctx context.Context, tx *sql.Tx, products []domain.Product) []domain.Product {
	SQL := "insert into products (id_product, name, price, quantity, image_url, created_at) values"
	var vals []interface{}

	for _, product := range products {
		SQL += "(?,?,?,?,?,?),"
		vals = append(vals, product.IdProduct, product.Name, product.Price, product.Quantity, product.ImageUrl, product.CreatedAt)
	}

	SQL = SQL[0 : len(SQL)-1]

	stmt, err := tx.PrepareContext(ctx, SQL)
	helper.PanicIfError(err)
	defer stmt.Close()

	_, err = stmt.Exec(vals...)
	helper.PanicIfError(err)
	return products
}

func (repository RepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product domain.Product) domain.Product {
	SQL := "update products set name = ?, price = ?, quantity = ?, image_url = ? where id_product = ?"
	_, err := tx.ExecContext(ctx, SQL, product.Name, product.Price, product.Quantity, product.ImageUrl, product.IdProduct)
	helper.PanicIfError(err)

	return product
}

func (repository RepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, product domain.Product) {
	SQL := "delete from products where id_product = ?"
	_, err := tx.ExecContext(ctx, SQL, product.IdProduct)
	helper.PanicIfError(err)
}

func (repository RepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, idProduct string) (domain.Product, error) {
	SQL := "select id_product, name, price, quantity, image_url, created_at from products where id_product = ?"
	rows, err := tx.QueryContext(ctx, SQL, idProduct)
	helper.PanicIfError(err)

	defer rows.Close()

	product := domain.Product{}
	if rows.Next() {
		err := rows.Scan(&product.IdProduct, &product.Name, &product.Price, &product.Quantity, &product.ImageUrl, &product.CreatedAt)
		helper.PanicIfError(err)
		return product, nil
	} else {
		return product, errors.New("product not found")
	}
}

func (repository RepositoryImpl) Findall(ctx context.Context, tx *sql.Tx, pagination domain.Pagination) []domain.Product {
	SQL := fmt.Sprintf(`
	select id_product, name, price, quantity, image_url, created_at from products 
	order by created_at desc limit %d,%d`, pagination.Page, pagination.Limit)
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)

	defer rows.Close()

	products := []domain.Product{}
	for rows.Next() {
		product := domain.Product{}
		err := rows.Scan(&product.IdProduct, &product.Name, &product.Price, &product.Quantity, &product.ImageUrl, &product.CreatedAt)
		helper.PanicIfError(err)
		products = append(products, product)
	}

	return products
}

func (repository RepositoryImpl) DeleteTable(ctx context.Context, tx *sql.Tx) {
	SQL := "delete from products where id_product not in ( select id_product where id_product = '08eea32bec0311ec9d4d0242ac130002' )"
	_, err := tx.ExecContext(ctx, SQL)
	helper.PanicIfError(err)
}
