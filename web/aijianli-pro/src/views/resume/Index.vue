<template>
  <div class="resume-list-container">
    <div class="list-header">
      <div class="header-content">
        <h1>我的简历</h1>
        <button class="create-btn" @click="createResume">
          <span class="icon">➕</span>
          创建简历
        </button>
      </div>
    </div>

    <div class="list-content">
      <div class="search-filters">
        <div class="search-box">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="搜索简历..."
            class="search-input"
            @input="filterResumes"
          />
          <span class="search-icon">🔍</span>
        </div>
        
        <div class="filters">
          <select v-model="statusFilter" @change="filterResumes" class="filter-select">
            <option value="">全部状态</option>
            <option value="draft">草稿</option>
            <option value="completed">已完成</option>
            <option value="published">已发布</option>
          </select>
          
          <select v-model="templateFilter" @change="filterResumes" class="filter-select">
            <option value="">全部模板</option>
            <option value="1">经典模板</option>
            <option value="2">现代简约</option>
            <option value="3">创意设计</option>
            <option value="4">专业商务</option>
          </select>
        </div>
      </div>

      <div class="resume-grid">
        <div v-for="resume in filteredResumes" :key="resume.ID" class="resume-card">
          <div class="card-header">
            <h3>{{ resume.name }}</h3>
            <div class="card-actions">
              <button class="action-btn view-btn" @click="viewResume(resume)">
                <span class="icon">👁️</span>
                查看
              </button>
              <button class="action-btn edit-btn" @click="editResume(resume)">
                <span class="icon">✏️</span>
                编辑
              </button>
              <button class="action-btn delete-btn" 
                      @click="deleteResume(resume)"
                      :disabled="deletingResumeId === resume.ID">
                <span v-if="deletingResumeId === resume.ID" class="loading-spinner"></span>
                <span v-else class="icon">🗑️</span>
                {{ deletingResumeId === resume.ID ? '删除中...' : '删除' }}
              </button>
            </div>
          </div>
          
          <div class="card-content">
            <div class="resume-preview">
              <ResumeShow :content="resume.content" />
            </div>
          </div>
          
          <div class="card-footer">
            <div class="resume-info">
              <span class="template-badge">
                模板：{{ getTemplateName(resume.template_id) }}
              </span>
              <span class="status-badge" :class="getStatusClass(resume.status)">
                {{ getStatusText(resume.status) }}
              </span>
            </div>
            <div class="resume-actions">
              <button class="use-btn" @click="useForInterview(resume)">
                <span class="icon">🎯</span>
                用于面试
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="filteredResumes.length === 0" class="empty-state">
        <div class="empty-icon">📄</div>
        <h3>暂无简历</h3>
        <p>点击"创建简历"开始制作您的第一份简历</p>
        <button class="create-btn" @click="createResume">
          <span class="icon">➕</span>
          创建简历
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getResumeList, deleteResume as deleteResumeAPI, type ResumeListItem } from '@/service/resume'
import ResumeShow from '@/components/ResumeShow.vue'

const router = useRouter()

// 简历数据
const resumes = ref<ResumeListItem[]>([])
const filteredResumes = ref<ResumeListItem[]>([])
const searchQuery = ref('')
const statusFilter = ref('')
const templateFilter = ref('')
const deletingResumeId = ref<number | null>(null) // 添加删除状态

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

const filterResumes = () => {
  let filtered = resumes.value
  console.log('开始过滤，原始数据:', filtered) // 添加调试信息

  // 搜索过滤
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(resume => 
      resume.name.toLowerCase().includes(query) ||
      resume.content.toLowerCase().includes(query)
    )
    console.log('搜索过滤后:', filtered) // 添加调试信息
  }

  // 状态过滤
  if (statusFilter.value && statusFilter.value !== '') {
    const status = parseInt(statusFilter.value)
    if (!isNaN(status)) {
      filtered = filtered.filter(resume => resume.status === status)
      console.log('状态过滤后:', filtered) // 添加调试信息
    }
  }

  // 模板过滤
  if (templateFilter.value && templateFilter.value !== '') {
    const templateId = parseInt(templateFilter.value)
    if (!isNaN(templateId)) {
      filtered = filtered.filter(resume => resume.template_id === templateId)
      console.log('模板过滤后:', filtered) // 添加调试信息
    }
  }

  filteredResumes.value = filtered
  console.log('最终过滤结果:', filteredResumes.value) // 添加调试信息
}

const createResume = () => {
  router.push('/resume/create')
}

const viewResume = (resume: ResumeListItem) => {
  // 使用后端返回的正确ID字段
  const resumeId = resume.ID
  console.log('查看简历，ID:', resumeId, '简历名称:', resume.name)
  router.push(`/resume/${resumeId}`)
}

const editResume = (resume: ResumeListItem) => {
  // 使用后端返回的正确ID字段
  const resumeId = resume.ID
  console.log('编辑简历，ID:', resumeId, '简历名称:', resume.name)
  router.push(`/resume/edit/${resumeId}`)
}

const deleteResume = async (resume: ResumeListItem) => {
  // 如果正在删除，直接返回
  if (deletingResumeId.value === resume.ID) {
    return
  }

  if (!confirm('确定要删除这份简历吗？此操作不可恢复。')) {
    return
  }

  // 设置删除状态
  deletingResumeId.value = resume.ID

  try {
    // 使用从service导入的deleteResume函数
    const response = await deleteResumeAPI(resume.ID)
    
    // DELETE请求可能没有返回数据，或者返回空对象
    if (response && (response.data?.code === 1000 || response.status === 200)) {
      // 从数组中删除
      const index = resumes.value.findIndex(r => r.ID === resume.ID)
      if (index > -1) {
        resumes.value.splice(index, 1)
        filterResumes()
      }
      alert('简历删除成功')
    } else {
      alert('简历删除失败，请重试')
    }
  } catch (error) {
    console.error('删除简历失败:', error)
    alert('删除简历失败，请重试')
  } finally {
    // 清除删除状态
    deletingResumeId.value = null
  }
}

