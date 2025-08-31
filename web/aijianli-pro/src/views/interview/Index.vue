<template>
  <div class="interview-container">
    <div class="interview-header">
      <div class="header-content">
        <div class="header-left">
          <button class="back-btn" @click="goBack">
            <span class="icon">←</span>
            返回
          </button>
          <h1>面试管理</h1>
        </div>
        <div class="header-right">
          <button class="create-btn" @click="showCreateForm = true">
            <span class="icon">+</span>
            创建新面试
          </button>
        </div>
      </div>
    </div>

    <div class="interview-content">
      <!-- 创建面试表单 -->
      <div v-if="showCreateForm" class="create-form-overlay">
        <div class="create-form-modal">
          <div class="modal-header">
            <h3>创建新面试</h3>
            <button class="close-btn" @click="showCreateForm = false">×</button>
          </div>
          
          <form @submit.prevent="createInterview" class="create-form">
            <div class="form-group">
              <label for="candidate">候选人姓名 *</label>
              <input
                id="candidate"
                v-model="interviewForm.candidate"
                type="text"
                placeholder="请输入候选人姓名"
                required
              />
            </div>
            
            <div class="form-group">
              <label for="position">面试职位 *</label>
              <input
                id="position"
                v-model="interviewForm.position"
                type="text"
                placeholder="请输入面试职位"
                required
              />
            </div>
            
            <div class="form-group">
              <label for="jobDescription">职位描述 *</label>
              <textarea
                id="jobDescription"
                v-model="interviewForm.job_description"
                placeholder="请描述职位要求和职责"
                rows="4"
                required
              ></textarea>
            </div>
            
            <div class="form-row">
              <div class="form-group">
                <label for="interviewTime">面试时间 *</label>
                <input
                  id="interviewTime"
                  v-model="interviewForm.time"
                  type="datetime-local"
                  required
                />
              </div>
              
              <div class="form-group">
                <label for="duration">预计时长</label>
                <select v-model="interviewForm.duration">
                  <option value="30">30分钟</option>
                  <option value="45">45分钟</option>
                  <option value="60">1小时</option>
                  <option value="90">1.5小时</option>
                  <option value="120">2小时</option>
                </select>
              </div>
            </div>
            

            
            <div class="form-actions">
              <button type="button" class="cancel-btn" @click="showCreateForm = false">
                取消
              </button>
              <button type="submit" class="submit-btn" :disabled="creating">
                <span v-if="creating" class="loading-spinner"></span>
                {{ creating ? '创建中...' : '创建面试' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 筛选和搜索 -->
      <div class="filter-section">
        <div class="search-box">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索候选人姓名或职位..."
            @input="filterInterviews"
          />
          <span class="search-icon">🔍</span>
        </div>
        
        <div class="filter-options">
          <select v-model="statusFilter" @change="filterInterviews">
            <option value="all">全部状态</option>
            <option value="planned">待开始</option>
                            <option value="interviewing">进行中</option>
            <option value="completed">已完成</option>
            <option value="canceled">已取消</option>
          </select>
          
          <select v-model="dateFilter" @change="filterInterviews">
            <option value="all">全部时间</option>
            <option value="today">今天</option>
            <option value="tomorrow">明天</option>
            <option value="week">本周</option>
            <option value="month">本月</option>
          </select>
        </div>
      </div>

      <!-- 面试列表 -->
      <div class="interview-list">
        <div v-if="filteredInterviews.length === 0" class="empty-state">
          <div class="empty-icon">🎯</div>
          <h3>暂无面试安排</h3>
          <p>您还没有安排任何面试，开始创建您的第一个面试吧！</p>
          <button class="create-first-btn" @click="showCreateForm = true">
            立即创建
          </button>
        </div>

        <div v-else class="interview-grid">
          <div v-for="interview in filteredInterviews" :key="interview.id" class="interview-card">
            <div class="interview-header-card">
              <div class="interview-title">
                <h3>{{ interview.candidate }}</h3>
                <span class="position">{{ interview.position }}</span>
              </div>
              <div class="interview-status" :class="getStatusClass(interview.status)">
                {{ getStatusText(interview.status) }}
              </div>
            </div>

            <div class="interview-info">
              <div class="info-item">
                <span class="label">面试时间：</span>
                <span class="value">{{ formatTime(interview.time) }}</span>
              </div>
              
              <div class="info-item">
                <span class="label">职位描述：</span>
                <span class="value">{{ formatDescription(interview.job_description) }}</span>
              </div>
              

            </div>

            <div class="interview-actions">
              <button class="action-btn view-btn" @click="viewInterview(interview.id)">
                <span class="icon">👁️</span>
                查看详情
              </button>
              
              <button 
                v-if="interview.status === 'planned' || interview.status === 'interviewing'" 
                class="action-btn start-btn" 
                @click="startInterview(interview.id)"
              >
                <span class="icon">▶️</span>
                开始面试
              </button>
              
              <button 
                v-if="interview.status === 'planned'" 
                class="action-btn edit-btn" 
                @click="editInterview(interview.id)"
              >
                <span class="icon">✏️</span>
                编辑
              </button>
              
              <button 
                v-if="interview.status === 'planned'" 
                class="action-btn cancel-btn" 
                @click="cancelInterview(interview.id)"
              >
                <span class="icon">❌</span>
                取消
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="filteredInterviews.length > 0" class="pagination">
        <button 
          class="page-btn" 
          :disabled="currentPage === 1"
          @click="changePage(currentPage - 1)"
        >
          上一页
        </button>
        
        <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
        
        <button 
          class="page-btn" 
          :disabled="currentPage === totalPages"
          @click="changePage(currentPage + 1)"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { createMeeting, getMeetingList, updateMeeting, deleteMeeting } from '@/service/meeting'

const router = useRouter()

// 创建表单
const showCreateForm = ref(false)
const creating = ref(false)
const interviewForm = ref({
  candidate: '',
  position: '',
  job_description: '',
  time: '',
  duration: '60'
})

// 搜索和筛选
const searchQuery = ref('')
const statusFilter = ref('all')
const dateFilter = ref('all')

// 分页
const currentPage = ref(1)
const pageSize = 12

// 面试数据
const interviews = ref<any[]>([])
const filteredInterviews = ref<any[]>([])

const getStatusClass = (status: string) => {
  const statusClasses: { [key: string]: string } = {
    planned: 'status-planned',
    interviewing: 'status-progress',
    completed: 'status-completed',
    canceled: 'status-canceled'
  }
  return statusClasses[status] || 'status-planned'
}

const getStatusText = (status: string) => {
  const statusTexts: { [key: string]: string } = {
    planned: '待开始',
    interviewing: '进行中',
    completed: '已完成',
    canceled: '已取消'
  }
  return statusTexts[status] || '待开始'
}

const formatTime = (timestamp: number) => {
  if (!timestamp) return '未设置'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

const formatDescription = (description: string) => {
  if (!description) return '未填写'
  return description.length > 100 ? description.substring(0, 100) + '...' : description
}

const goBack = () => {
  router.go(-1)
}

const createInterview = async () => {
  if (!interviewForm.value.candidate || !interviewForm.value.position || 
      !interviewForm.value.job_description || !interviewForm.value.time) {
    alert('请填写必填字段')
    return
  }

  creating.value = true

  try {
    const interviewData = {
      candidate: interviewForm.value.candidate,
      position: interviewForm.value.position,
      job_description: interviewForm.value.job_description,
      time: new Date(interviewForm.value.time).getTime() / 1000,
      status: 'planned'
    }

    await createMeeting(interviewData)
    alert('面试创建成功！')
    showCreateForm.value = false
    
    // 重置表单
    interviewForm.value = {
      candidate: '',
      position: '',
      job_description: '',
      time: '',
      duration: '60'
    }
    
    // 重新加载数据
    loadInterviews()

  } catch (error: any) {
    console.error('创建面试失败:', error)
    alert('创建面试失败，请重试')
  } finally {
    creating.value = false
  }
}

const viewInterview = (id: string) => {
  router.push(`/interview/${id}`)
}

const startInterview = (id: string) => {
  router.push(`/interview/room/${id}`)
}

const editInterview = (id: string) => {
  router.push(`/interview/edit/${id}`)
}

const cancelInterview = async (id: string) => {
  if (!confirm('确定要取消这场面试吗？此操作不可恢复。')) {
    return
  }

  try {
    // 调用删除面试API
    const response = await deleteMeeting(parseInt(id))
    
    if (response.data.code === 1000) {
      // 从本地列表中移除面试
      const index = interviews.value.findIndex(i => i.id === parseInt(id))
      if (index > -1) {
        interviews.value.splice(index, 1)
        filterInterviews()
      }
      
      alert('面试已成功取消')
    } else {
      alert('取消面试失败: ' + response.data.msg)
    }
  } catch (error) {
    console.error('取消面试失败:', error)
    alert('取消面试失败，请重试')
  }
}

const filterInterviews = () => {
  let filtered = [...interviews.value]

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(interview => 
      interview.candidate?.toLowerCase().includes(query) ||
      interview.position?.toLowerCase().includes(query)
    )
  }

  // 状态过滤
  if (statusFilter.value !== 'all') {
    filtered = filtered.filter(interview => interview.status === statusFilter.value)
  }

  // 日期过滤
  if (dateFilter.value !== 'all') {
    const now = new Date()
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const tomorrow = new Date(today.getTime() + 24 * 60 * 60 * 1000)
    const weekStart = new Date(today.getTime() - today.getDay() * 24 * 60 * 60 * 1000)
    const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)

    filtered = filtered.filter(interview => {
      const interviewTime = new Date(interview.time * 1000)
      
      switch (dateFilter.value) {
        case 'today':
          return interviewTime >= today && interviewTime < tomorrow
        case 'tomorrow':
          return interviewTime >= tomorrow && interviewTime < new Date(tomorrow.getTime() + 24 * 60 * 60 * 1000)
        case 'week':
          return interviewTime >= weekStart
        case 'month':
          return interviewTime >= monthStart
        default:
          return true
      }
    })
  }

  filteredInterviews.value = filtered
  currentPage.value = 1
}

const changePage = (page: number) => {
  currentPage.value = page
}

const totalPages = computed(() => {
  return Math.ceil(filteredInterviews.value.length / pageSize)
})

const loadInterviews = async () => {
  try {
    const response = await getMeetingList()
    console.log('面试列表响应:', response)
    
    // API返回格式：{ code: 1000, msg: "success", data: [...] }
    if (response.data.code === 1000) {
      if (response.data.data && Array.isArray(response.data.data)) {
        // 直接访问 response.data.data，因为这是面试数组
        interviews.value = response.data.data
        console.log('面试数据:', interviews.value)
        filterInterviews()
      } else {
        console.error('获取面试列表失败: 数据结构不正确', response)
        loadMockData()
        filterInterviews()
      }
    } else {
      console.error('获取面试列表失败: 状态码不正确', response.data.code, response.data.msg)
      loadMockData()
      filterInterviews()
    }
  } catch (error: any) {
    console.error('获取面试列表失败:', error)
    loadMockData()
    filterInterviews()
  }
}

const loadMockData = () => {
  interviews.value = [
    {
      id: 1,
      candidate: '张三',
      position: '高级软件工程师',
      job_description: '负责公司核心产品的架构设计和开发，要求5年以上Java开发经验，熟悉Spring Boot、MySQL等技术栈。',
      time: Date.now() / 1000 + 3600, // 1小时后
      status: 'planned',

    },
    {
      id: 2,
      candidate: '李四',
      position: '产品经理',
      job_description: '负责产品规划和需求分析，与开发团队协作推进产品迭代，要求3年以上产品经验。',
      time: Date.now() / 1000 + 7200, // 2小时后
      status: 'planned',

    }
  ]
}

onMounted(() => {
  console.log('面试管理页面加载')
  loadInterviews()
})
</script>

<style scoped>
.interview-container {
  min-height: 100vh;
  background: #F8FAFC;
}

.interview-header {
  background: white;
  border-bottom: 1px solid #E5E7EB;
  padding: 1rem 0;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
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

.header-left h1 {
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
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

.interview-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

/* 创建表单模态框 */
.create-form-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.create-form-modal {
  background: white;
  border-radius: 15px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid #E5E7EB;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1F2937;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #6B7280;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: #F3F4F6;
  color: #374151;
}

.create-form {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #374151;
  font-size: 0.9rem;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 2px solid #E5E7EB;
  border-radius: 10px;
  font-size: 1rem;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #2563EB;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.cancel-btn, .submit-btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.cancel-btn {
  background: #F3F4F6;
  color: #374151;
  border: 1px solid #D1D5DB;
}

.cancel-btn:hover {
  background: #E5E7EB;
}

.submit-btn {
  background: #2563EB;
  color: white;
}

.submit-btn:hover:not(:disabled) {
  background: #1D4ED8;
}

.submit-btn:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 筛选和搜索 */
.filter-section {
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

.search-box input {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 2.5rem;
  border: 2px solid #E5E7EB;
  border-radius: 10px;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.search-box input:focus {
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
}

.filter-options {
  display: flex;
  gap: 1rem;
}

.filter-options select {
  padding: 0.75rem 1rem;
  border: 2px solid #E5E7EB;
  border-radius: 10px;
  font-size: 0.9rem;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
}

.filter-options select:focus {
  outline: none;
  border-color: #2563EB;
}

/* 面试列表 */
.interview-list {
  margin-bottom: 2rem;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state h3 {
  margin: 0 0 1rem 0;
  font-size: 1.5rem;
  color: #374151;
}

.empty-state p {
  margin: 0 0 2rem 0;
  color: #6B7280;
  font-size: 1rem;
}

.create-first-btn {
  background: #2563EB;
  color: white;
  border: none;
  padding: 1rem 2rem;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.create-first-btn:hover {
  background: #1D4ED8;
  transform: translateY(-2px);
}

.interview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 1.5rem;
}

.interview-card {
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
}

.interview-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px -3px rgba(0, 0, 0, 0.1);
}

.interview-header-card {
  padding: 1.5rem;
  border-bottom: 1px solid #F3F4F6;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.interview-title h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1F2937;
}

.position {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background: #DBEAFE;
  color: #1E40AF;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 500;
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

.interview-info {
  padding: 1.5rem;
}

.info-item {
  display: flex;
  margin-bottom: 0.75rem;
}

.info-item:last-child {
  margin-bottom: 0;
}

.info-item .label {
  font-weight: 500;
  color: #6B7280;
  min-width: 80px;
  flex-shrink: 0;
}

.info-item .value {
  color: #374151;
  flex: 1;
}

.interview-actions {
  padding: 1.5rem;
  border-top: 1px solid #F3F4F6;
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.action-btn {
  flex: 1;
  min-width: 120px;
  padding: 0.75rem 1rem;
  border: none;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.view-btn {
  background: #DBEAFE;
  color: #1E40AF;
}

.view-btn:hover {
  background: #BFDBFE;
}

.start-btn {
  background: #D1FAE5;
  color: #065F46;
}

.start-btn:hover {
  background: #A7F3D0;
}

.edit-btn {
  background: #FEF3C7;
  color: #D97706;
}

.edit-btn:hover {
  background: #FDE68A;
}

.cancel-btn {
  background: #FEE2E2;
  color: #DC2626;
}

.cancel-btn:hover {
  background: #FECACA;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
}

.page-btn {
  padding: 0.75rem 1.5rem;
  border: 2px solid #E5E7EB;
  background: white;
  color: #374151;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.3s ease;
}

.page-btn:hover:not(:disabled) {
  border-color: #2563EB;
  color: #2563EB;
}

.page-btn:disabled {
  background: #F9FAFB;
  color: #9CA3AF;
  cursor: not-allowed;
}

.page-info {
  font-size: 0.9rem;
  color: #6B7280;
  font-weight: 500;
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .filter-section {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-box {
    max-width: none;
  }
  
  .filter-options {
    justify-content: center;
  }
  
  .interview-grid {
    grid-template-columns: 1fr;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .interview-actions {
    flex-direction: column;
  }
  
  .action-btn {
    min-width: auto;
  }
  
  .interview-content {
    padding: 1rem;
  }
}
</style>
