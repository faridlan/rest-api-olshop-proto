package r_uuid

import (
	"context"
	"database/sql"
	"errors"

	"github.com/faridlan/rest-api-olshop-proto/helper"
	"github.com/faridlan/rest-api-olshop-proto/model/domain"
)

type RepositoryImpl struct {
}

func NewUuidRepository() Repository {
	return RepositoryImpl{}
}

func (repository RepositoryImpl) CreteUui(ctx context.Context, tx *sql.Tx) (domain.Uuid, error) {
	SQL := "select REPLACE(UUID(),'-','')"
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)

	defer rows.Close()
	uuid := domain.Uuid{}
	if rows.Next() {
		err := rows.Scan(&uuid.Uuid)
		helper.PanicIfError(err)
		return uuid, nil
	} else {
		return uuid, errors.New("failed create uuid")
	}
}
