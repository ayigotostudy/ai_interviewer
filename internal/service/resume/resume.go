package resumeService

import (
	"ai_jianli_go/internal/dao"
	"ai_jianli_go/logs"
	"ai_jianli_go/types/model"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"
	"context"
	"encoding/json"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"

	"github.com/flosch/pongo2/v5"
)

type ResumeService struct {
	dao *dao.ResumeDAO
}

func NewResumeService(dao *dao.ResumeDAO) *ResumeService {
	return &ResumeService{dao: dao}
}

func (s *ResumeService) CreateResume(ctx context.Context, req *req.CreateResumeRequest)(*model.Resume, int64) {
	resume := model.NewResume(req.UserID, req.Name, req.TemplateID)

	// 根据模板ID获取模板内容
	// 这里假设需要调用模板服务或DAO来获取模板
	// 为简化示例，这里暂时使用模拟数据
	templateModel, err := s.dao.GetTemplate(uint(req.TemplateID))
	if err != nil {
		logs.SugarLogger.Errorf("获取模板失败: %v", err)
		return nil, common.CodeCreateResumeFail
	}

	// 一次性将req参数赋值给prompt变量
	promptVars := map[string]any{
		"role":       "专业简历撰写助手",
		"template":   templateModel.Content,
		"basic_info": req.BasicInfo,
		"work":       req.WorkExp,
		"demo":       req.ProjectExp,
		"comment":    req.SelfEval,
		"score":      req.Awards,
		"info":       req.TargetJob,
	}

	// 创建模板，使用 FString 格式
	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要专业切合实际地回答（最终目标是使简历更容易通过筛选）。你的目标是帮助求职者生成一份专业严谨的简历， 你将根据简历模版（jinja格式）返回对应的json数据， 如果模版有些内容用户没有提供，则帮它生成（可自己扩写）"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("简历模版: {template}\n基本信息:{basic_info}\n工作经历:{work}\n项目简历:{demo}\n个人评价:{comment}\n获奖情况:{score}\n岗位信息:{info}"),
	)

	// 使用模板生成消息
	messages, err := template.Format(context.Background(), promptVars)
	if err != nil {
		logs.SugarLogger.Errorf("生成消息失败: %v", err)
		return nil, common.CodeCreateResumeFail
	}

	// 从配置中获取API密钥，而不是硬编码
	// 这里假设已经有一个配置服务可以获取API密钥
	// apiKey := "" // 实际应用中应该从配置或环境变量中获取

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   "gpt-4o",
		BaseURL: "https://api.vveai.com/v1",
		APIKey:  "sk-npfmWk7VxIyeWYt23c5dCc49E7C343E487913c3e71E30b81",
	})

	if err != nil {
		logs.SugarLogger.Errorf("创建聊天模型失败: %v", err)
		return nil, common.CodeCreateResumeFail
	}

	res, err := chatModel.Generate(ctx, messages)

	if err != nil {
		logs.SugarLogger.Errorf("创建聊天模型失败: %v", err)
		return nil, common.CodeCreateResumeFail
	}

	// 解析生成的内容
	output, err := Write(templateModel.Content, res.Content)
	if err != nil {
		logs.SugarLogger.Errorf("解析生成的内容失败: %v", err)
		return nil, common.CodeCreateResumeFail
	}

	// 保存简历
	resume.Content = output
	err = s.dao.CreateResume(resume)
	if err != nil {
		logs.SugarLogger.Errorf("创建简历失败: %v", err)
		return nil, common.CodeCreateResumeFail
	}

	return resume, common.CodeSuccess
}

// Write 解析生成的内容
func Write(input string, jsonStr string) (string, error) {
	var data map[string]interface{}

	jsonStr, _ = strings.CutPrefix(jsonStr, "```json")
	jsonStr, _ = strings.CutSuffix(jsonStr, "```")

	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", err
	}

	tpl := pongo2.Must(pongo2.FromBytes([]byte(input)))

	output, err := tpl.Execute(pongo2.Context(data))

	if err != nil {
		return "", err
	}

	return output, nil
}
