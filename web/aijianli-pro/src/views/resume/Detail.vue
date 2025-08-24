<template>
  <div class="resume-detail-container">
    <div class="resume-header">
      <div class="header-content">
        <button class="back-btn" @click="goBack">
          <span class="icon">←</span>
          返回
        </button>
        <h1>简历详情</h1>
        <div class="header-actions">
          <button class="edit-btn" @click="editResume">
            <span class="icon">✏️</span>
            编辑
          </button>
          <button class="download-btn" @click="downloadResume">
            <span class="icon">⬇️</span>
            下载
          </button>
        </div>
      </div>
    </div>

    <div class="resume-content">
      <div class="resume-preview">
        <div class="resume-header-preview">
          <h2>{{ resume.name || '未命名简历' }}</h2>
          <div class="template-badge">
            模板：{{ getTemplateName(resume.template_id) }}
          </div>
        </div>

        <div class="resume-content-preview">
          <ResumeShow :content="resume.content" />
        </div>
      </div>

      <div class="resume-sidebar">
        <div class="sidebar-section">
          <h4>简历信息</h4>
          <div class="info-item">
            <span class="label">创建时间：</span>
            <span class="value">{{ formatTime(resume.CreatedAt) }}</span>
          </div>
          <div class="info-item">
            <span class="label">最后更新：</span>
            <span class="value">{{ formatTime(resume.UpdatedAt) }}</span>
          </div>
          <div class="info-item">
            <span class="label">状态：</span>
            <span class="value status-badge" :class="getStatusClass(resume.status)">
              {{ getStatusText(resume.status) }}
            </span>
          </div>
        </div>

        <div class="sidebar-section">
          <h4>快捷操作</h4>
          <button class="action-btn primary" @click="useForInterview">
            <span class="icon">🎯</span>
            用于面试
          </button>
          <button class="action-btn secondary" @click="duplicateResume">
            <span class="icon">📋</span>
            复制简历
          </button>
          <button class="action-btn export" @click="exportToPDF" :disabled="exporting">
            <span v-if="exporting" class="loading-spinner"></span>
            <span v-else class="icon">📄</span>
            {{ exporting ? '导出中...' : '导出PDF' }}
          </button>
          <button class="action-btn danger" 
                  @click="deleteResume"
                  :disabled="isDeleting">
            <span v-if="isDeleting" class="loading-spinner"></span>
            <span v-else class="icon">🗑️</span>
            {{ isDeleting ? '删除中...' : '删除简历' }}
          </button>
        </div>

        <div class="sidebar-section">
          <h4>相关面试</h4>
          <div v-if="relatedInterviews.length > 0">
            <div v-for="interview in relatedInterviews" :key="interview.id" class="interview-item">
              <div class="interview-header">
                <h5>{{ interview.candidate }}</h5>
                <span class="interview-status" :class="getInterviewStatusClass(interview.status)">
                  {{ getInterviewStatusText(interview.status) }}
                </span>
              </div>
              <p class="interview-position">{{ interview.position }}</p>
              <p class="interview-time">{{ formatTime(interview.time) }}</p>
            </div>
          </div>
          <p v-else class="no-data">暂无相关面试</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getResumeDetail, deleteResume as deleteResumeAPI, type ResumeDetailResponse } from '@/service/resume'
import ResumeShow from '@/components/ResumeShow.vue'

const router = useRouter()
const route = useRoute()

// 简历数据
const resume = ref<ResumeDetailResponse>({
  ID: 0,
  CreatedAt: '',
  UpdatedAt: '',
  DeletedAt: null,
  user_id: 0,
  name: '',
  content: '',
  template_id: 1,
  status: 0
})
const relatedInterviews = ref<any[]>([])
const isDeleting = ref(false) // 添加删除状态
const exporting = ref(false) // 添加导出状态

// 模板名称映射
const templateNames: { [key: number]: string } = {
  1: '经典模板',
  2: '现代简约',
  3: '创意设计',
  4: '专业商务'
}

const getTemplateName = (templateId: number) => {
  return templateNames[templateId] || '未知模板'
}

const getStatusClass = (status: number) => {
  const statusClasses: { [key: number]: string } = {
    0: 'status-draft',
    1: 'status-completed',
    2: 'status-published'
  }
  return statusClasses[status] || 'status-draft'
}

