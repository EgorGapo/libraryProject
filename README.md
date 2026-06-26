# Library

Сервис каталога книг и авторов на Go. Предоставляет gRPC API и REST-интерфейс поверх него через gRPC-Gateway, хранит данные в PostgreSQL и обеспечивает консистентность операций на уровне транзакций.

## Возможности

- gRPC API и автоматически сгенерированный REST-слой (gRPC-Gateway) с единой схемой из `.proto`
- Управление авторами и книгами со связью «многие ко многим»
- Потоковая выдача книг автора (server-side streaming)
- Хранение в PostgreSQL: миграции, индексы, генерация UUID на уровне БД, триггеры `updated_at`
- Атомарность операций через транзакции (книга и её связи с авторами пишутся одним коммитом)
- Валидация запросов на уровне API (`protoc-gen-validate`)
- Конфигурация через переменные окружения
- Структурированное логирование (`zap`)
- Корректное завершение по `SIGINT` / `SIGTERM`
- Юнит-тесты с моками и интеграционные тесты

## Стек

| Категория | Технологии |
|-----------|-----------|
| Язык | Go |
| API | gRPC, gRPC-Gateway (REST-to-gRPC) |
| База данных | PostgreSQL (драйвер `pgx`) |
| Миграции | `goose` |
| Валидация | `protoc-gen-validate` |
| Логирование | `zap` |
| Контейнеризация | Docker, docker-compose |

## Архитектура

Проект построен по принципам чистой архитектуры (на основе [go-clean-template](https://github.com/evrone/go-clean-template)): транспортный слой, доменная логика (use cases) и репозитории разделены и зависят через интерфейсы. Это позволяет подменять реализацию хранилища и покрывать логику тестами с моками.

```
cmd/library        — точка входа
internal/app       — сборка зависимостей, запуск gRPC и HTTP-gateway, graceful shutdown
internal/controller— gRPC-хендлеры (адаптация proto ↔ домен)
internal/usecase   — бизнес-логика и валидация
internal/usecase/repository — интерфейс хранилища + реализации (in-memory, postgres)
internal/entity    — доменные сущности и ошибки
db/migrations      — SQL-миграции (goose)
api/library        — proto-схема API
config             — конфигурация из env
```

## API

REST-пути генерируются из gRPC-описания. Идентификаторы — в формате UUID.

| Метод | REST | gRPC | Описание |
|-------|------|------|----------|
| POST | `/v1/library/book` | `AddBook` | Добавить книгу |
| PUT | `/v1/library/book` | `UpdateBook` | Обновить книгу |
| GET | `/v1/library/book/{id}` | `GetBookInfo` | Получить книгу по ID |
| POST | `/v1/library/author` | `RegisterAuthor` | Зарегистрировать автора |
| PUT | `/v1/library/author` | `ChangeAuthorInfo` | Изменить данные автора |
| GET | `/v1/library/author/{id}` | `GetAuthorInfo` | Получить автора по ID |
| GET | `/v1/library/author_books/{author_id}` | `GetAuthorBooks` | Книги автора (стрим) |

### Валидация

- ID книги и автора — корректный UUID
- Имя автора соответствует `^[A-Za-z0-9]+( [A-Za-z0-9]+)*$`, длина от 1 до 512 символов

После генерации `swagger.json` REST-схему можно открыть в [Swagger Editor](https://editor.swagger.io/).

## Модель данных

Три таблицы: `author`, `book` и связующая `author_book`.

- `author` и `book` — UUID в качестве первичного ключа (`DEFAULT uuid_generate_v4()`), поля `created_at` / `updated_at` с триггером автообновления времени изменения
- `author_book` — композитный первичный ключ `(author_id, book_id)`, внешние ключи с `ON DELETE CASCADE`, отдельный индекс на `book_id`
- Индексы на имя автора и имя книги

При создании книги вставка записи в `book` и связей в `author_book` выполняется в одной транзакции, поэтому несогласованных данных не остаётся даже при ошибке.

## Конфигурация

Сервис настраивается через переменные окружения:

| Переменная | Назначение |
|------------|-----------|
| `GRPC_PORT` | Порт gRPC-сервера |
| `GRPC_GATEWAY_PORT` | Порт REST-сервера (gRPC-Gateway) |
| `POSTGRES_HOST` | Хост PostgreSQL |
| `POSTGRES_PORT` | Порт PostgreSQL |
| `POSTGRES_DB` | Имя базы данных |
| `POSTGRES_USER` | Пользователь |
| `POSTGRES_PASSWORD` | Пароль |
| `POSTGRES_MAX_CONN` | Максимальный размер пула соединений |

Строка подключения формируется в виде:

```
postgres://user:password@host:port/dbname?sslmode=disable&pool_max_conns=10
```

## Запуск

### Требования

- Go
- Docker и docker-compose

### Локально

Поднять базу данных:

```bash
docker-compose up -d
```

Миграции применяются автоматически при старте сервиса (загружаются через `go:embed`). После запуска в логах видно накат:

```
OK   001_create_author_table.sql
OK   002_create_author_name_index.sql
OK   003_create_book_table.sql
OK   004_create_book_name_index.sql
OK   005_create_author_book_table.sql
OK   006_create_author_book_book_id_index.sql
goose: successfully migrated database to version: 6
```

После этого сервис принимает запросы на `GRPC_PORT` (gRPC) и `GRPC_GATEWAY_PORT` (REST).

## Разработка

Для локальной работы используется `Makefile`:

```bash
make all      # линтер + тесты
make test     # только тесты
make lint     # только линтер
make generate # генерация кода из .proto
make build    # сборка
```

## Тестирование

Доменная логика покрыта юнит-тестами с использованием сгенерированных моков; интеграционные тесты проверяют сервис целиком, включая работу с базой данных. Тесты лежат в `integration-test/`.


