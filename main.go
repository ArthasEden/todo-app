package main

import (
	"context"
	"fmt"
	simpleconnection "todo-app/repository/simpleConnection"
	simplesql "todo-app/repository/simpleSql"
)

func main() {
	ctx := context.Background()
	conn, err := simpleconnection.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}

	if err := simplesql.CreateTable(ctx, conn); err != nil {
		panic(err)
	}

	fmt.Println("Подключение к базе данных прошло успешно!")
}
