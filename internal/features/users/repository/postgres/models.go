package users_postgres_repository

// UserModel — структура для маппинга строки таблицы `todoapp.users` в Go-тип.
// Порядок полей совпадает с порядком столбцов в SELECT-запросах репозитория.
type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}