const getStatusText = (status: number) => {
  const statusTexts: { [key: number]: string } = {
    0: '草稿',
    1: '已完成',
    2: '已发布'
  }
  return statusTexts[status] || '草稿'
}

const getInterviewStatusClass = (status: string) => {
  const statusClasses: { [key: string]: string } = {
    planned: 'status-planned',
    in_progress: 'status-progress',
    completed: 'status-completed',
    canceled: 'status-canceled'
  }
  return statusClasses[status] || 'status-planned'
}

const getInterviewStatusText = (status: string) => {
  const statusTexts: { [key: string]: string } = {
    planned: '待开始',
    in_progress: '进行中',
    completed: '已完成',
    canceled: '已取消'
  }
  return statusTexts[status] || '待开始'
}

const formatTime = (timestamp: string) => {
  if (!timestamp) return '未知时间'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

const goBack = () => {
  router.go(-1)
}

const editResume = () => {
  router.push(`/resume/edit/${route.params.id}`)
}

const downloadResume = () => {
  // 这里实现下载功能
  alert('下载功能开发中...')
}

const useForInterview = () => {
  router.push('/interview')
}

const duplicateResume = () => {
  // 这里实现复制功能
  alert('复制功能开发中...')
}

const exportToPDF = async () => {
  exporting.value = true
  
  try {
    // 使用html2pdf.js库导出PDF
    const element = document.querySelector('.resume-content-preview')
    if (!element) {
      throw new Error('找不到简历内容')
    }
    
    // 动态导入html2pdf.js
    const html2pdf = await import('html2pdf.js')
    
    const opt = {
      margin: [10, 10, 10, 10],
      filename: `${resume.value.name || '我的简历'}.pdf`,
      image: { type: 'jpeg', quality: 0.98 },
      html2canvas: { 
        scale: 2,
        useCORS: true,
        letterRendering: true
      },
      jsPDF: { 
        unit: 'mm', 
        format: 'a4', 
        orientation: 'portrait' 
      }
    }
    
    // @ts-ignore
    await html2pdf.default().from(element).set(opt).save()
    
    alert('PDF导出成功！')
  } catch (error) {
    console.error('PDF导出失败:', error)
    alert('PDF导出失败，请重试')
  } finally {
    exporting.value = false
  }
}

const deleteResume = async () => {
  if (isDeleting.value) {
    return
  }
  isDeleting.value = true

  if (!confirm('确定要删除这份简历吗？此操作不可恢复。')) {
    isDeleting.value = false
    return
  }

  try {
    const response = await deleteResumeAPI(parseInt(route.params.id as string))
    
    // DELETE请求可能没有返回数据，或者返回空对象
    if (response && (response.data?.code === 1000 || response.status === 200)) {
      alert('简历删除成功')
      router.push('/resume')
    } else {
      alert('简历删除失败，请重试')
    }
  } catch (error) {
    console.error('删除简历失败:', error)
    alert('删除简历失败，请重试')
  } finally {
    isDeleting.value = false
  }
}

const loadResumeDetail = async () => {
  try {
    const response = await getResumeDetail(route.params.id as string)
    console.log('API响应:', response) // 添加调试信息
    
    // API返回格式：{ code: 1000, msg: "success", data: {...} }
    if (response.data.code === 1000) {
      resume.value = response.data.data
      console.log('简历数据加载成功:', resume.value)
    } else {
      console.error('获取简历详情失败:', response.data.msg)
      alert('获取简历详情失败，请重试')
    }
  } catch (error) {
    console.error('获取简历详情失败:', error)
    // 使用模拟数据
    resume.value = {
      ID: parseInt(route.params.id as string),
      CreatedAt: new Date().toISOString(),
      UpdatedAt: new Date().toISOString(),
      DeletedAt: null,
      user_id: 0, // 使用默认值，而不是route.params.id
      name: '软件工程师简历',
      content: `# 张三
男 ｜ 28岁 ｜ 软件工程师 ｜ 本科 ｜ 13800138000 ｜ zhangsan@example.com

## 自我评价
5年Java开发经验，熟悉微服务架构，有大型项目经验，具备良好的团队协作能力和问题解决能力。

## 工作经历

::: start
**腾讯科技**

:::
**高级软件工程师**

:::
**2020-2022**

::: end
负责微信支付系统的开发和维护

::: start
**阿里巴巴**

:::
**软件工程师**

:::
**2018-2020**

::: end
参与电商平台的开发

## 项目经历

::: start
**微信支付系统**

:::
**微服务架构**

:::
**2020-2022**

::: end
负责支付核心模块的设计和实现

::: start
**电商平台**

:::
**React+Node.js**

:::
**2018-2020**

::: end
参与用户端和商家端的开发

## 专业技能
Java, Spring Boot, MySQL, Redis, Docker, Kubernetes, 微服务架构

## 目标职位
高级软件工程师

## 获奖情况
优秀员工奖`,
      template_id: 1,
      status: 1
    }
    console.log('使用模拟数据:', resume.value)
  }
}

const loadRelatedInterviews = async () => {
  try {
    // 这里应该调用API获取相关面试
    // const response = await getRelatedInterviews(route.params.id)
    // relatedInterviews.value = response.data
    
    // 使用模拟数据
    relatedInterviews.value = [
      {
        id: 1,
        candidate: '张三',
        position: '高级软件工程师',
        status: 'completed',
        time: new Date().toISOString()
      }
    ]
  } catch (error) {
    console.error('获取相关面试失败:', error)
  }
}

onMounted(() => {
  console.log('简历详情页面加载，路由参数:', route.params)
  console.log('简历ID:', route.params.id)
  loadResumeDetail()
  loadRelatedInterviews()
})
</script>

<style scoped>
.resume-detail-container {
  min-height: 100vh;
  background: #F8FAFC;
}

.resume-header {
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

.edit-btn, .download-btn {
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

.edit-btn {
  background: #FEF3C7;
  color: #D97706;
}

.edit-btn:hover {
  background: #FDE68A;
}

.download-btn {
  background: #DBEAFE;
  color: #1E40AF;
}

.download-btn:hover {
  background: #BFDBFE;
}

.resume-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 2rem;
}

.resume-preview {
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  padding: 2rem;
}

.resume-header-preview {
  text-align: center;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 2px solid #374151;
}

.resume-header-preview h2 {
  margin: 0 0 1rem 0;
  font-size: 2.5rem;
  font-weight: bold;
  color: #1F2937;
}

.template-badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  background: #DBEAFE;
  color: #1E40AF;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 500;
}

.resume-content-preview {
  font-size: 1rem;
  line-height: 1.6;
}

.resume-sidebar {
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

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.info-item:last-child {
  margin-bottom: 0;
}

.info-item .label {
  font-weight: 500;
  color: #6B7280;
  font-size: 0.9rem;
}

.info-item .value {
  color: #374151;
  font-size: 0.9rem;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 500;
}

.status-draft {
  background: #FEF3C7;
  color: #D97706;
}

.status-completed {
  background: #D1FAE5;
  color: #065F46;
}

.status-published {
  background: #DBEAFE;
  color: #1E40AF;
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

.action-btn.export {
  background: #E0F2FE;
  color: #1E40AF;
}

.action-btn.export:hover {
  background: #C6EAFE;
}

/* 加载动画 */
.loading-spinner {
  display: inline-block;
  width: 12px;
  height: 12px;
  border: 2px solid #DC2626;
  border-radius: 50%;
  border-top-color: transparent;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 禁用状态 */
.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.action-btn:disabled:hover {
  background: inherit;
}

.interview-item {
  padding: 1rem;
  background: #F9FAFB;
  border-radius: 10px;
  margin-bottom: 0.75rem;
}

.interview-item:last-child {
  margin-bottom: 0;
}

.interview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.interview-header h5 {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: #1F2937;
}

.interview-status {
  padding: 0.25rem 0.75rem;
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

.interview-position {
  margin: 0 0 0.25rem 0;
  color: #6B7280;
  font-size: 0.9rem;
}

.interview-time {
  margin: 0;
  color: #9CA3AF;
  font-size: 0.8rem;
}

.no-data {
  color: #9CA3AF;
  font-style: italic;
  text-align: center;
  margin: 1rem 0;
}

@media (max-width: 1024px) {
  .resume-content {
    grid-template-columns: 1fr;
  }
  
  .resume-sidebar {
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
  
  .edit-btn, .download-btn {
    width: 100%;
  }
  
  .resume-content {
    padding: 1rem;
  }
  
  .resume-preview {
    padding: 1.5rem;
  }
  
  .resume-header-preview h2 {
    font-size: 2rem;
  }
}
</style>
