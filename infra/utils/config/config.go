package config

import (
	"time"
)

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
	MaxIdleCon  int
	MaxOpenCon  int
	MaxIdleTime int64
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

type AuthCodeJwt struct {
	Secret string `yaml:"secret"`
	Issuer string `yaml:"issuer"`
	Expire int    `yaml:"expire"`
}

type QiNiuYun struct {
	Host          string `yaml:"host"`
	AccessKey     string `yaml:"accessKey"`
	SecretKey     string `yaml:"secretKey"`
	TopicBucket   string `yaml:"topicBucket"`
	CommentBucket string `yaml:"commentBucket"`
	ReplyBucket   string `yaml:"replyBucket"`
	DefaultExpire uint64 `yaml:"defaultExpire"`
}

type RabbitMq struct {
	Host       string `yaml:"host"`
	Exchange   string `yaml:"exchange"`
	RoutingKey string `yaml:"routingKey"`
	Queue      string `yaml:"queue"`
}

type MqConsumer struct {
	Host       string `yaml:"host"`
	Exchange   string `yaml:"exchange"`
	RoutingKey string `yaml:"routingKey"`
	Queue      string `yaml:"queue"`
}
