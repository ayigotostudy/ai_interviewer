package resumeService

import (
	"ai_jianli_go/component"
	"ai_jianli_go/internal/dao"
	"ai_jianli_go/logs"
	"ai_jianli_go/types/model"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"
	"context"
	"encoding/json"
	"strings"

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
	logs.SugarLogger.Infof("promptVars: %v\n", promptVars)
// 创建模板，使用 FString 格式
template := prompt.FromMessages(schema.FString,
    // 系统消息模板 - 增强角色定义、输出要求和优化策略
    schema.SystemMessage(`你是一名专业的简历优化与生成专家。你的核心目标是帮助求职者制作一份能高效通过HR筛选和自动化招聘系统（ATS）的专业简历。
    
请严格遵循以下原则生成和优化内容：
1.  **内容扩写与丰富**: 对于用户未提供或简略的内容，基于其现有信息和目标岗位{info}进行合理、专业的扩写，确保内容饱满且真实可信。
2.  **成就量化**: 强烈倾向于使用具体数据、百分比、金额、时间等量化指标来突出成就和贡献。例如，将“负责提升系统性能”优化为“通过代码重构和数据库优化，将系统响应时间减少30%”。
3.  **关键词优化**: 精准分析并融入岗位描述{info}中的关键词和行业特定术语，以提高简历在ATS系统中的匹配度。
4.  **STAR法则**: 在描述工作经历{work}和项目经验{demo}时，运用STAR（情境-Situation, 任务-Task, 行动-Action, 结果-Result）法则来构建内容，使其结构清晰、重点突出。
5.  **专业技能表达**: 优化技能描述，使用“精通”、“熟练掌握”、“熟悉”等程度副词，并结合具体应用场景或技术栈。例如，将“会使用Python”优化为“熟练掌握Python进行数据分析和自动化脚本开发”。
6.  **语言与格式**: 使用专业、主动、有力的动词开头描述经历和成就。保持整体语言简洁、清晰，避免冗余。严格遵循JSON格式输出。
7.  **真实性原则**: 所有扩写和优化内容必须基于用户提供的信息点，不得虚构不存在的重要经历或技能。

请根据提供的简历模板{template}的结构和要求，填充和优化内容，并返回**严格符合模板规范**的JSON数据, 生成数据的字段要与模版需要字段对应。`),

    // 用户消息模板 (参数保持不变，但通过系统提示引导更丰富的输出)
    schema.UserMessage("简历模版: {template}\n基本信息:{basic_info}\n工作经历:{work}\n项目经历:{demo}\n个人评价:{comment}\n获奖情况:{score}\n岗位信息:{info}"),
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
	chatModel := component.GetAIComponent().GetChatModel("gpt-4o")

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

func (s *ResumeService) GetResumeTemplate(ctx context.Context) ([]*model.Template, int64) {
	templates, err := s.dao.GetResumeTemplateList()
	if err != nil {
		logs.SugarLogger.Errorf("获取简历模板失败: %v", err)
		return nil, common.CodeGetResumeTemplateFail
	}
	return templates, common.CodeSuccess
}

func (s *ResumeService) GetResume(ctx context.Context, id uint) (*model.Resume, int64) {
	resume, err := s.dao.GetResume(id)
	if err != nil {
		logs.SugarLogger.Errorf("获取简历失败: %v", err)
		return nil, common.CodeGetResumeFail
	}
	return resume, common.CodeSuccess
}

func (s *ResumeService) GetResumeList(ctx context.Context, userID uint) ([]*model.Resume, int64) {
	resumes, err := s.dao.GetResumeList(ctx, userID)
	if err != nil {
		logs.SugarLogger.Errorf("获取简历列表失败: %v", err)
		return nil, common.CodeGetResumeListFail
	}
	return resumes, common.CodeSuccess
}

func (s *ResumeService) DeleteResume(ctx context.Context, id uint) int64 {
	err := s.dao.DeleteResume(id)
	if err != nil {
		logs.SugarLogger.Errorf("删除简历失败: %v", err)
		return common.CodeDeleteResumeFail
	}
	return common.CodeSuccess
}

func (s *ResumeService) UpdateResume(ctx context.Context, req *req.UpdateResumeRequest) int64 {
	err := s.dao.UpdateResume(req.ID, req.Content)
	if err != nil {
		logs.SugarLogger.Errorf("更新简历失败: %v", err)
		return common.CodeUpdateResumeFail
	}
	return common.CodeSuccess
}