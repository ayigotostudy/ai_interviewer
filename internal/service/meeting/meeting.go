package meetingService

import (
	"ai_jianli_go/internal/dao"
	"ai_jianli_go/logs"
	"ai_jianli_go/pkg/rag"
	"ai_jianli_go/types/model"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

type MeetingService struct {
	dao *dao.MeetingDAO
}

func NewMeetingService(dao *dao.MeetingDAO) *MeetingService {
	return &MeetingService{dao: dao}
}

func (s *MeetingService) Create(request *req.CreateMeetingReq) int64 {
	meeting := &model.Meeting{
		UserID:         request.UserID,
		Candidate:      request.Candidate,
		Position:       request.Position,
		JobDescription: request.JobDescription,
		Time:           request.Time,
		Status:         request.Status,
		Remark:         request.Remark,
	}
	err := s.dao.Create(meeting)
	if err != nil {
		logs.SugarLogger.Errorf("创建面试记录失败: %v", err)
		return common.CodeCreateMeetingFail
	}
	return common.CodeSuccess
}

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
	}
	if request.Remark != "" {
		meeting.Remark = request.Remark
	}

	err = s.dao.Update(meeting)
	if err != nil {
		logs.SugarLogger.Errorf("更新面试记录失败: %v", err)
		return common.CodeUpdateMeetingFail
	}
	return common.CodeSuccess
}

func (s *MeetingService) Get(id uint) (*model.Meeting, int64) {
	meeting, err := s.dao.GetByID(id)
	if err != nil {
		logs.SugarLogger.Errorf("获取面试记录失败: %v", err)
		return nil, common.CodeMeetingNotExist
	}
	return meeting, common.CodeSuccess
}

func (s *MeetingService) List() ([]model.Meeting, int64) {
	meetings, err := s.dao.List()
	if err != nil {
		logs.SugarLogger.Errorf("获取面试列表失败: %v", err)
		return nil, common.CodeServerBusy
	}
	return meetings, common.CodeSuccess
}

func (s *MeetingService) Delete(id uint) int64 {
	err := s.dao.Delete(id)
	if err != nil {
		logs.SugarLogger.Errorf("删除面试记录失败: %v", err)
		return common.CodeServerBusy
	}
	return common.CodeSuccess
}

func (s *MeetingService) UploadResume(request *req.UploadResumeReq) int64 {
	err := s.dao.UploadResume(request.MeetingID, request.Resume)
	if err != nil {
		logs.SugarLogger.Errorf("上传简历失败: %v", err)
		return common.CodeServerBusy
	}
	return common.CodeSuccess
}

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

	if meeting.Resume == "" {
		logs.SugarLogger.Errorf("面试记录 %d 未上传简历", request.MeetingID)
		return "", common.CodeResumeNotExist
	}

	// 1. 获取简历内容
	resume, code := s.GetResume(request.MeetingID)
	if code != common.CodeSuccess {
		return "", code
	}

	// 2. 获取面试信息
	meeting, code = s.Get(request.MeetingID)
	if code != common.CodeSuccess {
		return "", code
	}

	// 3. 获取历史对话
	ctx := context.Background()
	memory := rag.NewRedisMemory(rag.RedisMemoryConfig{
		MaxWindowSize: 20,
		RedisOptions: &redis.Options{
			Addr:     "124.222.151.35:6379",
			Password: "123456",
		},
	})

	con := memory.GetConversation(fmt.Sprintf("%d", request.MeetingID), false)

	// 检查面试轮数
	if con.GetRoundCount() >= 20 {
		meeting.Status = "已完成"
		if err := s.dao.Update(meeting); err != nil {
			logs.SugarLogger.Errorf("更新面试轮数失败: %v", err)
			return "", common.CodeServerBusy
		}
		return "", common.CodeInterviewRoundLimit
	}

	if con.GetLastConversationsKnowledge() == "" {
		con.SetLastConversationKnowledge(resume)
	}
	// 4. 获取知识库相关内容
	retriever := rag.GetRetriever()
	docs, err := retriever.Retrieve(ctx, con.GetLastConversationsKnowledge())
	if err != nil {
		logs.SugarLogger.Errorf("获取检索器失败: %v", err)
		return "", common.CodeServerBusy
	}

	// 5. 构建上下文
	context := ""
	if len(docs) > 0 {
		contextParts := make([]string, len(docs))
		for i, doc := range docs {
			contextParts[i] = fmt.Sprintf("文档片段[%d]:\n%s\n", i+1, doc.Content)
		}
		context = strings.Join(contextParts, "\n---\n")
	}

	// 6. 创建对话模型
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   "gpt-4o",
		BaseURL: "https://api.vveai.com/v1",
		APIKey:  "sk-Xs5rROO2htFLGMJh407b42F505Fe4c89A8510f7608E52c2f",
	})
	if err != nil {
		logs.SugarLogger.Errorf("模型初始化失败: %v", err)
		return "", common.CodeServerBusy
	}

	// 7. 创建提示模板
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

	// 8. 构建提示
	prompt := map[string]any{
		"context":         context,
		"answer":          request.Answer,
		"resume":          resume,
		"history":         con.String(),
		"job_description": meeting.JobDescription,
	}

	// 9. 生成回答
	messages, err := template.Format(ctx, prompt)
	if err != nil {
		logs.SugarLogger.Errorf("生成回答失败: %v", err)
		return "", common.CodeServerBusy
	}

	con.Append(schema.UserMessage(request.Answer))
	res, err := chatModel.Generate(ctx, messages)
	if err != nil {
		return "", common.CodeServerBusy
	}

	// 10. 提取知识点并更新对话
	knowledgePoint := extractKnowledgePoint(res.Content)
	con.SetLastConversationKnowledge(knowledgePoint)
	con.Append(res)

	// 11. 更新面试记录和总结
	meeting.InterviewRecord = con.String()
	meeting.InterviewSummary = res.Content

	// 如果达到最大轮数，更新面试状态为已完成
	if con.GetRoundCount() >= 20 {
		meeting.Status = "已完成"
	}

	if err := s.dao.Update(meeting); err != nil {
		logs.SugarLogger.Errorf("更新面试记录失败: %v", err)
		return "", common.CodeServerBusy
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
