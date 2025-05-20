package config

import "time"

type ServiceConfig struct {
	// Application Configuration
	ApplicationName        string `mapstructure:"APPLICATION_NAME"`
	ApplicationEngine      string `mapstructure:"APPLICATION_EGINE"` // currently only support Gin (default) and Fiber
	ApplicationVersion     string `mapstructure:"APPLICATION_VERSION"`
	ApplicationTribeName   string `mapstructure:"APPLICATION_TRIBE_NAME"`
	ApplicationEnvrionment string `mapstructure:"APPLICATION_ENVRIONMENT"`

	RestPort               string `mapstructure:"REST_PORT"`
	RestBodyLimit          int64  `mapstructure:"REST_BODY_LIMIT"`           // in megabyte
	RestStrictRoute        bool   `mapstructure:"REST_STRICT_ROUTE"`         // /foo and /foo/ is different when enabled
	RestReadTimeout        int64  `mapstructure:"REST_READ_TIMEOUT"`         // in seconds, 0 is unlimited
	RestWriteTimeout       int64  `mapstructure:"REST_WRITE_TIMEOUT"`        // in seconds, 0 is unlimited
	RestIdleTimeout        int64  `mapstructure:"REST_IDLE_TIMEOUT"`         // in seconds, 0 is unlimited
	RestRouteCaseSensitive bool   `mapstructure:"REST_ROUTE_CASE_SENSITIVE"` // /Foo and /foo is different when enabled

	TelemetryMeterInterval      time.Duration `mapstructure:"TELEMETRY_METER_INTERVAL"`
	TelemetryEnableRuntimeMeter bool          `mapstructure:"TELEMETRY_ENABLE_RUNTIME_METER"`
}
