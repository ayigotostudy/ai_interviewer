package req

type CreateMeetingReq struct {
	UserID         uint   `json:"user_id"`                      // 用户ID
	Candidate      string `json:"candidate" binding:"required"` // 候选人
	Position       string `json:"position" binding:"required"`  // 职位
	JobDescription string `json:"job_description"`              // 职位描述
	Time           int64  `json:"time"`                         // 面试时间
	Status         string `json:"status"`                       // 面试状态
	Remark         string `json:"remark"`                       // 备注
	WikiID          uint   `json:"wiki_id"`                      // 知识库ID
}

type UpdateMeetingReq struct {
	ID               uint   `json:"id" binding:"required"` // 面试ID
	UserID           uint   `json:"user_id"`               // 用户ID
	Candidate        string `json:"candidate"`             // 候选人
	Position         string `json:"position"`              // 职位
	JobDescription   string `json:"job_description"`       // 职位描述
	Time             int64  `json:"time"`                  // 面试时间
	Status           string `json:"status"`                // 面试状态
	Remark           string `json:"remark"`                // 备注
	InterviewRecord  string `json:"interview_record"`      // 面试记录
	InterviewSummary string `json:"interview_summary"`     // 面试总结
}

type GetMeetingReq struct {
	ID uint `json:"id" binding:"required"`
}

type UploadResumeReq struct {
	UserID    uint   `json:"user_id"`                       // 用户ID
	MeetingID uint   `json:"meeting_id" binding:"required"` // 面试ID
	Resume    string `json:"resume" binding:"required"`     // 简历内容
}

type AIInterviewReq struct {
	UserID    uint   `json:"user_id"`                       // 用户ID
	MeetingID uint   `json:"meeting_id" binding:"required"` // 面试ID
	Answer    string `json:"answer" binding:"required"`     // 应聘者回答
}

type GetRemarkReq struct {
	UserID    uint   `json:"user_id"`                       // 用户ID
	MeetingID uint   `json:"meeting_id" binding:"required"` // 面试ID
}