const useForInterview = (resume: ResumeListItem) => {
  // 使用后端返回的正确ID字段
  const resumeId = resume.ID
  console.log('用于面试，ID:', resumeId, '简历名称:', resume.name)
  // 可以将简历ID作为查询参数传递
  router.push(`/interview?resumeId=${resumeId}`)
}

const loadResumes = async () => {
  try {
    const response = await getResumeList()
    console.log('API响应:', response) // 添加调试信息
    
    // API返回格式：{ code: 1000, msg: "success", data: [...] }
    if (response.data.code === 1000) {
      if (response.data.data && Array.isArray(response.data.data)) {
        // 直接访问 response.data.data，因为这是简历数组
        resumes.value = response.data.data
        console.log('简历数据:', resumes.value) // 添加调试信息
        filterResumes()
      } else {
        console.error('获取简历列表失败: 数据结构不正确', response)
        loadMockData()
        filterResumes()
      }
    } else {
      console.error('获取简历列表失败: 状态码不正确', response.data.code, response.data.msg)
      loadMockData()
      filterResumes()
    }
  } catch (error) {
    console.error('获取简历列表失败:', error)
    loadMockData()
    filterResumes()
  }
}

onMounted(() => {
  // 先添加测试数据确保组件能正常渲染
  resumes.value = [
    {
      user_id: 1,
      name: '测试简历',
      content: `# 张三
男 ｜ 28岁 ｜ 软件工程师 ｜ 本科 ｜ 13800138000 ｜ zhangsan@example.com

## 自我评价
5年Java开发经验，熟悉微服务架构。

## 工作经历
**腾讯科技** - 高级软件工程师 (2020-2022)
负责微信支付系统的开发和维护

## 专业技能
Java, Spring Boot, MySQL, Redis, Docker`,
      template_id: 1,
      status: 1
    }
  ]
  filterResumes()
  
  // 然后尝试加载真实数据
  loadResumes()
})
</script>

<style scoped>
.resume-list-container {
  min-height: 100vh;
  background: #F8FAFC;
}

.list-header {
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

.header-content h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1F2937;
}

.create-btn {
  background: #2563EB;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
}

.create-btn:hover {
  background: #1D4ED8;
}

.list-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
}

.search-filters {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  gap: 1rem;
}

.search-box {
  position: relative;
  flex: 1;
  max-width: 400px;
}

.search-input {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 2.5rem;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-size: 0.9rem;
  color: #374151;
  background: white;
  transition: all 0.3s ease;
}

.search-input:focus {
  outline: none;
  border-color: #2563EB;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  color: #9CA3AF;
  font-size: 1rem;
}

.filters {
  display: flex;
  gap: 1rem;
}

.filter-select {
  padding: 0.75rem;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-size: 0.9rem;
  color: #374151;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
}

.filter-select:focus {
  outline: none;
  border-color: #2563EB;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.resume-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 2rem;
}

.resume-card {
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
}

.resume-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px -3px rgba(0, 0, 0, 0.1);
}

.card-header {
  padding: 1.5rem;
  border-bottom: 1px solid #E5E7EB;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #1F2937;
}

.card-actions {
  display: flex;
  gap: 0.5rem;
}

.action-btn {
  padding: 0.5rem;
  border: none;
  border-radius: 6px;
  font-size: 0.8rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  transition: all 0.3s ease;
}

.view-btn {
  background: #DBEAFE;
  color: #1E40AF;
}

.view-btn:hover {
  background: #BFDBFE;
}

.edit-btn {
  background: #FEF3C7;
  color: #D97706;
}

.edit-btn:hover {
  background: #FDE68A;
}

.delete-btn {
  background: #FEE2E2;
  color: #DC2626;
}

.delete-btn:hover {
  background: #FECACA;
}

.card-content {
  padding: 1.5rem;
  max-height: 300px;
  overflow: hidden;
}

.resume-preview {
  font-size: 0.85rem;
  line-height: 1.4;
}

.card-footer {
  padding: 1.5rem;
  border-top: 1px solid #E5E7EB;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.resume-info {
  display: flex;
  gap: 0.75rem;
}

.template-badge {
  padding: 0.25rem 0.75rem;
  background: #F3F4F6;
  color: #6B7280;
  border-radius: 15px;
  font-size: 0.8rem;
  font-weight: 500;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
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

.use-btn {
  background: #10B981;
  color: white;
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

.use-btn:hover {
  background: #059669;
}

/* 加载动画 */
.loading-spinner {
  display: inline-block;
  width: 12px;
  height: 12px;
  border: 2px solid #ffffff;
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
  transform: none;
  box-shadow: none;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.empty-state h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1F2937;
}

.empty-state p {
  margin: 0 0 2rem 0;
  color: #6B7280;
  font-size: 1rem;
}

@media (max-width: 768px) {
  .search-filters {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-box {
    max-width: none;
  }
  
  .filters {
    justify-content: center;
  }
  
  .resume-grid {
    grid-template-columns: 1fr;
  }
  
  .card-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .card-actions {
    justify-content: center;
  }
  
  .card-footer {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .resume-info {
    justify-content: center;
  }
  
  .use-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
