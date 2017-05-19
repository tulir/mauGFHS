package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"maunium.net/go/maulogger"

	"gopkg.in/yaml.v2"
)

// Config is the base configuration container.
type Config struct {
	Database  DBConfig         `yaml:"database"`
	Addresses []ListenLocation `yaml:"listen"`
	Logging   LogConfig        `yaml:"logging"`
}

// LogConfig contains logging configurations
type LogConfig struct {
	Directory      string `yaml:"directory"`
	FileNameFormat string `yaml:"fileNameFormat"`
	FileMode       uint32 `yaml:"fileMode"`
}

// GetFileFormat returns a mauLogger-compatible logger file format based on the data in the struct.
func (lc LogConfig) GetFileFormat() maulogger.LoggerFileFormat {
	path := lc.FileNameFormat
	if len(lc.Directory) > 0 {
		path = lc.Directory + "/" + path
	}

	return func(now string, i int) string {
		return fmt.Sprintf(path, now, i)
	}
}

// Configure configures a mauLogger instance with the data in this struct.
func (lc LogConfig) Configure(log *maulogger.Logger) {
	log.FileFormat = lc.GetFileFormat()
	log.FileMode = os.FileMode(lc.FileMode)
}

// ListenLocation is a location where the server should listen.
type ListenLocation struct {
	Address      string `yaml:"address"`
	Port         uint8  `yaml:"port"`
	TrustHeaders bool   `yaml:"trustHeaders"`
}

// DBConfig contains connection information for the database.
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     uint8  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// GetDSN gets the SQL data source name for this database config.
func (config DBConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
}

var config Config

// GetConfig returns the singleton Config instance.
func GetConfig() Config {
	return config
}

// Open opens the config at the given path.
func Open(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &config)
}
