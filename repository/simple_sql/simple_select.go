package simple_sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func SelectRows(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	SELECT * FROM public.tasks
	ORDER BY id ASC 
	LIMIT 100
	`

	rows, err := conn.Query(ctx, sqlQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var text string
		var completed bool
		var created_at time.Time
		var completed_at *time.Time

		err := rows.Scan(
			&id,
			&title,
			&text,
			&completed,
			&created_at,
			&completed_at,
		)
		if err != nil {
			return err
		}

		fmt.Println(id, title, text, completed, created_at, completed_at)
	}

	return nil
}
