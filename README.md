# Track Selection — Сервис выбора образовательного трека

> Интеллектуальная платформа для подбора образовательного трека в университете.  
> Заполни профиль, укажи цели и навыки — система подберёт подходящие треки на основе метода PROMETHEE.

---

## Технологический стек

### Бэкенд (Go)

- **Язык**: Go (Golang)
- **API стиль**: RESTful
- **Архитектура**: Clean Architecture + DDD + Event-Driven — слои `domain` (агрегаты, value objects, domain events, интерфейсы репозиториев), `application` (use cases), `infrastructure` (HTTP, PostgreSQL, JWT, event bus), `bootstrap` (DI)
- **HTTP роутер**: `gorilla/mux`
- **JWT**: `golang-jwt/jwt` — Bearer-токены, время жизни 24 часа
- **Event Bus**: внутренняя шина событий (in-memory)
- **Миграции**: `pressly/goose` — версионированные SQL-миграции
- **Graceful Shutdown**: корректное завершение по `SIGTERM` / `Ctrl+C`

### Алгоритм рекомендаций

Рекомендации треков строятся методом **PROMETHEE** (Preference Ranking Organization METHod for Enrichment Evaluations) — многокритериальный метод принятия решений.

Критерии сравнения:
- Совпадение профессиональных целей студента и трека
- Академические оценки vs. минимальные требования трека
- Навыки студента (tech / math / soft) vs. желаемые навыки трека
- Стиль обучения
- Наличие сертификатов
- Перспективы трудоустройства и отзывы выпускников

### База данных

| СУБД       | Назначение                                                                   |
|------------|------------------------------------------------------------------------------|
| PostgreSQL | Пользователи, студенты, администраторы, предпочтения, треки, выборы треков  |

Треки хранятся с полным учебным планом в JSONB (`curriculum`) — поддерживаются линейные треки (`single`) и ветвящиеся (`branching`) с рекурсивными подветками.

### Фронтенд

- **React 18** + **Vite**
- Vanilla CSS (без UI-библиотек, без Tailwind)
- Адаптивная вёрстка (брейкпоинты: 1023px / 767px / 479px)
- Поддержка светлой / тёмной темы (системная + ручная)
- Токен и данные пользователя хранятся в `localStorage`

### Инфраструктура

- **Docker** — контейнеризация всех сервисов
- **Docker Compose** — оркестрация: Go-бэкенд, PostgreSQL, Nginx
- **Nginx** — reverse-proxy, раздача статики (сборка фронтенда из `/static`)

---

## Структура проекта

```
track-selection/
├── cmd/
│   ├── httpBack/main.go          # Точка входа HTTP-сервера (порт 8080)
│   └── migrate/main.go           # Точка входа для миграций (Goose)
├── internal/
│   ├── domain/                   # Сущности, value objects, интерфейсы репозиториев
│   │   ├── auth/                 # Пользователь, JWT-интерфейс
│   │   ├── student/              # Студент, предпочтения, профиль, checker
│   │   └── track/                # Трек, учебный план (curriculum), требования
│   ├── application/              # Use Cases (бизнес-сценарии)
│   │   ├── auth/                 # Регистрация, вход
│   │   ├── student/              # Предпочтения, рекомендации, выбор трека
│   │   └── track/                # CRUD треков
│   ├── infrastructure/           # HTTP-хендлеры, PostgreSQL-репозитории, JWT-сервис
│   │   ├── http/handlers/        # REST-хендлеры
│   │   ├── http/middleware/      # Auth middleware, context middleware
│   │   ├── persistence/postgres/ # Реализации репозиториев
│   │   ├── jwt/                  # JWT-реализация
│   │   └── eventbus/             # In-memory шина событий + подписчики
│   ├── config/                   # Конфигурация из env
│   └── bootstrap/app.go          # DI: сборка всего приложения
├── migrations/                   # SQL и Go миграции (Goose)
├── frontend/                     # React-приложение
│   └── src/
│       ├── App.jsx               # Корневой компонент: auth-состояние, модалки, тема
│       ├── pages/Home.jsx        # Каталог треков: вкладки, поиск, фильтры, карточки
│       ├── components/           # Header, Footer, Modal, TrackCard, ProfileModal и др.
│       ├── services/auth.js      # API-клиент + auth-хелперы
│       └── index.css             # Все стили — единый файл, BEM-подобные классы
├── nginx/nginx.conf
├── docker-compose.yml
├── Makefile
└── api.yaml                      # OpenAPI 3.0 спецификация
```

---

## Обзор API

Все маршруты — под префиксом `/api`.

### Аутентификация

| Метод | Путь        | Auth | Описание                            |
|-------|-------------|------|-------------------------------------|
| POST  | `/register` | —    | Регистрация (`student` или `admin`) |
| POST  | `/login`    | —    | Вход → возвращает `{ token }`       |

### Пользователь и профиль

