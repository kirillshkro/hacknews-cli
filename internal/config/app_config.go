package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
)

// AppConfig представляет полную конфигурацию микросервиса.
// Все значения загружаются из переменных окружения через библиотеку cleanenv.
type AppConfig struct {
	// APIConfig содержит параметры HTTP API (Gin).
	API APIConfig `env:"API"`

	// DatabaseConfig содержит параметры подключения к PostgreSQL.
	Database DatabaseConfig `env:"DATABASE"`

	// ScraperConfig содержит параметры фоновой синхронизации с Firebase API.
	Scraper ScraperConfig `env:"SCRAPER"`

	// CacheConfig содержит параметры кэширования.
	Cache CacheConfig `env:"CACHE"`
}

// APIConfig описывает параметры HTTP-сервера Gin.
type APIConfig struct {
	// Port — порт, на котором слушает HTTP-сервер.
	Port int `env:"PORT" env-default:"8080"`

	// Mode — режим работы Gin (debug, release, test).
	Mode string `env:"MODE" env-default:"release"`

	// ReadTimeout — таймаут чтения запроса (для мобильных сетей).
	ReadTimeout time.Duration `env:"READ_TIMEOUT" env-default:"3s"`

	// WriteTimeout — таймаут записи ответа.
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" env-default:"3s"`

	// IdleTimeout — таймаут простоя соединения.
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`

	// MaxRequestBodySize — максимальный размер тела запроса в байтах.
	MaxRequestBodySize int64 `env:"MAX_REQUEST_BODY_SIZE" env-default:"10485760"` // 10MB
}

// DatabaseConfig описывает параметры подключения к PostgreSQL.
type DatabaseConfig struct {
	// Host — хост PostgreSQL.
	Host string `env:"HOST" env-default:"localhost"`

	// Port — порт PostgreSQL.
	Port int `env:"PORT" env-default:"5432"`

	// User — имя пользователя.
	User string `env:"USER" env-default:"postgres"`

	// Password — пароль пользователя.
	Password string `env:"PASSWORD" env-default:"postgres"`

	// DBName — имя базы данных.
	DBName string `env:"DB_NAME" env-default:"hackernews"`

	// SSLMode — режим SSL (disable, require, verify-ca, verify-full).
	SSLMode string `env:"SSL_MODE" env-default:"disable"`

	// MaxOpenConns — максимальное количество открытых соединений.
	MaxOpenConns int `env:"MAX_OPEN_CONNS" env-default:"25"`

	// MaxIdleConns — максимальное количество connections в пуле ожидания.
	MaxIdleConns int `env:"MAX_IDLE_CONNS" env-default:"5"`

	// ConnMaxLifetime — максимальное время жизни соединения.
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" env-default:"1h"`

	// ConnMaxIdleTime — максимальное время простоя соединения.
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME" env-default:"10m"`
}

// RedisConfig описывает параметры подключения к Redis.
type RedisConfig struct {
	// Host — хост Redis.
	Host string `env:"HOST" env-default:"localhost"`

	// Port — порт Redis.
	Port int `env:"PORT" env-default:"6379"`

	// Password — пароль (если есть).
	Password string `env:"PASSWORD" env-default:""`

	// DB — номер базы данных.
	DB int `env:"DB" env-default:"0"`

	// PoolSize — размер пула соединений.
	PoolSize int `env:"POOL_SIZE" env-default:"10"`

	// MinIdleConns — минимальное количество connections в ожидании.
	MinIdleConns int `env:"MIN_IDLE_CONNS" env-default:"2"`

	// DialTimeout — таймаут установки соединения.
	DialTimeout time.Duration `env:"DIAL_TIMEOUT" env-default:"5s"`

	// ReadTimeout — таймаут чтения из Redis.
	ReadTimeout time.Duration `env:"READ_TIMEOUT" env-default:"2s"`

	// WriteTimeout — таймаут записи в Redis.
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" env-default:"2s"`

	// PoolTimeout — таймаут получения соединения из пула.
	PoolTimeout time.Duration `env:"POOL_TIMEOUT" env-default:"4s"`
}

// ScraperConfig описывает параметры фоновой синхронизации.
type ScraperConfig struct {
	// APIBaseURL — базовый URL Firebase Hacker News API.
	APIBaseURL string `env:"API_BASE_URL" env-default:"https://hacker-news.firebaseio.com/v0"`

	// SyncInterval — интервал синхронизации списков (в секундах).
	SyncInterval int `env:"SYNC_INTERVAL" env-default:"30"`

	// WorkerPoolSize — количество воркеров в пуле для скачивания отдельных элементов.
	WorkerPoolSize int `env:"WORKER_POOL_SIZE" env-default:"10"`

	// RequestTimeout — таймаут запроса к Firebase API (в секундах).
	RequestTimeout int `env:"REQUEST_TIMEOUT" env-default:"3"`

	// BatchSize — количество элементов для параллельной обработки в одной пачке.
	BatchSize int `env:"BATCH_SIZE" env-default:"50"`

	// RetryAttempts — количество попыток повторной передачи при ошибке.
	RetryAttempts int `env:"RETRY_ATTEMPTS" env-default:"3"`

	// RetryDelay — задержка между повторными попытками (в секундах).
	RetryDelay int `env:"RETRY_DELAY" env-default:"2"`
}

// CacheConfig описывает параметры кэширования.
type CacheConfig struct {
	// StoriesListTTL — TTL для кэша списков статей (в секундах).
	// По ТЗ: 60 секунд.
	StoriesListTTL int `env:"STORIES_LIST_TTL" env-default:"60"`

	// CommentTreeTTL — TTL для кэша "горячих" веток обсуждений (в секундах).
	// По ТЗ: 5 минут.
	CommentTreeTTL int `env:"COMMENT_TREE_TTL" env-default:"300"`

	// HotStoriesCount — количество "горячих" статей, деревья которых кешируются.
	HotStoriesCount int `env:"HOT_STORIES_COUNT" env-default:"10"`

	// InvalidationBuffer — буфер времени для инвалидации кэша (в секундах).
	InvalidationBuffer int `env:"INVALIDATION_BUFFER" env-default:"5"`
}

// Load загружает конфигурацию из переменных окружения.
// Возвращает ошибку, если какие-либо обязательные переменные не заданы.
func Load() (*AppConfig, error) {
	cfg := &AppConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// DSN возвращает DSN для подключения к PostgreSQL через sqlx.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// SQLXPoolConfig возвращает конфигурацию пула соединений для sqlx.
// Данная функция должна вызываться после установки соединения.
// Возвращаем nil здесь как заглушку, реальная логика будет в internal/db.
func (c *DatabaseConfig) SQLXPoolConfig() *sqlx.DB {
	return nil
}
