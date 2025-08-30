package req

// CreateResumeRequest 创建简历请求参数
type CreateResumeRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`      // 用户ID
	Name       string `json:"name" binding:"required,min=1,max=50"` // 简历名称
	BasicInfo  string `json:"basic_info" binding:"required"`  // 基本信息（姓名、年龄、学历、联系方式等）
	WorkExp    string `json:"work_exp"`                        // 工作经历（公司名称、职位、工作时间、工作内容等）
	ProjectExp string `json:"project_exp"`                     // 项目经历（项目名称、角色、时间、项目描述等）
	SelfEval   string `json:"self_eval"`                       // 个人评价
	Awards     string `json:"awards"`                          // 获奖情况
	TargetJob  string `json:"target_job"`                      // 目标岗位信息
	TemplateID int    `json:"template_id" binding:"required"` // 模板ID
}

// UpdateResumeRequest 更新简历请求参数
type UpdateResumeRequest struct {
	ID         int64   `json:"id" binding:"required"`          // 简历ID
	Content    string `json:"content"`                        // 简历内容
}

// GetResumeListRequest 获取简历列表请求参数
type GetResumeListRequest struct {
	UserID uint `json:"user_id" binding:"required"` // 用户ID
}

// GetResumeDetailRequest 获取简历详情请求参数
type GetResumeDetailRequest struct {
	ID uint `json:"id" binding:"required"` // 简历ID
}

// DeleteResumeRequest 删除简历请求参数
type DeleteResumeRequest struct {
	ID uint `json:"id" binding:"required"` // 简历ID
}