| Метод | Путь                     | Auth   | Описание                       |
|-------|--------------------------|--------|--------------------------------|
| GET   | `/me`                    | Bearer | Данные текущего пользователя   |
| GET   | `/me/profile-completion` | Bearer | Статус заполненности профиля   |
| GET   | `/me/info`               | Bearer | Предпочтения студента          |
| POST  | `/me/edit-info`          | Bearer | Обновить предпочтения студента |

### Треки

| Метод  | Путь                  | Auth         | Описание        |
|--------|-----------------------|--------------|-----------------|
| GET    | `/all-tracks`         | —            | Все треки       |
| POST   | `/new-track`          | Bearer admin | Создать трек    |
| PUT    | `/edit-track/{id}`    | Bearer admin | Обновить трек   |
| DELETE | `/delete-track/{id}`  | Bearer admin | Удалить трек    |

### Рекомендации и выбор трека

| Метод  | Путь                           | Auth   | Описание                              |
|--------|--------------------------------|--------|---------------------------------------|
| GET    | `/student/recommendations`     | Bearer | Рекомендации по методу PROMETHEE      |
| POST   | `/student/select-track`        | Bearer | Выбрать трек                          |
| GET    | `/student/selected-tracks`     | Bearer | Список выбранных треков               |
| DELETE | `/student/unselect-track/{id}` | Bearer | Отменить выбор трека                  |

> Формат ошибок: `{ "error": { "code": "...", "message": "..." } }`

---

## Структура трека

Каждый трек содержит:

| Поле                   | Тип       | Описание                                       |
|------------------------|-----------|------------------------------------------------|
| `name`                 | string    | Название трека                                 |
| `description`          | string    | Описание                                       |
| `curriculum`           | JSONB     | Учебный план по годам (single / branching)     |
| `requirements`         | JSONB     | Минимальные оценки по предметам                |
| `teachers`             | JSONB     | Список преподавателей                          |
| `difficulty`           | int (1–5) | Сложность трека                                |
| `employment_prospects` | int (1–10)| Перспективы трудоустройства                    |
| `alumni_reviews`       | int (1–10)| Оценки выпускников                             |
| `has_certificates`     | 0 / 1     | Наличие сертификатов                           |
| `learning_style`       | 1 / 2 / 3 | Стиль обучения (теория / практика / смешанный) |
| `desired_tech_skills`  | int (1–10)| Нужный уровень технических навыков             |
| `desired_math_skills`  | int (1–10)| Нужный уровень математических навыков          |
| `desired_soft_skills`  | int (1–10)| Нужный уровень soft skills                     |
| `professional_goals`   | []int     | Профессиональные цели, которым подходит трек   |

---

## Профиль студента

Перед получением рекомендаций студент заполняет профиль:

**Оценки** (от 2 до 5):  
Информатика, Программирование, Иностранный язык, Физика, АИГ, Математический анализ, Алгоритмы и структуры данных, Базы данных, Дискретная математика

**Навыки** (от 0 до 10):  
Базы данных, Системная архитектура, Алгоритмическое программирование, Публичные выступления, Тестирование, Аналитика, Machine Learning, Знание ОС, Исследовательские проекты

**Дополнительно:**
- Профессиональные цели (массив int)
- Стиль обучения: `1` — теория, `2` — практика, `3` — смешанный
- Желание получить сертификат: `0` / `1`

---

## Переменные окружения

| Переменная              | Описание                                           |
|-------------------------|----------------------------------------------------|
| `STANDART_PG_USER`      | Суперпользователь PostgreSQL                       |
| `STANDART_PG_PASSWORD`  | Пароль суперпользователя                           |
| `STANDART_PG_DB_NAME`   | Имя базы данных                                    |
| `PG_USERNAME_FOR_APP`   | Пользователь БД для приложения                     |
| `PG_USERPASS_FOR_APP`   | Пароль пользователя БД приложения                  |
| `PG_HOST`               | Хост PostgreSQL (в Docker: `postgres`)             |
| `PG_PORT`               | Порт PostgreSQL (обычно `5432`)                    |
| `PG_SSLMODE`            | Режим SSL (`disable` для локальной разработки)     |
| `APP_PORT`              | Порт Go-приложения (обычно `8080`)                 |
| `NGINX_PORT`            | Внешний порт Nginx                                 |
| `JWT_SECRET_KEY`        | Секрет для подписи JWT-токенов                     |
| `LOGS_LEVEL_APP`        | Уровень логирования приложения (`INFO`, `DEBUG`, `WARN`, `ERROR`) |
| `LOGS_LEVEL_MIGRATE`    | Уровень логирования миграций                       |

---

## Запуск

### Через Docker Compose

```bash
git clone <repo-url>
cd track-selection

# Создать .env на основе примера
cp .env.example .env
# Отредактировать .env при необходимости

# Запустить все сервисы (postgres → migrate → app → nginx)
docker-compose up --build
```

Приложение будет доступно по адресу `http://localhost:<NGINX_PORT>`.

### Разработка фронтенда

