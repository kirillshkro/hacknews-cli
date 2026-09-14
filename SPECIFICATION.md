# Техническое задание (ТЗ)
## Микросервис-агрегатор Hacker News API (Golang / Gin / sqlx / PostgreSQL / Redis)

### 1. Архитектурная концепция и Стек
Проект реализуется в виде изолированного микросервиса на языке **Golang**, обеспечивающего API для веб- и мобильных платформ. Приложение работает по принципу **асинхронного агрегатора**: запросы пользователей обслуживаются мгновенно из локальной инфраструктуры (Redis/PostgreSQL), пока фоновые процессы непрерывно синхронизируют данные с Firebase Hacker News API.

*   **API Framework:** Gin Gonic (высокая производительность, удобные Middleware).
*   **Database Mapper:** `jmoiron/sqlx` (типизированный маппинг в Go-структуры, контроль сырых SQL-запросов).
*   **Основная СУБД:** PostgreSQL 15+ (эффективная работа с древовидными структурами через CTE и индексация).
*   **Кэш / Брокер:** Redis (снижение нагрузки на СУБД, координация распределенных воркеров).

---

### 2. Схема хранения данных (PostgreSQL)

Для бесшовного маппинга через `sqlx` используется единая таблица `items`, повторяющая полиморфную сущность оригинального API, с добавлением связей.

```sql
CREATE TYPE item_type AS ENUM ('story', 'comment', 'job', 'poll', 'pollopt');

CREATE TABLE items (
    id BIGINT PRIMARY KEY,
    deleted BOOLEAN DEFAULT FALSE,
    type item_type NOT NULL,
    by_user VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    text TEXT,
    dead BOOLEAN DEFAULT FALSE,
    parent_id BIGINT REFERENCES items(id),
    poll_id BIGINT REFERENCES items(id),
    url TEXT,
    score INT DEFAULT 0,
    title VARCHAR(500),
    descendants INT DEFAULT 0, -- Общее кол-во комментариев
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для быстрой выборки списков и построения деревьев
CREATE INDEX idx_items_type_score ON items(type, score DESC) WHERE deleted = FALSE AND dead = FALSE;
CREATE INDEX idx_items_parent ON items(parent_id);
CREATE INDEX idx_items_created ON items(created_at DESC);
```

---

### 3. Функциональные компоненты системы

#### 3.1. Фоновый слой (Background Sync Worker Pool)
*   **Пул Горутин:** На старте приложения запускается воркер-пул (на базе `golang.org/x/sync/errgroup` или каналов).
*   **Синхронизация списков:** Раз в 30 секунд воркер запрашивает `/topstories.json` и `/newstories.json`. Полученный массив ID (до 500 штук) временно сохраняется в **Redis** (ключ `hn:top_ids`).
*   **Потоковый парсинг (Scraper):** Воркер сравнивает ID из сети с ID в PostgreSQL. Для новых или изменившихся ID запускаются горутины, которые скачивают индивидуальные JSON объекта (`/item/{id}.json`), парсят их и делают `INSERT ... ON CONFLICT (id) DO UPDATE` через `sqlx`.
*   **Распределенная блокировка:** Перед началом цикла синхронизации воркер выставляет ключ в Redis через `SET NX PX 5000` (Mutex), чтобы инстансы приложения в кластере не дублировали запросы к Firebase API.

#### 3.2. Слой кэширования (Redis)
Внедряется двухслойное кэширование:
1.  **Кэш списков (Краткосрочный):** Эндпоинты `/api/v1/stories/top` и `/api/v1/stories/new` кэшируют сериализованный в JSON ответ в Redis на **60 секунд**.
2.  **Кэш "Горячих" веток обсуждений:** Дерево комментариев для топ-10 статей кэшируется в Redis на **5 минут**. При добавлении новых комментариев воркером, кэш конкретной статьи инвалидируется (сбрасывается).

#### 3.3. HTTP API Слой (Gin Контроллеры)
Реализует эндпоинты для фронтенда и мобильного приложения.

*   **`GET /api/v1/stories?type=[top|new]&page=1&limit=20`**
    *   *Логика:* Проверяет наличие страницы в кэше Redis. Если нет — идет в PostgreSQL, делает выборку с `LIMIT` и `OFFSET`, пишет в Redis, отдает клиенту.
*   **`GET /api/v1/stories/:id/comments`**
    *   *Логика:* Запрашивает дерево комментариев. Из PostgreSQL данные вытягиваются рекурсивным запросом (CTE), собирающим все дочерние элементы для `parent_id = :id`. Из плоского массива Go-код собирает древовидную структуру JSON.

---

### 4. Оптимизация под Mobile & Web клиенты

1.  **Сжатие данных (Трафик):** Подключить middleware `github.com/gin-contrib/gzip`. Длинные текстовые ветки комментариев Hacker News сжимаются до 70-80%, что критично для мобильных сетей.
2.  **Экономия CPU:** Заменить стандартный пакет `encoding/json` на высокопроизводительный **`github.com/bytedance/sonic`** (или `json-iterator`), прописав тег сборки в Gin. Это ускорит сериализацию тяжелых деревьев комментариев.
3.  **Таймауты и Контексты:** Все запросы от Gin к PostgreSQL и Redis, а также запросы воркера к Firebase API обязаны принимать `context.Context` с жестким `WithTimeout` (до 2-3 секунд), исключая утечки ресурсов и горутин.
