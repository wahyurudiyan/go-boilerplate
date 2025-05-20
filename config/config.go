package config

import (
	"time"

	"github.com/wahyurudiyan/go-boilerplate/pkg/redis"
	"github.com/wahyurudiyan/go-boilerplate/pkg/sql"
)

type ServiceConfig struct {
	// Application Configuration
	ApplicationName        string `mapstructure:"APPLICATION_NAME"`
	ApplicationEngine      string `mapstructure:"APPLICATION_EGINE"` // currently only support Gin (default) and Fiber
	ApplicationVersion     string `mapstructure:"APPLICATION_VERSION"`
	ApplicationTribeName   string `mapstructure:"APPLICATION_TRIBE_NAME"`
	ApplicationEnvrionment string `mapstructure:"APPLICATION_ENVRIONMENT"`

	GrpcPort    string        `mapstructure:"GRPC_PORT"`
	GrpcTimeout time.Duration `mapstructure:"GRPC_TIMEOUT"`

	RestPort               string `mapstructure:"REST_PORT"`
	RestBodyLimit          int64  `mapstructure:"REST_BODY_LIMIT"`           // [FIBER ONLY] in megabyte
	RestStrictRoute        bool   `mapstructure:"REST_STRICT_ROUTE"`         // [FIBER ONLY] /foo and /foo/ is different when enabled
	RestReadTimeout        int64  `mapstructure:"REST_READ_TIMEOUT"`         // [FIBER ONLY] in seconds, 0 is unlimited
	RestWriteTimeout       int64  `mapstructure:"REST_WRITE_TIMEOUT"`        // [FIBER ONLY] in seconds, 0 is unlimited
	RestIdleTimeout        int64  `mapstructure:"REST_IDLE_TIMEOUT"`         // [FIBER ONLY] in seconds, 0 is unlimited
	RestRouteCaseSensitive bool   `mapstructure:"REST_ROUTE_CASE_SENSITIVE"` // [FIBER ONLY] /Foo and /foo is different when enabled

	TelemetryMeterInterval      time.Duration `mapstructure:"TELEMETRY_METER_INTERVAL"`
	TelemetryEnableRuntimeMeter bool          `mapstructure:"TELEMETRY_ENABLE_RUNTIME_METER"`

	Redis    redis.RedisConfig `mapstructure:",squash"`
	Database sql.SQLConfig     `mapstructure:",squash"`
}
