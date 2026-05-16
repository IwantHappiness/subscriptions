# Subscriptions API

REST API для управления пользовательскими подписками.

Сервис позволяет создавать, получать, обновлять и удалять подписки, а также считать суммарную стоимость подписок за выбранный период с фильтрацией по пользователю и названию сервиса.

## Стек

- Go
- PostgreSQL
- Gorilla Mux
- pgx
- Docker Compose
- golang-migrate
- OpenAPI / Swagger UI

## Возможности

- CRUD для подписок.
- Хранение подписок в PostgreSQL.
- Даты подписок в формате `MM-YYYY`.
- Подсчет суммарной стоимости за период.
- Фильтрация расчета по `user_id` и `service_name`.
- Swagger UI с описанием API.
- HTTP-логи запросов: method, path, status, duration.

## Модель подписки

```json
{
  "id": 1,
  "service_name": "Spotify Premium",
  "price": 100,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "01-2024",
  "end_date": "03-2024"
}
```

`end_date` может быть `null`. Такая подписка считается активной без даты окончания.

## Подсчет стоимости

Ручка `GET /api/v1/subscriptions/total-price` считает сумму помесячно.

Пример: подписка стоит `100`, активна с `01-2024` по `03-2024`. Для периода `01-2024..03-2024` итог будет `300`.

В расчет попадают только месяцы, в которых период подписки пересекается с выбранным периодом.

## API

Базовый URL:

```text
http://localhost:8080
```

Основные ручки:

```text
GET    /api/v1/subscriptions?limit=100&offset=0
POST   /api/v1/subscriptions
GET    /api/v1/subscriptions/{id}
PUT    /api/v1/subscriptions/{id}
DELETE /api/v1/subscriptions/{id}
GET    /api/v1/subscriptions/total-price
```

Для списка подписок доступны параметры пагинации:

- `limit` - количество записей на странице, по умолчанию `100`, максимум `100`.
- `offset` - количество записей, которые нужно пропустить, по умолчанию `0`.

Пример запроса на подсчет стоимости:

```text
GET /api/v1/subscriptions/total-price?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Spotify%20Premium&from=01-2024&to=03-2024
```

## Swagger

После запуска сервиса документация доступна по адресам:

```text
http://localhost:8080/swagger/
http://localhost:8080/swagger/openapi.json
```

## Переменные окружения

Пример `.env`:

```env
HTTP_ADDR=:8080
LOG_LEVEL=info

POSTGRES_USER=postgres
POSTGRES_PASSWORD=sub-test-123
POSTGRES_DB=test-sub
```

`LOG_LEVEL` поддерживает значения `debug`, `info`, `warn`, `error`. Значение по умолчанию: `info`.

**Файл `.env` присутствует только в целях разработки. Не используйте его в продакшене.**

Приложение использует `DATABASE_DSN`. В Docker Compose он собирается автоматически:

```text
postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable
```

## Запуск через Docker Compose

Собрать окружение:

```bash
make env-build
```

Запустить приложение, PostgreSQL и миграции:

```bash
make env-up
```

Остановить окружение:

```bash
make env-down
```

## Миграции

Применить миграции:

```bash
make migrate-up
```

Откатить миграции:

```bash
make migrate-down
```

Создать новую миграцию:

```bash
make migrate-create seq=add_some_table
```

## Локальная разработка

Запустить тесты:

```bash
go test ./...
```

Запустить приложение локально:

```bash
go run ./cmd/api
```

При локальном запуске без Docker нужно передать `DATABASE_DSN`, указывающий на доступный PostgreSQL.

## Структура проекта

```text
cmd/api                         Точка входа приложения
internal/domain/subscription    Доменная модель
internal/usecase/subscription   Бизнес-логика и валидация
internal/repository/postgres    Работа с PostgreSQL
internal/transport/http         HTTP router, handlers, DTO, Swagger
migrations                      SQL-миграции
```
