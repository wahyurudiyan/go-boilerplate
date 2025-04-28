package sql

type SQLConfig struct {
	// Database Configuration
	DatabaseName      string `mapstructure:"DATABASE_NAME"`
	DatabaseHost      string `mapstructure:"DATABASE_HOST"`
	DatabasePort      string `mapstructure:"DATABASE_PORT"`
	DatabaseDriver    string `mapstructure:"DATABASE_DRIVER"` // e.g., "mysql, postgres, etc"
	DatabaseCharset   string `mapstructure:"DATABASE_CHARSET"`
	DatabaseUsername  string `mapstructure:"DATABASE_USERNAME"`
	DatabasePassword  string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseParsetime bool   `mapstructure:"DATABASE_PARSETIME"`

	DatabaseMaxOpenConnection     int `mapstructure:"DATABASE_MAX_OPEN_CONNECTION"`
	DatabaseMaxIdleConnection     int `mapstructure:"DATABASE_MAX_IDLE_CONNECTION"`
	DatabaseMaxIdleTimeConnection int `mapstructure:"DATABASE_MAX_IDLE_TIME_CONNECTION"` // in seconds or minutes as needed
	DatabaseMaxLifetimeConnection int `mapstructure:"DATABASE_MAX_LIFETIME_CONNECTION"`  // in seconds or minutes as needed
}
