01.04 
commit "infra setup"

Текущая архитектура:

Хост (WSL)
    │
    ├── todoapp-postgres-port-forwarder (порт 5432 на хосте)
    │       │
    │       └── перенаправляет на → todoapp-env-postgres:5432
    │
    └── todoapp-env-postgres (внутренний порт 5432, НЕ доступен с хоста напрямую)
--------------------------------------------------------------

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

Ход действий:

1. Создаём ветку infra/env-setup

2. Создаём файл docker-compose
    2.1 Описываем сервис БД todoapp-postgres
    2.2 Описываем сервис Миграций для БД todoapp-postgres-migrate
    2.3 Описываем сервис Прокси для связи между хостом и БД port-forwarder

3. Создаём Makefile
    БД:
    3.1 Описываем команду env-up
    3.2 Описываем команду env-down
    3.3 Описываем команду env-cleanup

    Прокси - Socat:
    3.3 Описываем команду env-port-forward
    3.3 Описываем команду env-port-close

    Миграции для БД
    3.4 Описываем команду migrate-create
    3.5 Описываем команду migrate-up
    3.6 Описываем команду migrate-down
    3.6 Описываем команду migrate-action

4. Создаём .env
    4.1 Создаём переменную POSTGRES_USER
    4.2 Создаём переменную POSTGRES_PASSWORD
    4.3 Создаём переменную POSTGRES_DB

5. Создаём .gitignore
    5.1 Прописываем ".env" и "out"

6. Создаём и заполяем файлы миграции
    6.1 Up 01
    6.2 Down 01
