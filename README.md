Проблемы и решения:

    1. Migrations
        Проблема: Не удалось сохранить файлы миграции  
        Решение: Прописать следующую команду для Смена владельца папки с root на arthas для доступа VS Code:
                sudo chown -R arthas:arthas ~/projects/todo-app - 
    2. Volumes
        Проблема: Не заполнялась папка out локально файлами из-за проблем с доступом
        Решение: Прописать следующие команды в Makefile для env-up:
                sleep 2 
                sudo chmod -R 777 ~/projects/todo-app/out/

--------------------------------------------------------------

RoadMap:

1. Настраиваем окружение

    1.1. Создаём ветку infra/env-setup

    1.2. Создаём файл docker-compose
        1.2.1 Описываем сервис БД todoapp-postgres
        1.2.2 Описываем сервис Миграций для БД todoapp-postgres-migrate
        1.2.3 Описываем сервис Прокси для связи между хостом и БД port-forwarder

    1.3. Создаём Makefile
    БД:
        1.3.1 Описываем команду env-up
        1.3.2 Описываем команду env-down
        1.3.3 Описываем команду env-cleanup

    Прокси - Socat:
        1.3.3 Описываем команду env-port-forward
        1.3.3 Описываем команду env-port-close

    Миграции для БД
        1.3.4 Описываем команду migrate-create
        1.3.5 Описываем команду migrate-up
        1.3.6 Описываем команду migrate-down
        1.3.6 Описываем команду migrate-action

    1.4. Создаём .env
        1.4.1 Создаём переменную POSTGRES_USER
        1.4.2 Создаём переменную POSTGRES_PASSWORD
        1.4.3 Создаём переменную POSTGRES_DB

    1.5. Создаём .gitignore
        1.5.1 Прописываем ".env" и "out"

    1.6. Создаём и заполяем файлы миграции
        1.6.1 Up 01
        1.6.2 Down 01



2. Общие файлы

    2.1 Создаём Логгер 
        2.1.1 Создаём директорию internal/core/logger.go
            2.1.2 Создаём структуру конфига для логгера, функцию-констурктор (lib: envconfig) и функцию "NewConfigMust"
            2.1.3 Создаём структуру логгера, функцию-констурктор и метод закрытия файла
            2.1.4 Создаём метод With, который будет отдавать наш собственный логгер
    
    2.2. Создаём MiddleWare
        2.2.1 Создаём директорию transport/http/middleware/middleware.go
            2.2.2 Создаём аллиас "Middleware" который является func(http.Handler) http.Handler
        2.2.3 Создаём директорию transport/http/middleware/common.go
            2.2.4 Создаём RequestID() Middleware
            2.2.5 Создаём Logger() Middleware
            2.2.5 Создаём Panic() Middleware

