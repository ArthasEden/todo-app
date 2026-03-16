package simple_connection

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, "postgres://sergey:1234@192.168.23.131:5432/sergey")
}
