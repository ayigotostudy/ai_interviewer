package config

import (
	"ai_jianli_go/logs"
	"os"

	"github.com/spf13/viper"
)

type AIConfig struct {
	Items map[string]AIConfigItem `yaml:"Items"`
}

type AIConfigItem struct {
	Model   string `yaml:"Model"`
	BaseURL string `yaml:"BaseURL"`
	APIKey  string `yaml:"APIKey"`
}

var aiConfig AIConfig

func GetAIConfig() *AIConfig {
	return &aiConfig
}

// InitAIConfig 初始化AI配置
func InitAIConfig() {
    workdir, _ := os.Getwd()
    viper.SetConfigName("ai_config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(workdir + "/config")

    if err := viper.ReadInConfig(); err != nil {
        logs.SugarLogger.Errorf("Failed to read AI config file: %v", err)
        return
    }

    modelsConfig := make(map[string]AIConfigItem)
    if err := viper.Unmarshal(&modelsConfig); err != nil {
        logs.SugarLogger.Errorf("Failed to unmarshal AI config: %v", err)
        return
    }

    logs.SugarLogger.Infof("AI config loaded successfully: %+v", modelsConfig)

	aiConfig.Items = modelsConfig
}