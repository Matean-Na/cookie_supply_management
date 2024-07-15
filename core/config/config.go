package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Server struct {
	Host string
	Port int
}

type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	SslMode  string
}

type Token struct {
	SecretKey string
}

type Dir struct {
	Media    string
	Template string
	Seeder   string
}

type Redis struct {
	Address  string
	Database int
	Password string
}

type Config struct {
	Server
	Database
	Token
	Dir
	Redis
}

var config Config

func Load() (*Config, error) {
	return LoadFromFile(".")
}

func LoadFromFile(path string) (*Config, error) {
	if err := godotenv.Load(path + "/.env"); err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	config = Config{
		Server: Server{
			Host: viper.GetString("SERVER_HOST"),
			Port: viper.GetInt("SERVER_PORT"),
		},
		Database: Database{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			Username: viper.GetString("DB_USERNAME"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SslMode:  viper.GetString("DB_SSLMODE"),
		},
		Token: Token{
			SecretKey: viper.GetString("TOKEN_SECRET_KEY"),
		},
		Dir: Dir{
			Media:    viper.GetString("DIR_MEDIA"),
			Template: viper.GetString("DIR_TEMPLATE"),
			Seeder:   viper.GetString("DIR_SEEDER"),
		},
		Redis: Redis{
			Address:  viper.GetString("REDIS_ADDR"),
			Database: viper.GetInt("REDIS_DB"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
	}

	return &config, nil
}

func Get() *Config {
	return &config
}
