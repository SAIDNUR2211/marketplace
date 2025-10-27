# Marketplace API - Документация проекта

## 🎯 О проекте

**Marketplace API** - это высоконагруженный REST API сервис для управления онлайн-маркетплейсом. Сервис предоставляет полный функционал для управления пользователями, магазинами, товарами и заказами.

### Ключевые особенности:
- 🔐 **JWT аутентификация** с ролевой моделью доступа
- 🏪 **Мульти-тенантность** - поддержка множества магазинов
- 🛒 **Полный цикл заказов** - от создания до выполнения
- 💾 **PostgreSQL** - надежное хранение данных
- ⚡ **Redis** - кеширование для высокой производительности
- 🐳 **Docker** - контейнеризация и легкое развертывание
- 📚 **Swagger** - полная документация API

---

## 🏗️ Архитектура проекта

```
marketplace/
├── cmd/main.go                 # Точка входа
├── internal/
│   ├── configs/               # Конфигурация
│   ├── controller/            # HTTP handlers
│   ├── service/               # Бизнес-логика
│   ├── repository/            # Работа с БД
│   ├── models/                # Модели данных
│   └── contracts/             # Интерфейсы
├── migrations/                # Миграции БД
├── pkg/                       # Вспомогательные пакеты
└── utils/                     # Утилиты
```

### 🎪 Паттерны проектирования:
- **Clean Architecture** - разделение на слои
- **Repository Pattern** - абстракция доступа к данным
- **Dependency Injection** - инъекция зависимостей
- **Factory Pattern** - создание сервисов

---

## 📊 Модель данных

### Сущности системы:

#### 👥 Users (Пользователи)
```sql
id, username, email, password, role, phone, created_at, updated_at
```
**Роли:** `USER`, `SHOPKEPER`, `ADMIN`

#### 🏪 Shops (Магазины)
```sql
id, name, slug, owner_id, description, created_at, updated_at
```

#### 📦 Products (Товары)
```sql
id, name, slug, price, currency, quantity, shop_id, active, created_at
```

#### 🛒 Orders (Заказы)
```sql
id, user_id, shop_id, total, currency, status, note, created_at
```

#### 📋 Order Items (Позиции заказа)
```sql
id, order_id, product_id, quantity, unit_price, total_price
```

---

## 🔐 Система аутентификации и авторизации

### JWT Токены:
- **Access Token** - 720 минут (для запросов)
- **Refresh Token** - 35 дней (для обновления)

### Ролевая модель:
- **USER** - базовые права (просмотр, заказы)
- **SHOPKEPER** - управление магазинами и товарами
- **ADMIN** - полный доступ, управление пользователями

---

## 🚀 API Endpoints

### 🔑 Аутентификация
| Метод | Endpoint | Описание | Доступ |
|-------|----------|----------|---------|
| POST | `/auth/sign-up` | Регистрация | Public |
| POST | `/auth/sign-in` | Вход | Public |
| GET | `/auth/refresh` | Обновление токенов | Public |

### 👥 Пользователи
| Метод | Endpoint | Описание | Доступ |
|-------|----------|----------|---------|
| PUT | `/admin/users/{id}/role` | Изменить роль | ADMIN |

### 🏪 Магазины
| Метод | Endpoint | Описание | Доступ |
|-------|----------|----------|---------|
| POST | `/api/v1/shops` | Создать магазин | USER+ |
| GET | `/api/v1/shops/{id}` | Получить магазин | USER+ |
| GET | `/api/v1/shops` | Список магазинов | USER+ |
| PUT | `/api/v1/shops/{id}` | Обновить магазин | OWNER/ADMIN |
| DELETE | `/api/v1/shops/{id}` | Удалить магазин | OWNER/ADMIN |

### 📦 Товары
| Метод | Endpoint | Описание | Доступ |
|-------|----------|----------|---------|
| POST | `/api/v1/products` | Создать товар | SHOPKEPER+ |
| GET | `/api/v1/products/{id}` | Получить товар | Public |
| GET | `/api/v1/products` | Список товаров | Public |
| PUT | `/api/v1/products/{id}` | Обновить товар | OWNER/ADMIN |
| DELETE | `/api/v1/products/{id}` | Удалить товар | OWNER/ADMIN |

