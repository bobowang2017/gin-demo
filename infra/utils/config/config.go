package config

import "time"

type System struct {
	LogLevel string `yaml:"log-level"`
}

type Server struct {
	RunMode         string
	HttpPort        string
	LogSavePath     string
	LogMaxAge       time.Duration
	LogRotationTime time.Duration
	LogLevel        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	EnabledSwagger  bool
}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

type Redis struct {
	Host           string
	Password       string
	MaxIdle        int
	MaxActive      int
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}