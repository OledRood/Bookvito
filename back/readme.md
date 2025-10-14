# Bookvito Backend

Сервис обмена книгами на Go с чистой архитектурой.

## Структура проекта

```
back/
├── cmd/
│   └── api/
│       └── main.go              # Точка входа в приложение
├── internal/
│   ├── domain/                  # Слой бизнес-логики (entities & interfaces)
│   │   ├── entities.go          # Модели: User, Book, Exchange
│   │   └── repository.go        # Интерфейсы репозиториев
│   ├── usecase/                 # Бизнес-логика приложения
│   │   ├── user_usecase.go      # Use cases для пользователей
│   │   ├── book_usecase.go      # Use cases для книг
│   │   └── exchange_usecase.go  # Use cases для обменов
│   ├── repository/              # Реализация репозиториев
│   │   └── postgres/
│   │       ├── user_repository.go
│   │       ├── book_repository.go
│   │       └── exchange_repository.go
│   └── delivery/                # HTTP обработчики
│       └── http/
│           ├── router.go
│           ├── user_handler.go
│           ├── book_handler.go
│           └── exchange_handler.go
├── pkg/                         # Общие пакеты
│   └── database/
│       └── postgres.go          # Подключение к PostgreSQL
├── config/                      # Конфигурация
│   └── config.go
├── go.mod
├── .env.example
└── .gitignore
```

## Установка

1. Клонируйте репозиторий
2. Скопируйте `.env.example` в `.env` и настройте переменные окружения
3. Установите зависимости:

```bash
cd back
go mod download
```

## API Endpoints

### Users
- `POST /api/v1/users/register` - Регистрация пользователя
- `POST /api/v1/users/login` - Вход пользователя
- `GET /api/v1/users/:id` - Получить пользователя
- `PUT /api/v1/users/:id` - Обновить пользователя
- `DELETE /api/v1/users/:id` - Удалить пользователя
- `GET /api/v1/users` - Список пользователей

### Books
- `POST /api/v1/books` - Создать книгу
- `GET /api/v1/books/:id` - Получить книгу
- `GET /api/v1/books` - Список книг
- `GET /api/v1/books/search?q=query` - Поиск книг
- `GET /api/v1/books/available` - Доступные книги
- `PUT /api/v1/books/:id` - Обновить книгу
- `DELETE /api/v1/books/:id` - Удалить книгу
- `GET /api/v1/books/owner/:owner_id` - Книги владельца

### Exchanges
- `POST /api/v1/exchanges` - Создать запрос на обмен
- `GET /api/v1/exchanges/:id` - Получить обмен
- `GET /api/v1/exchanges` - Список обменов
- `GET /api/v1/exchanges/requester/:requester_id` - Обмены запросившего
- `GET /api/v1/exchanges/owner/:owner_id` - Обмены владельца
- `PUT /api/v1/exchanges/:id/accept` - Принять обмен
- `PUT /api/v1/exchanges/:id/reject` - Отклонить обмен
- `PUT /api/v1/exchanges/:id/complete` - Завершить обмен

## Запуск

### Локальный запуск

```bash
# Запуск сервера
go run cmd/api/main.go
```

### Запуск с Docker

```bash
# Запуск всех сервисов (API + PostgreSQL)
docker-compose up -d

# Остановка сервисов
docker-compose down

# Просмотр логов
docker-compose logs -f api

# Пересборка образа
docker-compose up -d --build
```

### Запуск только Docker образа

```bash
# Сборка образа
docker build -t bookvito-api .

# Запуск контейнера
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=bookvito \
  bookvito-api
```

Сервер будет доступен по адресу `http://localhost:8080`

## Используемые библиотеки

- **Gin** - HTTP фреймворк
- **GORM** - ORM для работы с БД
- **PostgreSQL Driver** - Драйвер PostgreSQL
- **bcrypt** - Хэширование паролей
- **godotenv** - Загрузка .env файлов

## Принципы чистой архитектуры

Проект следует принципам чистой архитектуры:

1. **Domain Layer** - бизнес-логика и интерфейсы (не зависит от внешних библиотек)
2. **Use Case Layer** - прикладная бизнес-логика
3. **Repository Layer** - доступ к данным (реализация интерфейсов)
4. **Delivery Layer** - HTTP handlers (точка входа для запросов)

Зависимости направлены внутрь: Delivery → UseCase → Domain ← Repository
