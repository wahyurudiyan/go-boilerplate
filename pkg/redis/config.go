package redis

import (
	"time"
)

type RedisConfig struct {
	// Redis Configuration
	RedisDB           int      `mapstructure:"REDIS_DB"`
	RedisAddr         string   `mapstructure:"REDIS_ADDR"`
	RedisUsername     string   `mapstructure:"REDIS_USERNAME"`
	RedisPassword     string   `mapstructure:"REDIS_PASSWORD"`
	RedisIsCluster    bool     `mapstructure:"REDIS_IS_CLUSTER"`
	RedisClusterAddrs []string `mapstructure:"REDIS_CLUSTER_ADDRS"`

	// Connection pool settings
	RedisPoolSize        int           `mapstructure:"REDIS_POOL_SIZE"`
	RedisMinIdleConns    int           `mapstructure:"REDIS_MIN_IDLE_CONNS"`
	RedisPoolTimeout     time.Duration `mapstructure:"REDIS_POOL_TIMEOUT"`
	RedisMaxRetries      int           `mapstructure:"REDIS_MAX_RETRIES"`
	RedisMinRetryBackoff time.Duration `mapstructure:"REDIS_MIN_RETRY_BACKOFF"`
	RedisMaxRetryBackoff time.Duration `mapstructure:"REDIS_MAX_RETRY_BACKOFF"`

	// Connection health and monitoring
	RedisDialTimeout  time.Duration `mapstructure:"REDIS_DIAL_TIMEOUT"`
	RedisReadTimeout  time.Duration `mapstructure:"REDIS_READ_TIMEOUT"`
	RedisWriteTimeout time.Duration `mapstructure:"REDIS_WRITE_TIMEOUT"`
	RedisIdleTimeout  time.Duration `mapstructure:"REDIS_IDLE_TIMEOUT"`

	// TLS Configuration
	RedisEnableTLS             bool `mapstructure:"REDIS_ENABLE_TLS"`
	RedisTLSInsecureSkipVerify bool `mapstructure:"REDIS_TLS_INSECURE_SKIP_VERIFY"`

	// Cluster specific configurations
	RedisMaxRedirects   int  `mapstructure:"REDIS_MAX_REDIRECTS"`
	RedisRouteByLatency bool `mapstructure:"REDIS_ROUTE_BY_LATENCY"`
	RedisRouteRandomly  bool `mapstructure:"REDIS_ROUTE_RANDOMLY"`
	RedisReadOnly       bool `mapstructure:"REDIS_READ_ONLY"`
}
