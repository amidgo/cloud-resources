package config

import (
	"log"
	"time"

	_ "embed"

	"gopkg.in/yaml.v3"
)

type resourceConfig struct {
	Load     float32 `yaml:"load"`
	MinDelta float32 `yaml:"min_delta"`
	MaxDelta float32 `yaml:"max_delta"`
	Offload  float32 `yaml:"offload"`
	Init     float32 `yaml:"init"`
}

type HealthLoadConfig struct {
	CPU resourceConfig `yaml:"cpu"`
	RAM resourceConfig `yaml:"ram"`
}

type ResourceTypeConfig struct {
	HealthLoad      HealthLoadConfig `yaml:"health_load"`
	MaxMachineCount int              `yaml:"max_machine_count"`
}

// штука чтобы можно было парсить длительность из yaml
type YamlTimeDuration time.Duration

// имплементируем интерфейс yaml.Unmarshaler для парсинга yaml
func (y *YamlTimeDuration) UnmarshalYAML(value *yaml.Node) error {
	dr, err := time.ParseDuration(value.Value)
	*y = YamlTimeDuration(dr)
	return err
}

func (y YamlTimeDuration) String() string {
	return time.Duration(y).String()
}

type SchedulerConfig struct {
	CheckTime YamlTimeDuration `yaml:"check_time"`
}

type LoggerConfig struct {
	LogURL string `yaml:"log_url"`
}

type ApiConfig struct {
	Token string `yaml:"token"`
}

type AppConfig struct {
	Api       ApiConfig          `yaml:"api"`
	VM        ResourceTypeConfig `yaml:"vm"`
	DB        ResourceTypeConfig `yaml:"db"`
	Logger    LoggerConfig       `yaml:"logger"`
	Scheduler SchedulerConfig    `yaml:"scheduler"`
}

// в папке app должен быть файл example.yaml
// вам не обязательно оставлять те настройки которые там поставлены
// фишка в том что можно настраивать сервер с помощью конфигурационного файла

/*
если не вдаваться в душные подробности то этот комментарий парсит файл по пути app/config.yaml
и встраивает его в бинарник, при запуске переменная configData будет проинициализированна со значением которое содержалось в данном файле
в определении могут быть технические неточности так что оставляю ссылку на доку
см. https://pkg.go.dev/embed
*/
//go:embed app/config.yaml
var configData []byte

func ParseAppConfig() *AppConfig {
	var cnf AppConfig
	err := yaml.Unmarshal(configData, &cnf)
	if err != nil {
		log.Fatalf("failed unmarshal config data, err: %s", err)
	}
	return &cnf
}
