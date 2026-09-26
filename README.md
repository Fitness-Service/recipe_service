# Recipe Service

Микросервис управления рецептами фитнес-платформы. Отвечает за хранение, поиск и выдачу рецептов по заданным критериям (ингредиенты, цель питания). Работает по протоколу **gRPC**.

## Стек технологий

- **Язык:** Go
- **База данных:** PostgreSQL 18
- **Драйвер БД:** [sqlx](https://github.com/jmoiron/sqlx) + [lib/pq](https://github.com/lib/pq)
- **Транспорт:** [gRPC](https://grpc.io/) (HTTP/2 + Protocol Buffers)
- **Контракт:** [Protocol Buffers v3](https://protobuf.dev/)

## Структура проекта

```
recipe-service/
├── cmd/
│   └── main.go              # Точка входа: gRPC-сервер на :8082
├── internal/
│   ├── models/
│   │   └── recipe.go        # Доменная модель Recipe
│   ├── repository/
│   │   └── recipe.go        # Работа с PostgreSQL
│   └── handlers/
│       └── recipe.go        # HTTP-хендлеры (не используются в gRPC)
├── go.mod
├── go.sum
└── README.md
```

### Связанные репозитории

- [`shared`](https://github.com/Fitness-Service/shared) — общие модели, proto-контракты, события.
- [`user-service`](https://github.com/Fitness-Service/user-service) — сервис пользователей (HTTP/REST).
- [`api-gateway`](https://github.com/Fitness-Service/api-gateway) — HTTP→gRPC прокси для клиентов.

## Быстрый старт

### 1. Требования

- Go 1.22+
- PostgreSQL 16+
- `protoc` (для регенерации proto-кода)
- `grpcurl` (опционально, для тестирования)

### 2. Клонирование

```bash
git clone https://github.com/Fitness-Service/recipe-service.git
cd recipe-service
```

### 3. Настройка базы данных

Используется та же БД `fitness`, что и у `user-service`:

```sql
CREATE DATABASE fitness;
```

Таблица `recipes` создаётся автоматически при первом запуске сервиса (`CREATE TABLE IF NOT EXISTS` в `NewRecipeRepo`).

### 4. Переменные окружения

Создайте файл `.env` в корне проекта:

```env
DB_DSN=postgres://postgres:ВАШ_ПАРОЛЬ@localhost:5432/fitness?sslmode=disable
```

> Замените `ВАШ_ПАРОЛЬ` на пароль пользователя `postgres`.
> Если в пароле есть символы `@`, `:`, `/`, `#`, `?` — закодируйте их по URL.

### 5. Установка зависимостей

```bash
go mod tidy
```

### 6. Запуск

```bash
go run ./cmd
```

Ожидаемый вывод:

```
recipe-service gRPC started on :8082
```

Сервис будет доступен по адресу: **localhost:8082** (gRPC).

## 📡 gRPC API

Контракт описан в [`shared/proto/recipe/recipe.proto`](https://github.com/Fitness-Service/shared/blob/main/proto/recipe.proto).

### Методы сервиса `RecipeService`

| Метод | Вход | Выход | Описание |
|-------|------|-------|----------|
| `GetAllRecipes` | `Empty` | `RecipeList` | Получить все рецепты |
| `GetRecipe` | `GetRecipeRequest` | `Recipe` | Получить рецепт по ID |
| `SuggestByIngredients` | `IngredientRequest` | `RecipeList` | Рецепты, содержащие указанные ингредиенты |
| `SuggestByGoal` | `GoalRequest` | `RecipeList` | Рецепты под цель (`cut` / `bulk` / `maintain`) |
| `CreateRecipe` | `CreateRecipeRequest` | `Recipe` | Создать новый рецепт |

### Сообщения

```protobuf
message Recipe {
    string id = 1;
    string name = 2;
    repeated string ingredients = 3;
    int32 calories = 4;
    double protein = 5;
    double fats = 6;
    double carbs = 7;
    string goal = 8;
    repeated string tags = 9;
    int32 prep_time = 10;
    string created_at = 11;
}

message RecipeList {
    repeated Recipe recipes = 1;
}
```

## Примеры вызовов

### Через `grpcurl`

Установка:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Список сервисов:
```bash
grpcurl -plaintext localhost:8082 list
```

Получить все рецепты:
```bash
grpcurl -plaintext \
  -import-path ../shared/proto \
  -proto recipe.proto \
  localhost:8082 recipe.RecipeService/GetAllRecipes
```

Создать рецепт:
```bash
grpcurl -plaintext \
  -import-path ../shared/proto \
  -proto recipe.proto \
  -d '{
    "name": "Овсянка с бананом",
    "ingredients": ["овсянка", "банан", "молоко"],
    "calories": 350,
    "protein": 12.5,
    "fats": 6.0,
    "carbs": 60.0,
    "goal": "maintain",
    "tags": ["завтрак", "быстро"],
    "prep_time": 10
  }' \
  localhost:8082 recipe.RecipeService/CreateRecipe
```

Найти рецепты по ингредиентам:
```bash
grpcurl -plaintext \
  -import-path ../shared/proto \
  -proto recipe.proto \
  -d '{"ingredients": ["курица", "рис"]}' \
  localhost:8082 recipe.RecipeService/SuggestByIngredients
```

### Через Go-клиент

```go
conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := pb.NewRecipeServiceClient(conn)

resp, err := client.GetAllRecipes(context.Background(), &pb.Empty{})
if err != nil {
    log.Fatal(err)
}

for _, r := range resp.Recipes {
    fmt.Println(r.Name, r.Calories, "kcal")
}
```

### Через Postman

1. Создайте новый **gRPC Request**.
2. URL: `localhost:8082`.
3. Импортируйте `recipe.proto`.
4. Выберите метод `RecipeService/GetAllRecipes`.
5. Нажмите **Invoke**.

## Схема базы данных

Таблица `recipes`:

| Поле | Тип | Описание |
|------|-----|----------|
| id | UUID | Первичный ключ |
| name | TEXT | Название рецепта |
| ingredients | TEXT[] | Массив ингредиентов |
| calories | INT | Калорийность |
| protein | FLOAT | Белки (г) |
| fats | FLOAT | Жиры (г) |
| carbs | FLOAT | Углеводы (г) |
| goal | TEXT | Цель (`cut` / `bulk` / `maintain`) |
| tags | TEXT[] | Теги (`завтрак`, `быстро`, `веган`) |
| prep_time | INT | Время приготовления (мин) |
| created_at | TIMESTAMP | Дата создания |

## Генерация proto-кода

Контракт лежит в `shared/proto/recipe/recipe.proto`. Для регенерации:

```bash
cd ../shared
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/recipe/recipe.proto
```

> Флаг `paths=source_relative` обязателен — иначе `protoc` создаст лишние папки `github.com/...`.

Для установки плагинов:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Безопасность

- Используется **gRPC без TLS** (`WithInsecure`) — для локальной разработки. В продакшене нужно добавить TLS-сертификаты.
- Пароль от БД **не должен** попадать в Git — используйте `.env` и `.gitignore`.

## Roadmap

- [ ] TLS для gRPC
- [ ] Интерцепторы для логирования и аутентификации
- [ ] Пагинация в `GetAllRecipes`
- [ ] Кэширование популярных рецептов (Redis)
- [ ] Docker-образ и docker-compose
- [ ] Покрытие тестами
