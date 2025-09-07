package meetingService

import (
	"ai_jianli_go/component"
	"ai_jianli_go/internal/dao"
	wikiService "ai_jianli_go/internal/service/wiki"
	"ai_jianli_go/logs"
	"ai_jianli_go/pkg/rag"
	"ai_jianli_go/types/model"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type MeetingService struct {
	dao *dao.MeetingDAO
}

func NewMeetingService(dao *dao.MeetingDAO) *MeetingService {
	return &MeetingService{dao: dao}
}

const (
	PLANED       = "planned"
	INTERVIEWING = "interviewing"
	COMPLETED    = "completed"
	CANCELED     = "canceled"
)

// 创建面试
func (s *MeetingService) Create(request *req.CreateMeetingReq) int64 {
	meeting := &model.Meeting{
		UserID:         request.UserID,
		Candidate:      request.Candidate,
		Position:       request.Position,
		JobDescription: request.JobDescription,
		Time:           request.Time,
		Status:         PLANED,
		WikiID:         request.WikiID,
	}
	err := s.dao.Create(meeting)
	if err != nil {
		logs.SugarLogger.Errorf("创建面试记录失败: %v", err)
		return common.CodeCreateMeetingFail
	}
	return common.CodeSuccess
}

// 更新面试
func (s *MeetingService) Update(request *req.UpdateMeetingReq) int64 {
	meeting, err := s.dao.GetByID(request.ID)
	if err != nil {
		logs.SugarLogger.Errorf("获取面试记录失败: %v", err)
		return common.CodeMeetingNotExist
	}

	if request.UserID != 0 {
		meeting.UserID = request.UserID
	}
	if request.Candidate != "" {
		meeting.Candidate = request.Candidate
	}
	if request.Position != "" {
		meeting.Position = request.Position
	}
	if request.JobDescription != "" {
		meeting.JobDescription = request.JobDescription
	}
	if request.Time != 0 {
		meeting.Time = request.Time
	}
	if request.Status != "" {
		meeting.Status = request.Status
		if request.Status == COMPLETED {
			memory := rag.NewRedisMemory(rag.RedisMemoryConfig{
				RedisOptions:  component.GetRedisDB(),
				MaxWindowSize: 20,
			})
			con := memory.GetConversation(fmt.Sprintf("%d", request.ID), false)
			meeting.InterviewNumber = con.GetRoundCount()
			meeting.InterviewRecord = con.String()
		}
	}

	err = s.dao.Update(meeting)
	if err != nil {
		logs.SugarLogger.Errorf("更新面试记录失败: %v", err)
		return common.CodeUpdateMeetingFail
	}
	return common.CodeSuccess
}

// 获取面试
func (s *MeetingService) Get(id uint) (*model.Meeting, int64) {
	meeting, err := s.dao.GetByID(id)
	if err != nil {
		logs.SugarLogger.Errorf("获取面试记录失败: %v", err)
		return nil, common.CodeMeetingNotExist
	}
	return meeting, common.CodeSuccess
}

// 获取面试列表
func (s *MeetingService) List() ([]model.Meeting, int64) {
	meetings, err := s.dao.List()
	if err != nil {
		logs.SugarLogger.Errorf("获取面试列表失败: %v", err)
		return nil, common.CodeServerBusy
	}
	return meetings, common.CodeSuccess
}

// 删除面试
func (s *MeetingService) Delete(id uint) int64 {
	err := s.dao.Delete(id)
	if err != nil {
		logs.SugarLogger.Errorf("删除面试记录失败: %v", err)
		return common.CodeDeleteMeetingFail
	}
	return common.CodeSuccess
}

// 上传简历
func (s *MeetingService) UploadResume(request *req.UploadResumeReq) int64 {
	err := s.dao.UploadResume(request.MeetingID, request.Resume)
	if err != nil {
		logs.SugarLogger.Errorf("上传简历失败: %v", err)
		return common.CodeUploadResumeFail
	}
	return common.CodeSuccess
}

// 获取简历
func (s *MeetingService) GetResume(meetingID uint) (string, int64) {
	resume, err := s.dao.GetResume(meetingID)
	if err != nil {
		return "", common.CodeResumeNotExist
	}
	return resume, common.CodeSuccess
}

// AI面试主流程
func (s *MeetingService) AIInterview(request *req.AIInterviewReq) (string, int64) {
	meeting, err := s.dao.GetByID(request.MeetingID)
	if err != nil {
		logs.SugarLogger.Errorf("获取面试记录失败: %v", err)
		return "", common.CodeMeetingNotExist
	}

	// 检查面试状态
	if meeting.Status == CANCELED {
		return "", common.CodeMeetingCompleted
	}

	if meeting.Resume == "" {
		return "", common.CodeResumeNotExist
	}

	// 获取历史对话
	ctx := context.Background()
	memory := rag.NewRedisMemory(rag.RedisMemoryConfig{
		MaxWindowSize: 20,
		RedisOptions:  component.GetRedisDB(),
	})

	con := memory.GetConversation(fmt.Sprintf("%d", request.MeetingID), false)

	// 检查面试轮数
	if con.GetRoundCount() >= 20 {
		meeting.Status = COMPLETED
		meeting.InterviewNumber = con.GetRoundCount()
		meeting.InterviewRecord = con.String()
		if err = s.dao.Update(meeting); err != nil {
			logs.SugarLogger.Errorf("更新面试轮数失败: %v", err)
			return "", common.CodeServerBusy
		}
		return "", common.CodeInterviewRoundLimit
	}

	if con.GetLastConversationsKnowledge() == "" {
		con.SetLastConversationKnowledge(meeting.Resume)
	}

	var wiki string
	var code int64
	if meeting.WikiID != 0 {
		wikiService := wikiService.NewWikiService(dao.NewWikiDAO(component.GetMySQLDB()))
		wiki, code = wikiService.Query(&req.QueryWikiRequest{
			UserID: request.UserID,
			RootId: meeting.WikiID,
			Query:  con.GetLastConversationsKnowledge(),
		})
		if code != common.CodeSuccess {
			return "", code
		}
	}else{
		wiki = con.GetLastConversationsKnowledge()
	}

	// 创建对话模型
	chatModel := component.GetAIComponent().GetChatModel("gpt-4o")

	// 创建提示模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(
			"你是一个专业面试官，需要完成以下任务：\n"+
				"1. 基于专业知识库内容提出精准问题\n"+
				"2. 对用户回答进行结构化评价（优点/不足）， 评价后要接着提出问题\n"+
				"3. 针对不足点给出专业解释\n"+
				"4. 根据用户回答生成1-4轮深度追问， 追问结束继续根据简历内容提问\n"+
				"5. 每次回答都需要返回要问的知识点（关键词）以便后续知识库检索提问\n\n"+
				"6. 追问每次只追问一道题目， 后续在根据用户回答继续追问，最多追问4轮\n"+
				"7. 如果用户表示不会， 请不要继续追问， 提问简历的其他知识点\n"+
				"8. 当前是第"+fmt.Sprintf("%d", con.GetRoundCount()+1)+"轮面试，总共20轮\n"+
				"9. 如果用户回答与面试内容无关， 请统一提醒它正在面试（返回知识点继承上次对话的）\n"+
				"当前知识库上下文：{context}\n\n"+
				"当前对话记录：{history}\n\n"+
				"用户简历内容:{resume}\n"+
				"职位描述:{job_description}\n"+
				"输出格式要求：\n"+
				"- 评价使用✅和❌标识优劣点\n",
		),
		schema.UserMessage("【应聘者回答】\n{answer}"),
		schema.AssistantMessage(
			"请按以下结构组织回答：\n"+
				"1. 评价总结（含具体不足点分析）\n"+
				"2. 可追问的知识点:knowledgepoint\n"+
				"3. 问题:question\n",
			[]schema.ToolCall{},
		),
	)

	// 构建提示
	prompt := map[string]any{
		"context":        wiki,
		"answer":          request.Answer,
		"resume":          meeting.Resume,
		"history":         con.String(),
		"job_description": meeting.JobDescription,
	}

	// 9. 生成回答
	messages, err := template.Format(ctx, prompt)
	if err != nil {
		logs.SugarLogger.Errorf("生成回答失败: %v", err)
		return "", common.CodeInterviewGenerateFail
	}

	con.Append(schema.UserMessage(request.Answer))
	res, err := chatModel.Generate(ctx, messages)
	if err != nil {
		logs.SugarLogger.Error(err)
		return "", common.CodeServerBusy
	}

	// 10. 提取知识点并更新对话
	knowledgePoint := extractKnowledgePoint(res.Content)
	con.SetLastConversationKnowledge(knowledgePoint)
	con.Append(res)

	// 如果达到最大轮数，更新面试状态为已完成
	if con.GetRoundCount() >= 20 {
		meeting.Status = COMPLETED
		meeting.InterviewRecord = con.String()
		meeting.InterviewNumber = con.GetRoundCount()
		if err := s.dao.Update(meeting); err != nil {
			logs.SugarLogger.Errorf("更新面试记录失败: %v", err)
			return "", common.CodeServerBusy
		}
	}

	return res.Content, common.CodeSuccess
}

