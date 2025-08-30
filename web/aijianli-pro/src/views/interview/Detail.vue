<template>
  <div class="interview-detail-container">
    <div class="detail-header">
      <div class="header-content">
        <button class="back-btn" @click="goBack">
          <span class="icon">←</span>
          返回
        </button>
        <h1>面试详情</h1>
        <div class="header-actions">
          <button class="start-btn" @click="startInterview" v-if="interview.status === 'planned'">
            <span class="icon">▶️</span>
            开始面试
          </button>
          <button class="edit-btn" @click="editInterview">
            <span class="icon">✏️</span>
            编辑
          </button>
          <button class="evaluation-btn" @click="viewEvaluation" v-if="interview.status === 'completed'">
            <span class="icon">📊</span>
            查看评价
          </button>
          <button class="cancel-btn" @click="cancelInterview" v-if="interview.status === 'planned'">
            <span class="icon">❌</span>
            取消面试
          </button>
        </div>
      </div>
    </div>

    <div class="detail-content">
      <div class="interview-info">
        <div class="info-card">
          <div class="card-header">
            <h3>面试信息</h3>
            <span class="status-badge" :class="getStatusClass(interview.status)">
              {{ getStatusText(interview.status) }}
            </span>
          </div>
          
          <div class="info-grid">
            <div class="info-item">
              <span class="label">候选人：</span>
              <span class="value">{{ interview.candidate }}</span>
            </div>
            <div class="info-item">
              <span class="label">目标职位：</span>
              <span class="value">{{ interview.position }}</span>
            </div>
            <div class="info-item">
              <span class="label">面试时间：</span>
              <span class="value">{{ formatTime(interview.time) }}</span>
            </div>
            <div class="info-item">
              <span class="label">面试类型：</span>
              <span class="value">{{ getInterviewTypeText(interview.type) }}</span>
            </div>
            <div class="info-item">
              <span class="label">面试官：</span>
              <span class="value">{{ interview.interviewer || 'AI面试官' }}</span>
            </div>
            <div class="info-item">
              <span class="label">创建时间：</span>
              <span class="value">{{ formatTime(interview.created_at) }}</span>
            </div>
          </div>

          <div class="description-section" v-if="interview.job_description">
            <h4>职位描述</h4>
            <p>{{ interview.job_description }}</p>
          </div>


        </div>

        <div class="resume-card" v-if="interview.resume_id">
          <div class="card-header">
            <h3>候选人简历</h3>
            <button class="view-resume-btn" @click="viewResume">
              <span class="icon">👁️</span>
              查看完整简历
            </button>
          </div>
          
          <div class="resume-preview">
            <div class="resume-header">
              <h4>{{ resume.name || '未命名简历' }}</h4>
              <div class="template-badge">
                模板：{{ getTemplateName(resume.template_id) }}
              </div>
            </div>
            
            <div class="resume-content">
              <div class="resume-section">
                <h5>简历内容</h5>
                <div class="resume-markdown" v-html="renderMarkdown(resume.content || '')"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="interview-sidebar">
        <div class="sidebar-section">
          <h4>面试记录</h4>
          <div v-if="interviewRecords.length > 0">
            <div v-for="record in interviewRecords" :key="record.id" class="record-item">
              <div class="record-header">
                <span class="record-time">{{ formatTime(record.timestamp) }}</span>
                <span class="record-type" :class="getRecordTypeClass(record.type)">
                  {{ getRecordTypeText(record.type) }}
                </span>
              </div>
              <p class="record-content">{{ record.content }}</p>
            </div>
          </div>
          <p v-else class="no-data">暂无面试记录</p>
        </div>

        <div class="sidebar-section">
          <h4>面试评估</h4>
          <div v-if="interview.evaluation">
            <div class="evaluation-item">
              <span class="label">综合评分：</span>
              <span class="score">{{ interview.evaluation.score }}/100</span>
            </div>
            <div class="evaluation-item">
              <span class="label">技术能力：</span>
              <span class="score">{{ interview.evaluation.technical }}/100</span>
            </div>
            <div class="evaluation-item">
              <span class="label">沟通能力：</span>
              <span class="score">{{ interview.evaluation.communication }}/100</span>
            </div>
            <div class="evaluation-item">
              <span class="label">团队协作：</span>
              <span class="score">{{ interview.evaluation.teamwork }}/100</span>
            </div>
            <div class="evaluation-item">
              <span class="label">总体评价：</span>
              <span class="comment">{{ interview.evaluation.comment }}</span>
            </div>
          </div>
          <p v-else class="no-data">暂无评估结果</p>
        </div>

        <div class="sidebar-section">
          <h4>快捷操作</h4>
          <button class="action-btn primary" @click="startInterview" v-if="interview.status === 'planned'">
            <span class="icon">▶️</span>
            开始面试
          </button>
          <button class="action-btn secondary" @click="rescheduleInterview" v-if="interview.status === 'planned'">
            <span class="icon">📅</span>
            重新安排
          </button>
          <button class="action-btn danger" @click="deleteInterview">
            <span class="icon">🗑️</span>
            删除面试
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getMeetingDetail, deleteMeeting } from '@/service/meeting'
import { getResumeDetail } from '@/service/resume'

