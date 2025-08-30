package component

import (
	"ai_jianli_go/config"
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

type AIComponent struct {
	ChatModel map[string]*openai.ChatModel
	config    *config.AIConfig
}

func NewAIComponent(aiConfig *config.AIConfig) *AIComponent {
	component := &AIComponent{
		ChatModel: make(map[string]*openai.ChatModel),
		config:    aiConfig,
	}
	component.initChatModel()
	return component
}

var AIComponentInstance *AIComponent

func (ac *AIComponent) GetChatModel(modelName string) *openai.ChatModel {
	return ac.ChatModel[modelName]
}

func (ac *AIComponent) AddChatModel(modelName string, model *openai.ChatModel) {
	ac.ChatModel[modelName] = model
}

func (ac *AIComponent) initChatModel() {
	ctx := context.Background()
	for modelName, modelConfig := range ac.config.Items {
		model, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
			Model:   modelConfig.Model,
			BaseURL: modelConfig.BaseURL,
			APIKey:  modelConfig.APIKey,
		})
		if err != nil {
			panic(err)
		}
		ac.ChatModel[modelName] = model
	}
}

func GetAIComponent() *AIComponent {
	return AIComponentInstance
}

func initAIComponent() {
	AIComponentInstance = NewAIComponent(config.GetAIConfig())
}