// 提取知识点
func extractKnowledgePoint(input string) string {
	// 查找"可追问的知识点："的位置
	prefix := "可追问的知识点："
	index := strings.Index(input, prefix)

	if index >= 0 {
		// 提取前缀之后的所有内容
		result := input[index+len(prefix):]
		// 去除开头可能的多余空格
		return strings.TrimSpace(result)
	}

	// 兼容可能使用英文冒号的情况
	prefix = "可追问的知识点:"
	index = strings.Index(input, prefix)
	if index >= 0 {
		result := input[index+len(prefix):]
		return strings.TrimSpace(result)
	}

	return "" // 未找到匹配
}

// 获取面试评价
func (s *MeetingService) GetRemark(ctx context.Context, req *req.GetRemarkReq) (string, int64) {
	meeting, err := s.dao.GetByID(req.MeetingID)
	if err != nil {
		logs.SugarLogger.Errorf("获取面试记录失败: %v", err)
		return "", common.CodeGetMeetingFail
	}

	if meeting.Status != COMPLETED {
		return "", common.CodeMeetingNotCompleted
	}

	if meeting.Remark != "" {
		return meeting.Remark, common.CodeSuccess
	}

	model := component.GetAIComponent().GetChatModel("gpt-4o")
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(
			"你是一个专业的面试官，需要根据面试记录以及岗位描述，生成胜任力维度得分、答题内容分析和总体得分, 面试文字评价和可改进点五大板块json数据。\n"+
				"重要要求：\n"+
				"1. 你必须只返回一个纯净的JSON对象，不要有任何额外的前缀、后缀、解释或Markdown代码块标记（如```json）。\n"+
				"2. JSON必须严格遵循我已提供的格式。\n"+
				"3. 不要返回任何非JSON文本。",
		),
		schema.UserMessage("面试记录和总结：\n{input}"),
		schema.AssistantMessage(
			"岗位描述：{job_description}\n返回数据格式：{output}\n",
			[]schema.ToolCall{},
		),
	)
	prompt := map[string]any{
		"input":           meeting.InterviewRecord,
		"job_description": meeting.JobDescription,
		"output":          output,
	}
	messages, err := template.Format(ctx, prompt)
	if err != nil {
		logs.SugarLogger.Errorf("生成提示失败: %v", err)
		return "", common.CodeInterviewGenerateFail
	}

	resp, err := model.Generate(ctx, messages)
	if err != nil {
		logs.SugarLogger.Errorf("生成面试评价失败: %v", err)
		return "", common.CodeInterviewGenerateFail
	}

	meeting.Remark = resp.Content
	if err := s.dao.Update(meeting); err != nil {
		logs.SugarLogger.Errorf("更新面试记录失败: %v", err)
		return "", common.CodeServerBusy
	}

	return resp.Content, common.CodeSuccess
}

