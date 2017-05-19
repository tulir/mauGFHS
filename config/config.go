// mauGFHS - A server that can serve as a backend for many kinds of services that only require file hosting.
// Copyright (C) 2017 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
	Database DBConfig       `yaml:"database"`
	Listen   ListenLocation `yaml:"listen"`
	Logging  LogConfig      `yaml:"logging"`
	DataPath string         `yaml:"dataPath"`
}

// LogConfig contains logging configurations
type LogConfig struct {
	Directory       string `yaml:"directory"`
	FileNameFormat  string `yaml:"fileNameFormat"`
	FileDateFormat  string `yaml:"fileDateFormat"`
	FileMode        uint32 `yaml:"fileMode"`
	TimestampFormat string `yaml:"timestampFormat"`
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
	log.TimeFormat = lc.TimestampFormat
}

// ListenLocation is a location where the server should listen.
type ListenLocation struct {
	Address      string `yaml:"address"`
	Port         uint8  `yaml:"port"`
	TrustHeaders bool   `yaml:"trustHeaders"`
	PathPrefix   string `yaml:"pathPrefix"`
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

// MainConfig is the main singleton config instance.
var MainConfig = &Config{}

// Open opens the config at the given path.
func Open(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, MainConfig)
}
