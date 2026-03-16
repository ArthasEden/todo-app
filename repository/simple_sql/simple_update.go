package simple_sql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func UpdateRow(ctx context.Context, conn *pgx.Conn) error {
	sqlUpdate := `
	UPDATE tasks
	SET text = ';)'
	WHERE completed = FALSE;
	`
	_, err := conn.Exec(ctx, sqlUpdate)

	return err
}
