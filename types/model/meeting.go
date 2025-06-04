package model

import "gorm.io/gorm"

// 面试表
type Meeting struct {
	gorm.Model
	UserID           uint   `json:"user_id" gorm:"index"` // 用户ID
	Candidate        string `json:"candidate"`            // 候选人
	Position         string `json:"position"`             // 职位
	JobDescription   string `json:"job_description"`      // 职位描述
	Time             int64  `json:"time"`                 // 面试时间
	Status           string `json:"status"`               // 面试状态
	Remark           string `json:"remark"`               // 备注
	Resume           string `json:"resume"`               // 简历内容
	InterviewRecord  string `json:"interview_record"`     // 面试记录
	InterviewSummary string `json:"interview_summary"`    // 面试总结
}
