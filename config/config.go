package config

import "fmt"

// Config is the base configuration container.
type Config struct {
	Database  DBConfig         `yaml:"database"`
	Addresses []ListenLocation `yaml:"listen"`
}

// ListenLocation is a location where the server should listen.
type ListenLocation struct {
	Address      string `yaml:"address"`
	Port         uint8  `yaml:"port"`
	TrustHeaders bool   `yaml:"trust-headers"`
}

// DBConfig contains connection information for the database.
type DBConfig struct {
	Address  string `yaml:"address"`
	Port     uint8  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// GetDSN gets the SQL data source name for this database config.
func (config DBConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Address, config.Port, config.Database)
}
