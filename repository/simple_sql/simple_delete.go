package simple_sql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func DeleteRow(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	DELETE FROM tasks
	WHERE id IN(3,5);
	`

	_, err := conn.Exec(ctx, sqlQuery)

	return err
}
