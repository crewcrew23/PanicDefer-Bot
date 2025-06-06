package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BotToken            string   `yaml:"bot_token" env-required:"true"`
	MqConfig            MQConfig `yaml:"rabbitMQ" env-required:"true"`
	Env                 string   `yaml:"env" env-required:"true"`
	DbPath              string   `yaml:"db" env-required:"true"`
	TimeToPing          int      `yaml:"time_to_ping" env-required:"true"`
	PingTTLStore        int      `yaml:"ping_ttl_store" env-required:"true"`
	AbnormalCoefficient float64  `yaml:"abnormal_coefficient" env-required:"true"`
	Worker              Worker   `yaml:"worker" env-required:"true"`
}

type MQConfig struct {
	Host   string     `yaml:"host" env-required:"true"`
	Topics TopicsConf `yaml:"topic" env-required:"true"`
}

type TopicsConf struct {
	FromServerTopic string `yaml:"fromServerName" env-required:"true"`
	FromWorkerTopic string `yaml:"fromWorkerName" env-required:"true"`
}

type Worker struct {
	PingWorker    int `yaml:"pingWorker" env-required:"true"`
	HistoryWorker int `yaml:"historyWorker" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfiPath()
	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist:" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Fail to parse config")
	}

	return &cfg
}

func fetchConfiPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config")
	flag.Parse()

	if res == "" {
		panic("flag config is required")
	}

	return res
}