```bash
cd frontend
npm install
npm run dev      # Dev-сервер с проксированием /api → :8080
```

```bash
cd frontend
npm run build    # Сборка → /static (раздаётся Nginx)
```

---

## Миграции

```bash
# Запустить только postgres
docker compose up postgres -d

# Применить все миграции
docker compose run --rm migrate up

# Откатить последнюю миграцию
docker compose run --rm migrate down

# Сбросить все миграции
docker compose run --rm migrate reset

# Текущая версия миграции
docker compose run --rm migrate version
```

---

## E2E тесты

Требуется запущенный бэкенд на порту `:8080`.

```bash
make e2e-auth       # Тесты аутентификации
make e2e-student    # Тесты профиля студента
make e2e-tracks     # Тесты CRUD треков (admin)
make e2e-recomm     # Тесты рекомендаций и выбора трека
```

---

## Архитектурные решения

### DDD (Domain-Driven Design)

Проект следует тактическим паттернам DDD. Домен разбит на несколько **Bounded Context'ов**, каждый живёт в своём пакете:

| Контекст | Пакет | Агрегаты / Сущности |
|----------|-------|---------------------|
| Аутентификация | `domain/auth` | `AuthUser` — учётная запись с email/password-хешем |
| Студент | `domain/student` | `Student`, `Preferences`, `ProfileCompletion`, `TrackSelection` |
| Администратор | `domain/admin` | `Admin` |
| Трек | `domain/track` | `Track` с вложенным `Curriculum` |
| Рекомендации | `domain/selection` | Доменная логика PROMETHEE (чистые функции) |

**Value Objects** — каждый идентификатор (`StudentID`, `AdminID`) и общий тип (`Email`) реализованы как value object: инкапсулируют валидацию, неизменяемы, сравниваются по значению, а не по ссылке.

```go
// Email — value object с валидацией на уровне домена
type Email struct{ value string }

func NewEmail(email string) (Email, error) { /* валидация */ }
func (e Email) Equals(other Email) bool   { return strings.EqualFold(e.value, other.value) }
```

**Domain Events** — агрегаты публикуют события через интерфейс `DomainEvent` (`domain/shared/events`). Конкретные события: `StudentRegisteredEvent`, `ProfileCompletedEvent`, `AdminRegisteredEvent`. Это позволяет реагировать на изменения состояния без прямых зависимостей между контекстами.

**Repository interfaces** — репозитории (`StudentRepository`, `TrackRepository`, `TrackSelectionRepository` и др.) объявлены как интерфейсы внутри домена. Инфраструктурный слой (`infrastructure/persistence/postgres`) предоставляет конкретные реализации — домен не знает ни о PostgreSQL, ни о pgx.

**Domain Service** — `ProfileChecker` (`domain/student`) — доменный сервис, инкапсулирующий инвариант «профиль считается заполненным»: все оценки от 2 до 5, хотя бы один навык > 0, указан стиль обучения и хотя бы одна профессиональная цель. Используется как в use case обновления предпочтений, так и в use case получения рекомендаций.

### Event Bus

При регистрации use case публикует событие (`student.registered` / `admin.registered`) в in-memory шину (`domain/shared/events/EventBus`). Подписчики (`infrastructure/eventbus/subscribers`) асинхронно создают связанные сущности: профиль студента и запись `ProfileCompletion`. Это разрывает транзакционную связность между агрегатами из разных контекстов.

### PROMETHEE

Алгоритм рекомендаций (`domain/selection/promethee.go`) реализован как чистая доменная логика без зависимостей на инфраструктуру. Каждый трек попарно сравнивается с остальными по каждому критерию; формируется положительный и отрицательный поток предпочтений; итоговый score (0–1) определяет место трека в рейтинге.

### Seeding

При старте `bootstrap/app.go` вызывает `postgres.SeedTracks()` — если таблица треков пуста, засеваются дефолтные треки. Это позволяет сразу запустить сервис с реальными данными без отдельной команды.

---

## Скриншоты

<div align="center">

### Главная страница — каталог треков
![Каталог треков](./photo/main.png)
*Список доступных образовательных треков с фильтрацией и поиском*

---

### Регистрация и вход
![Форма регистрации и входа](./photo/auth.png)
*Модальное окно аутентификации с вкладками «Вход» и «Регистрация»*

---

### Заполнение профиля студента
![Профиль студента](./photo/profile.png)
*Форма ввода оценок, навыков, целей и стиля обучения*

---

### Рекомендации треков
![Рекомендации](./photo/recommendations.png)
*Персонализированный рейтинг треков по методу PROMETHEE*

---

### Детальная карточка трека
![Карточка трека](./photo/track_details.png)
*Подробная информация о треке: учебный план, требования, преподаватели*

---

### Тёмная тема
![Тёмная тема](./photo/dark_theme.png)
*Интерфейс в тёмном режиме*

---

### Панель администратора — управление треками
![Управление треками](./photo/admin.png)
*Создание, редактирование и удаление треков (только для admin)*

</div>
