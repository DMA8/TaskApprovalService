package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	defaultConfig = "config/config.yaml"
)

type MongoConfig struct {
	URI             string `yaml:"uri"`
	URIFul          string `yaml:"uri_full"`
	TasksCollection string `yaml:"tasks_collection"`
	DB              string `yaml:"db"`
	Login           string `yaml:"login"`
	Password        string `yaml:"password"`
}

type HTTPConfig struct {
	URI               string `yaml:"uri"`
	AccessCookieName  string `yaml:"access_cookie_name"`
	RefreshCookieName string `yaml:"refresh_cookie_name"`
	APIVersion        string `yaml:"api_version"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type AuthGRPCConfig struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Transport string `yaml:"transport"`
}

type TasksEventGRPCConfig struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Transport string `yaml:"transport"`
}

type KafkaConfig struct {
	URL       string `yaml:"url"`
	TaskTopic string `yaml:"task_topic"`
	MailTopic string `yaml:"mail_topic"`
	GroupID   string `yaml:"group_id"`
	Partition int    `yaml:"partition"`
}

type Config struct {
	HTTP            HTTPConfig           `yaml:"http_server"`
	Mongo           MongoConfig          `yaml:"mongo"`
	Log             LogConfig            `yaml:"logging"`
	AuthGRPC        AuthGRPCConfig       `yaml:"auth_grpc"`
	TasksEventQueue TasksEventGRPCConfig `yaml:"tasks_events_grpc"`
	Kafka           KafkaConfig          `yaml:"kafka"`
}

var once sync.Once
var configG *Config

func NewConfig() *Config {
	var cfgPath string
	once.Do(func() {
		if c := os.Getenv("CFG_PATH"); c != "" {
			cfgPath = c
		} else {
			log.Println("if you are running localy -> export CFG_PATH=config/config_debug.yaml")
			cfgPath = defaultConfig
		}
		file, err := os.Open(filepath.Clean(cfgPath))
		if err != nil {
			log.Fatal("config file problem", err)
		}
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		configG = &Config{}
		err = decoder.Decode(configG)
		if err != nil {
			log.Fatal("config file problem", err)
		}
	})
	return configG
}
