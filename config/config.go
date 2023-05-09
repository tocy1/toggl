package config

import (
	"fmt"
	"os"
	"strconv"

	"code.cloudfoundry.org/lager"
)

type MariaDB struct {
	DSN      string
	Username string
	Password string
	Database string
	Host     string
	Port     int
}

type Server struct {
	BindAddress   string
	AdminUsername string
	AdminPassword string
}

const (
	envKeyLogLevel          envKey = "LOG_LEVEL"
	envKeyServerBindAddress envKey = "BIND_ADDRESS"
	envKeyServerAdminUser   envKey = "ADMIN_USER"
	envKeyServerAdminPass   envKey = "ADMIN_PASS"

	envKeyMariaDBHost     envKey = "MARIADB_HOST"
	envKeyMariaDBPort     envKey = "MARIADB_PORT"
	envKeyMariaDBDatabase envKey = "MARIADB_DB"
	envKeyMariaDBUser     envKey = "MARIADB_USER"
	envKeyMariaDBPass     envKey = "MARIADB_PASS"
)

type Config struct {
	MariaDB MariaDB
	Server  Server
}

type envKey string

func (k envKey) String() string {
	return string(k)
}

func LogLevel() lager.LogLevel {
	if value, ok := os.LookupEnv(envKeyLogLevel.String()); ok {
		level, err := lager.LogLevelFromString(value)
		if err == nil {
			return level
		}
		fmt.Printf("Invalid log level configured in Environment %v: %v. Falling back to INFO.", envKeyLogLevel, value)
	} else {
		fmt.Printf("No log level configured in Environment %v. Falling back to INFO.", envKeyLogLevel)
	}

	return lager.INFO
}

func GetConfig(logger lager.Logger) Config {
	l := logger.Session("configuration")
	c := Config{}
	c.readFromEnvironment()
	c.setDefaultValues()
	l.Info("read", lager.Data{"config": c})
	return c
}

func (c *Config) readFromEnvironment() {

	setIfEnvNotEmpty(&c.MariaDB.Host, envKeyMariaDBHost)
	setIfEnvNotEmptyInt(&c.MariaDB.Port, envKeyMariaDBPort)
	setIfEnvNotEmpty(&c.MariaDB.Database, envKeyMariaDBDatabase)
	setIfEnvNotEmpty(&c.MariaDB.Password, envKeyMariaDBPass)
	setIfEnvNotEmpty(&c.MariaDB.Username, envKeyMariaDBUser)

	setIfEnvNotEmpty(&c.Server.BindAddress, envKeyServerBindAddress)
	setIfEnvNotEmpty(&c.Server.AdminUsername, envKeyServerAdminUser)
	setIfEnvNotEmpty(&c.Server.AdminPassword, envKeyServerAdminPass)

}

func setIfEnvNotEmpty(to *string, key envKey) {
	if value, ok := os.LookupEnv(key.String()); ok {
		*to = value
	}
}

func setIfEnvNotEmptyInt(to *int, key envKey) {
	if value, ok := os.LookupEnv(key.String()); ok {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			*to = intValue
		}
	}
}

func (c *Config) setDefaultValues() {
	setDefault(&c.Server.BindAddress, ":8080")
}

func setDefault(to *string, value string) {
	if *to == "" {
		*to = value
	}
}
