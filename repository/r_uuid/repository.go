package r_uuid

import (
	"context"
	"database/sql"

	"github.com/faridlan/rest-api-olshop-proto/model/domain"
)

type Repository interface {
	CreteUui(ctx context.Context, tx *sql.Tx) (domain.Uuid, error)
}
