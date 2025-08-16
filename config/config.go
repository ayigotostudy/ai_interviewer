package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL `yaml:"mysql"`
	Redis `yaml:"redis"`
	//EmailInfo `yaml:"email"`
	Speech `yaml:"speech"`
}

type MySQL struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Pwd    string `yaml:"pwd"`
	DbName string `yaml:"dbname"`
	User   string `yaml:"user"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Pwd  string `yaml:"pwd"`
}

type EmailInfo struct {
	Addr  string `yaml:"addr"`
	Host  string `yaml:"host"`
	From  string `yaml:"from"`
	Email string `yaml:"email"`
	Auth  string `yaml:"auth"`
}

// Speech contains credentials for ASR service
type Speech struct {
	APIKey    string `yaml:"apiKey"`
	APISecret string `yaml:"apiSecret"`
	AppID     string `yaml:"appId"`
}

var config Config

func Init() {
	workdir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workdir + "/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err = viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	fmt.Println(config)
}

func GetMySQLConfig() MySQL {
	return config.MySQL
}

func GetRedisConfig() Redis {
	return config.Redis
}

// func GetEmailInfo() EmailInfo {
// 	return config.EmailInfo
// }

func GetSpeechConfig() Speech {
	return config.Speech
}
