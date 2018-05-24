package config

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/labstack/gommon/log"
)

var (
	Config config // holds the global app config.
	defaultConfigFile = "config/config.toml"
)

type config struct {
	ReleaseMode bool   `toml:"release_mode"`
	LogLevel string `toml:"log_level"`

	SessionStore string `toml:"session_store"`
	CacheStore string `toml:"cache_store"`

	// Application configuration
	App app

	// template
	Tmpl tmpl

	Server server

	// MySQL
	DB database `toml:"database"`

	// static resources
	Static static

	// Redis
	Redis redis

	// Memcached
	Memcached memcached

	// Opentracing
	Opentracing opentracing
}

type app struct {
	Name string `toml:"name"`
}

type server struct {
	Graceful bool   `toml:"graceful"`
	Addr     string `toml:"addr"`

	DomainApi    string `toml:"domain_api"`
	DomainWeb    string `toml:"domain_web"`
	DomainSocket string `toml:"domain_socket"`
}

type static struct {
	Type string `toml:"type"`
}

type tmpl struct {
	Type   string `toml:"type"`
	Data   string `toml:"data"`
	Dir    string `toml:"dir"`
	Suffix string `toml:"suffix"` // .html,.tpl
}

type database struct {
	Name     string `toml:"name"`
	UserName string `toml:"user_name"`
	Pwd      string `toml:"pwd"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
}

type redis struct {
	Server string `toml:"server"`
	Pwd    string `toml:"pwd"`
}

func init() {
}

// initConfig initializes the app configuration by first setting defaults,
// then overriding settings from the app config file, then overriding
// It returns an error if any.
func InitConfig(configFile string) error {
	if configFile == "" {
		configFile = defaultConfigFile
	}

	// Set defaults.
	Config = config{
		ReleaseMode: false,
		LogLevel:    "DEBUG",
	}

	if _, err := os.Stat(configFile); err != nil {
		return errors.New("config file err:" + err.Error())
	} else {
		log.Infof("load config from file:" + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return errors.New("config load err:" + err.Error())
		}
		_, err = toml.Decode(string(configBytes), &Config)
		if err != nil {
			return errors.New("config decode err:" + err.Error())
		}
	}

	log.Infof("config data:%v", Config)

	return nil
}

func GetLogLvl() log.Lvl {
	//DEBUG INFO WARN ERROR OFF
	switch Config.LogLevel {
	case "DEBUG":
		return log.DEBUG
	case "INFO":
		return log.INFO
	case "WARN":
		return log.WARN
	case "ERROR":
		return log.ERROR
	case "OF":
		return log.OFF
	}

	return log.DEBUG
}

const (
	// Template Type
	TEMPLATE = "TEMPLATE"

	// File
	FILE = "FILE"

	// Redis
	REDIS = "REDIS"

	// Cookie
	COOKIE = "COOKIE"

	// In Memory
	IN_MEMORY = "IN_MEMORY"
)
