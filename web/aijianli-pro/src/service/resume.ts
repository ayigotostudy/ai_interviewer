import request from '@/plugins/axiosInstance'

// 创建简历参数接口
export interface CreateResumeParams {
  user_id: number
  name: string
  basic_info: string
  work_exp: string
  project_exp: string
  self_eval: string
  awards: string
  target_job: string
  template_id: number
}

// 简历详情响应接口
export interface ResumeDetailResponse {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  user_id: number
  name: string
  content: string  // Markdown 格式的简历内容
  template_id: number
  status: number
}

// 简历列表项接口
export interface ResumeListItem {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  user_id: number
  name: string
  content: string  // Markdown 格式的简历内容
  template_id: number
  status: number
}

// 简历模板接口
export interface ResumeTemplate {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: null | string
  name: string
  show_content: string
}

// 创建简历
export const createResume = (data: CreateResumeParams) => {
  return request.post('/resume', data)
}

// 获取简历详情
export const getResumeDetail = (id: string) => {
  return request.get('/resume', {
    params: { id }
  })
}

// 获取简历列表
export const getResumeList = () => {
  return request.get('/resume/list')
}

// 获取简历模板
export const getResumeTemplates = () => {
  return request.get('/resume/template')
}

// 删除简历
export const deleteResume = (id: number) => {
  return request.delete('/resume', {
    data: { id }
  })
}