const output = `
{
  "overallEvaluation": {
    "score": 85,
    "maxScore": 100,
    "rating": "良好",
    "chartType": "gauge"
  },
  "competencyDimensions": {
    "chartType": "radar",
    "dimensions": [
      {
        "name": "沟通表达",
        "score": 82,
        "fullMark": 100
      },
      {
        "name": "逻辑思维",
        "score": 88,
        "fullMark": 100
      },
      {
        "name": "专业知识",
        "score": 90,
        "fullMark": 100
      },
      {
        "name": "学习能力",
        "score": 85,
        "fullMark": 100
      },
      {
        "name": "抗压性",
        "score": 80,
        "fullMark": 100
      },
      {
        "name": "团队合作",
        "score": 78,
        "fullMark": 100
      }
    ]
  },
  "answerAnalysis": {
    "keywordCloud": {
      "chartType": "wordcloud",
      "keywords": [
        {
          "text": "Spring Boot",
          "value": 32
        },
        {
          "text": "微服务",
          "value": 28
        },
        {
          "text": "MySQL",
          "value": 25
        },
        {
          "text": "分布式系统",
          "value": 22
        },
        {
          "text": "问题解决",
          "value": 19
        }
      ]
    },
    "jdMatch": {
      "chartType": "doughnut",
      "matchPercentage": 76,
      "matchedKeywords": ["Java", "Spring Boot", "MySQL", "云计算"],
      "missingKeywords": ["Redis", "消息队列", "容器化"]
    },
	"interviewEvaluation": "候选人整体表现良好，具备扎实的专业技术功底和清晰的逻辑思维能力。在面试过程中能够围绕Spring Boot、微服务架构和MySQL等核心技术栈展开深入讨论，展现出较强的系统设计能力和问题解决导向。候选人学习能力较好，对新知识保持求知欲，但在高压环境下的稳定性与团队协作意识有进一步提升空间。其技术能力与当前岗位要求有较高匹配度，但在分布式中间件和云原生技术领域存在经验缺口。",
  	"improvablePoints": [
    "团队协作能力有待加强：在跨部门沟通和团队项目协作中表现较为被动，需提升倾听他人意见、整合团队资源的意识与能力[1](@ref)。",
    "抗压性与情绪管理：在压力情境下表现出一定的紧张感，需增强应对复杂问题和紧迫任务的稳定性与韧性[5](@ref)。",
    "技术广度需扩展：缺乏Redis缓存应用、消息队列及容器化技术（如Docker/K8s）的实战经验，需针对性补充分布式系统相关知识[3](@ref)。",
    "表达精炼度不足：技术描述有时过于细节，需提升结构化表达和总结概括能力，增强与非技术人员的沟通效果[1](@ref)。",
    "岗位匹配度提升：虽然基础技能扎实，但仍需弥补JD中明确的‘消息队列’和‘容器化’要求，可通过快速学习或项目实践补足[3](@ref)。"
  	]
}
`