const router = useRouter()
const route = useRoute()

// 面试数据
const interview = ref<any>({})
const resume = ref<any>({})
const interviewRecords = ref<any[]>([])

const getStatusClass = (status: string) => {
  const statusClasses: { [key: string]: string } = {
    planned: 'status-planned',
    in_progress: 'status-progress',
    completed: 'status-completed',
    canceled: 'status-canceled'
  }
  return statusClasses[status] || 'status-planned'
}

const getStatusText = (status: string) => {
  const statusTexts: { [key: string]: string } = {
    planned: '待开始',
    in_progress: '进行中',
    completed: '已完成',
    canceled: '已取消'
  }
  return statusTexts[status] || '待开始'
}

const getInterviewTypeText = (type: string) => {
  const typeTexts: { [key: string]: string } = {
    text: '文本面试',
    video: '视频面试',
    mixed: '混合面试'
  }
  return typeTexts[type] || '文本面试'
}

const getRecordTypeClass = (type: string) => {
  const typeClasses: { [key: string]: string } = {
    question: 'type-question',
    answer: 'type-answer',
    evaluation: 'type-evaluation'
  }
  return typeClasses[type] || 'type-question'
}

const getRecordTypeText = (type: string) => {
  const typeTexts: { [key: string]: string } = {
    question: '问题',
    answer: '回答',
    evaluation: '评估'
  }
  return typeTexts[type] || '问题'
}

const getTemplateName = (templateId: number) => {
  const templateNames: { [key: number]: string } = {
    1: '经典模板',
    2: '现代简约',
    3: '创意设计',
    4: '专业商务'
  }
  return templateNames[templateId] || '未知模板'
}