### 🛒 Заказы
| Метод | Endpoint | Описание | Доступ |
|-------|----------|----------|---------|
| POST | `/api/v1/orders` | Создать заказ | USER+ |
| GET | `/api/v1/orders/{id}` | Получить заказ | OWNER/ADMIN |

---

## 💾 База данных

### PostgreSQL Конфигурация:
```json
{
  "host": "localhost",
  "port": "5432", 
  "user": "postgres",
  "database": "shope_db"
}
```


## ⚡ Кеширование

### Redis Конфигурация:
```json
{
  "addr": "localhost:6379",
  "password": "",
  "db": 0
}
```

### Стратегия кеширования:
- **Products** - кеширование на 1 час
- **Авто-инвалидация** при обновлении/удалении
- **Cache-Aside** паттерн

---

## 🐳 Docker развертывание

### Запуск в Docker:
```bash
# Сборка и запуск
docker-compose up --build

# Остановка
docker-compose down

# Просмотр логов
docker-compose logs -f app
```

### Сервисы:
- **app** - основное приложение (порт 7577)
- **db** - PostgreSQL (порт 5433)
- **redis** - Redis (порт 6379)

---

## 🛠️ Установка и запуск

### Локальная разработка:

1. **Клонирование репозитория:**
```bash
git clone <repository-url>
cd marketplace
```

2. **Настройка окружения:**
```bash
# Создать .env файл
cp .env.example .env

# Установить переменные
DB_PASSWORD=your_password
JWT_SECRET=your_super_secret_jwt_key_that_is_long
```

3. **Запуск БД:**
```bash
docker-compose up db redis -d
```

4. **Применение миграций:**
```bash
make migrate-up
```

5. **Запуск приложения:**
```bash
go run cmd/main.go
```

---

## 📋 Тестирование API

### Swagger документация:
```
http://localhost:7577/swagger/index.html
```

### Примеры запросов:

#### 🔐 Регистрация пользователя
```bash
curl -X POST http://localhost:7577/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Иван Иванов",
    "username": "ivanov",
    "password": "password123", 
    "email": "ivan@example.com",
    "phone": "+79991234567"
  }'
```

#### 🛒 Создание заказа
```bash
curl -X POST http://localhost:7577/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "note": "Срочная доставка",
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'
```

---

## 🔧 Технические особенности

### 🚀 Производительность:
- **Connection Pooling** - пул соединений с БД
- **Transaction Management** - управление транзакциями
- **Optimistic Locking** - предотвращение race condition

### 🔒 Безопасность:
- **Password Hashing** - bcrypt для паролей
- **SQL Injection Protection** - параметризованные запросы
- **XSS Protection** - валидация входных данных

---

## 🐛 Отладка и логирование

### Уровни логирования:
- **INFO** - информационные сообщения
- **ERROR** - ошибки приложения
- **DEBUG** - отладочная информация

### Структура логов:
```json
{
  "level": "info",
  "time": "2024-01-15T10:30:00Z",
  "message": "Order created successfully",
  "order_id": 123,
  "user_id": 456
}
```

---

## 📈 Мониторинг и метрики

### Health Check:
```bash
curl http://localhost:7577/health
```

### Response:
```json
{
  "status": "healthy",
  "service": "MARKETPLACE", 
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## 🎯 Дальнейшее развитие

### Планируемые улучшения:
1. **Elasticsearch** - полнотекстовый поиск товаров
2. **Kafka** - асинхронная обработка заказов  
3. **Prometheus** - сбор метрик производительности
4. **Kubernetes** - оркестрация контейнеров
5. **GraphQL** - альтернативный API интерфейс

---
SELECT id, username, role FROM users WHERE username = 'admin_user';
UPDATE users SET role = 'ADMIN' WHERE id = 1;
