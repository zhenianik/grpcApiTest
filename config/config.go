package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	config Config
	once   sync.Once
)

type Config struct {
	KafkaAddress string        `mapstructure:"KAFKA_ADDRESS"`
	PostgresUrl  string        `mapstructure:"POSTGRES_URL"`
	GrpcNetwork  string        `mapstructure:"GRPC_NETWORK"`
	GrpcAddress  string        `mapstructure:"GRPC_ADDRESS"`
	RedisHost    string        `mapstructure:"REDIS_HOST"`
	RedisDb      int           `mapstructure:"REDIS_DB"`
	RedisExpires time.Duration `mapstructure:"REDIS_EXPIRES"`
	LogLevel     string        `mapstructure:"LOGGER_LEVEL"`
}

func GetConfig() (*Config, error) {

	var err error

	once.Do(func() {

		viper.AutomaticEnv()

		config = Config{}

		if err = parseYamlConfig(&config, "config.yaml"); err != nil {
			log.Fatal(err)
		}
	})

	return &config, err
}

func parseYamlConfig(cfg interface{}, fileName string) error {
	if len(fileName) == 0 {
		return errors.New("config file not found")
	}

	filePath, err := filepath.Abs(fileName)

	if err != nil {
		return err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("config file '%s' not found", filePath)
	}

	v := viper.New()

	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return v.Unmarshal(&cfg)

}
