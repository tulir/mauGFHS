package config

import "fmt"

// Config is the base configuration container.
type Config struct {
	Database DBConfig `yaml:"database"`
}

type NamespaceConfig struct {
	Namespace  string `yaml:"namespace"`
	Extensions string
}

// DBConfig contains connection information for the database.
type DBConfig struct {
	Address  string `yaml:"address"`
	Port     uint8  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (config DBConfig) GetDSN() string {
	fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Address, config.Port, config.Database)
}
