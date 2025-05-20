package sql

import "time"

type SQLConfig struct {
	// Database Configuration
	DatabaseName      string `mapstructure:"database_name"`
	DatabaseHost      string `mapstructure:"database_host"`
	DatabasePort      string `mapstructure:"database_port"`
	DatabaseDriver    string `mapstructure:"database_driver"` // e.g., "mysql, postgres, etc"
	DatabaseCharset   string `mapstructure:"database_charset"`
	DatabaseUsername  string `mapstructure:"database_username"`
	DatabasePassword  string `mapstructure:"database_password"`
	DatabaseParsetime bool   `mapstructure:"database_parsetime"`

	DatabaseMaxOpenConnection     int           `mapstructure:"database_max_open_connection"`
	DatabaseMaxIdleConnection     int           `mapstructure:"database_max_idle_connection"`
	DatabaseMaxIdleTimeConnection time.Duration `mapstructure:"database_max_idle_time_connection"` // in seconds or minutes as needed
	DatabaseMaxLifetimeConnection time.Duration `mapstructure:"database_max_lifetime_connection"`  // in seconds or minutes as needed
}
