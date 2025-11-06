package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required,unset"`
	DB       string `env:"POSTGRES_DB,required"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" envDefault:"prefer"`
	TZ       string `env:"POSTGRES_TZ" envDefault:"UTC"`
	AppName  string `env:"POSTGRES_APP_NAME" envDefault:"go-hex"`

	PoolMaxConns              int           `env:"PGX_POOL_MAX_CONNS" envDefault:"4"`
	PoolMinConns              int           `env:"PGX_POOL_MIN_CONNS" envDefault:"0"`
	PoolMaxConnLifetime       time.Duration `env:"PGX_POOL_MAX_CONN_LIFETIME" envDefault:"1h"`
	PoolMaxConnIdleTime       time.Duration `env:"PGX_POOL_MAX_CONN_IDLE_TIME" envDefault:"30m"`
	PoolHealthCheck           time.Duration `env:"PGX_POOL_HEALTH_CHECK" envDefault:"1m"`
	PoolMaxConnLifetimeJitter time.Duration `env:"PGX_POOL_MAX_CONN_LIFETIME_JITTER" envDefault:"0s"`
	QueryDebug                bool          `env:"POSTGRES_QUERY_DEBUG" envDefault:"false"`
	SlowQueryThreshold        time.Duration `env:"POSTGRES_SLOW_QUERY_THRESHOLD" envDefault:"200ms"`
}

func (p *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		p.Host, p.Port, p.User, p.Password, p.DB, p.SSLMode, p.TZ,
	)
}

func (p *Config) PGXDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&"+
			"TimeZone=%s&"+
			"pool_max_conns=%d&"+
			"pool_min_conns=%d&"+
			"pool_max_conn_lifetime=%s&"+
			"pool_max_conn_idle_time=%s&"+
			"pool_health_check_period=%s&"+
			"pool_max_conn_lifetime_jitter=%s&"+
			"application_name=%s",
		p.User, p.Password, p.Host, p.Port, p.DB, p.SSLMode,
		p.TZ,
		p.PoolMaxConns,
		p.PoolMinConns,
		p.PoolMaxConnLifetime,
		p.PoolMaxConnIdleTime,
		p.PoolHealthCheck,
		p.PoolMaxConnLifetimeJitter,
		p.AppName,
	)
}
