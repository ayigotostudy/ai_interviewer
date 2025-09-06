package config

import (
	"ai_jianli_go/logs"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL `yaml:"mysql"`
	Redis `yaml:"redis"`
	//EmailInfo `yaml:"email"`
	Speech    `yaml:"speech"`
	LocalPath `yaml:"localPath"`
	Role      `yaml:"role"`
	RateLimit `yaml:"rateLimit"`
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

type LocalPath struct {
	Path string `yaml:"path"`
}

type Role struct {
	Model  string `yaml:"model"`
	Policy string `yaml:"policy"`
}

// RateLimit 限流配置结构体
type RateLimit struct {
	// 是否启用限流
	Enabled bool `yaml:"enabled"`

	// 语音识别接口配置
	Speech SpeechRateLimit `yaml:"speech"`

	// 通用API配置
	General GeneralRateLimit `yaml:"general"`

	// 文件上传配置
	Upload UploadRateLimit `yaml:"upload"`

	// 认证接口配置
	Auth AuthRateLimit `yaml:"auth"`
}

// RoleRateLimit 角色特定限流配置
type RoleRateLimit struct {
	Rate  int `yaml:"rate"`  // 该角色的QPS限制
	Burst int `yaml:"burst"` // 该角色的桶容量
}

// SpeechRateLimit 语音识别接口限流配置
type SpeechRateLimit struct {
	DefaultRate  int                      `yaml:"defaultRate"`  // 默认QPS限制
	DefaultBurst int                      `yaml:"defaultBurst"` // 默认桶容量
	RoleLimits   map[string]RoleRateLimit `yaml:"roleLimits"`   // 角色特定限流配置
	Enabled      bool                     `yaml:"enabled"`      // 是否启用
	SkipRoles    []string                 `yaml:"skipRoles"`    // 跳过限流的角色
}

// GeneralRateLimit 通用API限流配置
type GeneralRateLimit struct {
	DefaultRate  int                      `yaml:"defaultRate"`
	DefaultBurst int                      `yaml:"defaultBurst"`
	RoleLimits   map[string]RoleRateLimit `yaml:"roleLimits"`
	Enabled      bool                     `yaml:"enabled"`
	SkipRoles    []string                 `yaml:"skipRoles"`
}

// UploadRateLimit 文件上传限流配置
type UploadRateLimit struct {
	DefaultRate  int                      `yaml:"defaultRate"`
	DefaultBurst int                      `yaml:"defaultBurst"`
	RoleLimits   map[string]RoleRateLimit `yaml:"roleLimits"`
	Enabled      bool                     `yaml:"enabled"`
	SkipRoles    []string                 `yaml:"skipRoles"`
}

// AuthRateLimit 认证接口限流配置
type AuthRateLimit struct {
	DefaultRate  int                      `yaml:"defaultRate"`
	DefaultBurst int                      `yaml:"defaultBurst"`
	RoleLimits   map[string]RoleRateLimit `yaml:"roleLimits"`
	Enabled      bool                     `yaml:"enabled"`
	SkipRoles    []string                 `yaml:"skipRoles"`
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

	viper.SetConfigName("ai_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workdir + "/config")

	if err = viper.Unmarshal(&aiConfig); err != nil {
		panic(err)
	}

	logs.SugarLogger.Infof("config: %v", config)
	InitAIConfig()
}

func GetMySQLConfig() MySQL {
	return config.MySQL
}

func GetRedisConfig() Redis {
	return config.Redis
}

func GetLocalPathConfig() LocalPath {
	return config.LocalPath
}

func GetRoleConfig() Role {
	return config.Role
}

// func GetEmailInfo() EmailInfo {
// 	return config.EmailInfo
// }

func GetSpeechConfig() Speech {
	return config.Speech
}

func GetRateLimitConfig() RateLimit {
	return config.RateLimit
}