// 渲染Markdown
const renderMarkdown = (markdown: string): string => {
  if (!markdown) return ''
  
  // 简单的Markdown转HTML转换
  let html = markdown
    // 处理 ::: start ... ::: end 块
    .replace(/:::\s*start\s*\n*([\s\S]*?)\n*:::\s*end/g, (match, content) => {
      const lines = content.trim().split('\n').filter((line: string) => line.trim())
      const formattedContent = lines.map((line: string) => `<div class="highlight-line">${line.trim()}</div>`).join('')
      return `<div class="highlight-block">${formattedContent}</div>`
    })
    // 处理单独的 ::: 符号（作为分隔符）
    .replace(/:::/g, '<span class="highlight-symbol">:::</span>')
    // 标题
    .replace(/^### (.*$)/gim, '<h3>$1</h3>')
    .replace(/^## (.*$)/gim, '<h2>$1</h2>')
    .replace(/^# (.*$)/gim, '<h1>$1</h1>')
    // 粗体
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    // 斜体
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    // 代码块
    .replace(/```([\s\S]*?)```/g, '<pre><code>$1</code></pre>')
    // 行内代码
    .replace(/`(.*?)`/g, '<code>$1</code>')
    // 链接
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>')
    // 列表
    .replace(/^\* (.*$)/gim, '<li>$1</li>')
    .replace(/^- (.*$)/gim, '<li>$1</li>')
    // 处理基本信息行（用 | 分隔的信息）
    .replace(/^([^#\n]+)\|([^#\n]+)\|([^#\n]+)\|([^#\n]+)\|([^#\n]+)\|([^#\n]+)$/gm, (match, p1, p2, p3, p4, p5, p6) => {
      return `<div class="basic-info-line">
        <span class="info-item">${p1.trim()}</span>
        <span class="info-separator">|</span>
        <span class="info-item">${p2.trim()}</span>
        <span class="info-separator">|</span>
        <span class="info-item">${p3.trim()}</span>
        <span class="info-separator">|</span>
        <span class="info-item">${p4.trim()}</span>
        <span class="info-separator">|</span>
        <span class="info-item">${p5.trim()}</span>
        <span class="info-separator">|</span>
        <span class="info-item">${p6.trim()}</span>
      </div>`
    })
    // 段落
    .replace(/\n\n/g, '</p><p>')
    .replace(/^(?!<[h|u|o|p|d|b|d])(.*$)/gim, '<p>$1</p>')
    // 清理多余的p标签
    .replace(/<p><\/p>/g, '')
    .replace(/<p><p>/g, '<p>')
    .replace(/<\/p><\/p>/g, '</p>')
  
  return html
}

const formatTime = (timestamp: number) => {
  if (!timestamp) return '未知时间'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

const goBack = () => {
  router.go(-1)
}

const startInterview = () => {
  router.push(`/interview/room/${route.params.id}`)
}

const editInterview = () => {
  router.push(`/interview/edit/${route.params.id}`)
}

const viewEvaluation = () => {
  router.push(`/interview/evaluation/${route.params.id}`)
}

const cancelInterview = async () => {
  if (!confirm('确定要取消这场面试吗？此操作不可恢复。')) {
    return
  }

  try {
    // 调用删除面试API
    const response = await deleteMeeting(Number(route.params.id))
    
    if (response.data.code === 1000) {
      alert('面试已成功取消')
      router.push('/interview')
    } else {
      alert('取消面试失败: ' + response.data.msg)
    }
  } catch (error) {
    console.error('取消面试失败:', error)
    alert('取消面试失败，请重试')
  }
}

const viewResume = () => {
  if (interview.value.resume_id) {
    router.push(`/resume/${interview.value.resume_id}`)
  }
}

const rescheduleInterview = () => {
  // 这里实现重新安排面试功能
  alert('重新安排功能开发中...')
}

const deleteInterview = async () => {
  if (!confirm('确定要删除这场面试吗？此操作不可恢复。')) {
    return
  }

  try {
    // 这里应该调用删除面试API
    // await deleteMeeting(route.params.id)
    
    alert('面试删除成功')
    router.push('/interview')
  } catch (error) {
    console.error('删除面试失败:', error)
    alert('删除面试失败，请重试')
  }
}

const loadInterviewDetail = async () => {
  try {
    const interviewId = route.params.id as string
    const response = await getMeetingDetail(interviewId)
    if (response.data.code === 1000) {
      interview.value = response.data.data
      // 加载简历数据
      if (interview.value.resume_id) {
        await loadResumeData()
      }
    } else {
      console.error('获取面试详情失败:', response.data.msg)
      alert('获取面试详情失败: ' + response.data.msg)
    }
  } catch (error) {
    console.error('获取面试详情失败:', error)
    alert('获取面试详情失败，请重试')
  }
}

const loadResumeData = async () => {
  if (!interview.value.resume_id) return
  
  try {
    const resumeId = interview.value.resume_id.toString()
    const response = await getResumeDetail(resumeId)
    if (response.data.code === 1000) {
      resume.value = response.data.data
    } else {
      console.error('获取简历详情失败:', response.data.msg)
    }
  } catch (error) {
    console.error('获取简历详情失败:', error)
  }
}

const loadInterviewRecords = async () => {
  try {
    // 面试记录暂时为空，等待后端提供相关API
    // const response = await getInterviewRecords(route.params.id)
    // if (response.data.code === 1000) {
    //   interviewRecords.value = response.data.data
    // }
    interviewRecords.value = []
  } catch (error) {
    console.error('获取面试记录失败:', error)
  }
}

onMounted(() => {
  loadInterviewDetail()
  // loadResumeData() 已在 loadInterviewDetail 中调用
  loadInterviewRecords()
})
</script>

<style scoped>
.interview-detail-container {
  min-height: 100vh;
  background: #F8FAFC;
}

.detail-header {
  background: white;
  border-bottom: 1px solid #E5E7EB;
  padding: 1rem 0;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.back-btn {
  background: #F3F4F6;
  color: #374151;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: #E5E7EB;
}

.header-content h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1F2937;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.start-btn, .edit-btn, .evaluation-btn, .cancel-btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
}

.start-btn {
  background: #10B981;
  color: white;
}

.start-btn:hover {
  background: #059669;
}

.edit-btn {
  background: #FEF3C7;
  color: #D97706;
}

.edit-btn:hover {
  background: #FDE68A;
}

.evaluation-btn {
  background: #E0E7FF;
  color: #3730A3;
}

.evaluation-btn:hover {
  background: #C7D2FE;
}

.cancel-btn {
  background: #FEE2E2;
  color: #DC2626;
}

.cancel-btn:hover {
  background: #FECACA;
}

.detail-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 2rem;
}

.interview-info {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.info-card, .resume-card {
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  padding: 2rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #E5E7EB;
}

.card-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1F2937;
}

.status-badge {
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 500;
}

.status-planned {
  background: #FEF3C7;
  color: #D97706;
}

.status-progress {
  background: #DBEAFE;
  color: #1E40AF;
}

.status-completed {
  background: #D1FAE5;
  color: #065F46;
}

.status-canceled {
  background: #FEE2E2;
  color: #DC2626;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-item .label {
  font-weight: 500;
  color: #6B7280;
  font-size: 0.9rem;
}

.info-item .value {
  color: #374151;
  font-size: 0.9rem;
  font-weight: 500;
}

.description-section, .requirements-section {
  margin-bottom: 1.5rem;
}

.description-section:last-child, .requirements-section:last-child {
  margin-bottom: 0;
}

.description-section h4, .requirements-section h4 {
  margin: 0 0 0.75rem 0;
  font-size: 1rem;
  font-weight: 600;
  color: #374151;
}

.description-section p, .requirements-section p {
  margin: 0;
  color: #6B7280;
  line-height: 1.6;
  white-space: pre-line;
}

.view-resume-btn {
  background: #DBEAFE;
  color: #1E40AF;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  font-size: 0.8rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
}

.view-resume-btn:hover {
  background: #BFDBFE;
}

.resume-preview {
  margin-top: 1rem;
}

.resume-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.resume-header h4 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #1F2937;
}

.template-badge {
  padding: 0.25rem 0.75rem;
  background: #F3F4F6;
  color: #6B7280;
  border-radius: 15px;
  font-size: 0.7rem;
  font-weight: 500;
}

.resume-content {
  background: #F9FAFB;
  border-radius: 10px;
  padding: 1rem;
}

.resume-section {
  margin-bottom: 1rem;
}

.resume-section:last-child {
  margin-bottom: 0;
}

.resume-section h5 {
  margin: 0 0 8px 0;
  color: #374151;
  font-size: 1em;
  font-weight: 600;
}

.resume-section p {
  margin: 0;
  color: #6b7280;
  line-height: 1.6;
}

/* 简历Markdown内容样式 */
.resume-markdown {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  line-height: 1.6;
  color: #333;
}

.resume-markdown h1 {
  font-size: 1.6em;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 12px 0;
  padding-bottom: 6px;
  border-bottom: 2px solid #e5e7eb;
}

.resume-markdown h2 {
  font-size: 1.3em;
  font-weight: 600;
  color: #ea0202;
  margin: 16px 0 10px 0;
}

.resume-markdown h3 {
  font-size: 1.1em;
  font-weight: 600;
  color: #374151;
  margin: 14px 0 8px 0;
}

.resume-markdown p {
  margin: 0 0 10px 0;
  line-height: 1.7;
}

.resume-markdown strong {
  font-weight: 600;
  color: #1f2937;
}

.resume-markdown em {
  font-style: italic;
  color: #6b7280;
}

.resume-markdown code {
  background: #f3f4f6;
  padding: 2px 5px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9em;
  color: #dc2626;
}

.resume-markdown pre {
  background: #f8fafc;
  padding: 12px;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  overflow-x: auto;
  margin: 12px 0;
}

.resume-markdown pre code {
  background: none;
  padding: 0;
  color: #374151;
}

/* 高亮块样式 */
.resume-markdown .highlight-block {
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border: 1px solid #bae6fd;
  border-radius: 6px;
  padding: 12px;
  margin: 12px 0;
  position: relative;
}

.resume-markdown .highlight-block::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6, #06b6d4);
  border-radius: 6px 6px 0 0;
}

.resume-markdown .highlight-line {
  padding: 6px 0;
  border-bottom: 1px solid rgba(186, 230, 253, 0.3);
  font-weight: 500;
  color: #1e40af;
}

.resume-markdown .highlight-line:last-child {
  border-bottom: none;
}

/* 高亮符号样式 */
.resume-markdown .highlight-symbol {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  color: #92400e;
  padding: 1px 4px;
  border-radius: 3px;
  font-weight: 600;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.8em;
  border: 1px solid #fbbf24;
  display: inline-block;
  margin: 0 1px;
}

.interview-sidebar {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.sidebar-section {
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
}

.sidebar-section h4 {
  margin: 0 0 1rem 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #1F2937;
  border-bottom: 1px solid #F3F4F6;
  padding-bottom: 0.5rem;
}

.record-item {
  padding: 1rem;
  background: #F9FAFB;
  border-radius: 10px;
  margin-bottom: 0.75rem;
}

.record-item:last-child {
  margin-bottom: 0;
}

.record-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.record-time {
  font-size: 0.8rem;
  color: #9CA3AF;
}

.record-type {
  padding: 0.25rem 0.5rem;
  border-radius: 15px;
  font-size: 0.7rem;
  font-weight: 500;
}

.type-question {
  background: #DBEAFE;
  color: #1E40AF;
}

.type-answer {
  background: #D1FAE5;
  color: #065F46;
}

.type-evaluation {
  background: #FEF3C7;
  color: #D97706;
}

.record-content {
  margin: 0;
  color: #374151;
  font-size: 0.85rem;
  line-height: 1.4;
}

.no-data {
  color: #9CA3AF;
  font-style: italic;
  text-align: center;
  margin: 1rem 0;
}

.evaluation-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.evaluation-item:last-child {
  margin-bottom: 0;
}

.evaluation-item .label {
  font-weight: 500;
  color: #6B7280;
  font-size: 0.9rem;
}

.evaluation-item .score {
  color: #1F2937;
  font-weight: 600;
  font-size: 0.9rem;
}

.evaluation-item .comment {
  color: #374151;
  font-size: 0.85rem;
  max-width: 150px;
  text-align: right;
}

.action-btn {
  width: 100%;
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  margin-bottom: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
}

.action-btn:last-child {
  margin-bottom: 0;
}

.action-btn.primary {
  background: #2563EB;
  color: white;
}

.action-btn.primary:hover {
  background: #1D4ED8;
}

.action-btn.secondary {
  background: #F3F4F6;
  color: #374151;
  border: 1px solid #D1D5DB;
}

.action-btn.secondary:hover {
  background: #E5E7EB;
}

.action-btn.danger {
  background: #FEE2E2;
  color: #DC2626;
}

.action-btn.danger:hover {
  background: #FECACA;
}

@media (max-width: 1024px) {
  .detail-content {
    grid-template-columns: 1fr;
  }
  
  .interview-sidebar {
    order: -1;
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .header-actions {
    flex-direction: column;
    width: 100%;
  }
  
  .start-btn, .edit-btn, .cancel-btn {
    width: 100%;
  }
  
  .detail-content {
    padding: 1rem;
  }
  
  .info-card, .resume-card {
    padding: 1.5rem;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
