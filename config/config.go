package config

type ServiceConfig struct {
	// Application Configuration
	ApplicationName        string `mapstructure:"APPLICATION_NAME"`
	ApplicationEnvrionment string `mapstructure:"APPLICATION_ENVRIONMENT"`

	RestPort               string `mapstructure:"REST_PORT"`
	RestBodyLimit          int    `mapstructure:"REST_BODY_LIMIT"`           // in megabyte
	RestStrictRoute        bool   `mapstructure:"REST_STRICT_ROUTE"`         // /foo and /foo/ is different when enabled
	RestReadTimeout        int    `mapstructure:"REST_READ_TIMEOUT"`         // in seconds, 0 is unlimited
	RestWriteTimeout       int    `mapstructure:"REST_WRITE_TIMEOUT"`        // in seconds, 0 is unlimited
	RestRouteCaseSensitive bool   `mapstructure:"REST_ROUTE_CASE_SENSITIVE"` // /Foo and /foo is different when enabled
}
