<template>
  <div class="interview-edit-container">
    <div class="edit-header">
      <div class="header-content">
        <button class="back-btn" @click="goBack">
          <span class="icon">←</span>
          返回
        </button>
        <h1>编辑面试</h1>
        <div class="header-actions">
          <button class="save-btn" @click="saveInterview" :disabled="saving">
            <span class="icon">{{ saving ? '⏳' : '💾' }}</span>
            {{ saving ? '保存中...' : '保存面试' }}
          </button>
        </div>
      </div>
    </div>

    <div class="edit-content">
      <div class="edit-form">
        <div class="form-section">
          <h3>基本信息</h3>
          <div class="form-row">
            <div class="form-group">
              <label>候选人姓名</label>
              <input 
                v-model="interviewData.candidate" 
                type="text" 
                placeholder="请输入候选人姓名"
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label>目标职位</label>
              <input 
                v-model="interviewData.position" 
                type="text" 
                placeholder="请输入目标职位"
                class="form-input"
              />
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group">
              <label>面试时间</label>
              <input 
                v-model="interviewData.time" 
                type="datetime-local" 
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label>面试类型</label>
              <select v-model="interviewData.type" class="form-select">
                <option value="text">文本面试</option>
                <option value="video">视频面试</option>
                <option value="mixed">混合面试</option>
              </select>
            </div>
          </div>
          
          <div class="form-group">
            <label>面试官</label>
            <input 
              v-model="interviewData.interviewer" 
              type="text" 
              placeholder="请输入面试官姓名"
              class="form-input"
            />
          </div>
        </div>

        <div class="form-section">
          <h3>面试说明</h3>
          <div class="form-group">
            <textarea 
              v-model="interviewData.description" 
              placeholder="请描述面试的目的、重点考察内容等"
              class="form-textarea"
              rows="4"
            ></textarea>
          </div>
        </div>

        <div class="form-section">
          <h3>职位要求</h3>
          <div class="form-group">
            <textarea 
              v-model="interviewData.requirements" 
              placeholder="请列出该职位的主要要求，如技能、经验等"
              class="form-textarea"
              rows="4"
            ></textarea>
          </div>
        </div>

        <div class="form-section">
          <h3>面试状态</h3>
          <div class="form-group">
            <select v-model="interviewData.status" class="form-select">
              <option value="planned">待开始</option>
                              <option value="interviewing">进行中</option>
              <option value="completed">已完成</option>
              <option value="canceled">已取消</option>
            </select>
          </div>
        </div>
      </div>

      <div class="interview-preview">
        <div class="preview-header">
          <h3>面试预览</h3>
        </div>

        <div class="preview-content">
          <div class="preview-section">
            <h4>面试信息</h4>
            <div class="preview-item">
              <span class="label">候选人：</span>
              <span class="value">{{ interviewData.candidate || '未填写' }}</span>
            </div>
            <div class="preview-item">
              <span class="label">目标职位：</span>
              <span class="value">{{ interviewData.position || '未填写' }}</span>
            </div>
            <div class="preview-item">
              <span class="label">面试时间：</span>
              <span class="value">{{ formatPreviewTime(interviewData.time) }}</span>
            </div>
            <div class="preview-item">
              <span class="label">面试类型：</span>
              <span class="value">{{ getInterviewTypeText(interviewData.type) }}</span>
            </div>
            <div class="preview-item">
              <span class="label">面试官：</span>
              <span class="value">{{ interviewData.interviewer || '未指定' }}</span>
            </div>
            <div class="preview-item">
              <span class="label">状态：</span>
              <span class="value status-badge" :class="getStatusClass(interviewData.status)">
                {{ getStatusText(interviewData.status) }}
              </span>
            </div>
          </div>

          <div class="preview-section" v-if="interviewData.description">
            <h4>面试说明</h4>
            <p>{{ interviewData.description }}</p>
          </div>

          <div class="preview-section" v-if="interviewData.requirements">
            <h4>职位要求</h4>
            <p>{{ interviewData.requirements }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { updateMeeting } from '@/service/meeting'

const router = useRouter()
const route = useRoute()

// 面试数据
const interviewData = ref({
  candidate: '',
  position: '',
  time: '',
  type: 'text',
  interviewer: '',
  description: '',
  requirements: '',
  status: 'planned'
})

const saving = ref(false)

const getInterviewTypeText = (type: string) => {
  const typeTexts: { [key: string]: string } = {
    text: '文本面试',
    video: '视频面试',
    mixed: '混合面试'
  }
  return typeTexts[type] || '文本面试'
}

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

const formatPreviewTime = (time: string) => {
  if (!time) return '未设置'
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}

const goBack = () => {
  if (confirm('确定要离开吗？未保存的更改将丢失。')) {
    router.go(-1)
  }
}

const saveInterview = async () => {
  if (!interviewData.value.candidate.trim()) {
    alert('请输入候选人姓名')
    return
  }

  if (!interviewData.value.position.trim()) {
    alert('请输入目标职位')
    return
  }

  if (!interviewData.value.time) {
    alert('请选择面试时间')
    return
  }

  saving.value = true
  try {
    const apiData = {
      id: route.params.id,
      candidate: interviewData.value.candidate,
      position: interviewData.value.position,
      time: new Date(interviewData.value.time).getTime() / 1000,
      type: interviewData.value.type,
      interviewer: interviewData.value.interviewer,
      description: interviewData.value.description,
      requirements: interviewData.value.requirements,
      status: interviewData.value.status
    }

    await updateMeeting(apiData)
    alert('面试信息保存成功！')
    router.push('/interview')
  } catch (error) {
    console.error('保存面试信息失败:', error)
    alert('保存面试信息失败，请重试')
  } finally {
    saving.value = false
  }
}

const loadInterviewData = async () => {
  try {
    // 这里应该调用API获取面试详情
    // const response = await getMeetingDetail(route.params.id)
    // interviewData.value = response.data
    
    // 使用模拟数据
    interviewData.value = {
      candidate: '张三',
      position: '高级软件工程师',
      time: new Date(Date.now() + 86400000).toISOString().slice(0, 16), // 明天
      type: 'text',
      interviewer: 'AI面试官',
      description: '本次面试将重点考察候选人的技术能力、项目经验和团队协作能力。',
      requirements: '1. 5年以上Java开发经验\n2. 熟悉微服务架构\n3. 具备大型项目经验\n4. 良好的沟通能力',
      status: 'planned'
    }
  } catch (error) {
    console.error('获取面试详情失败:', error)
    alert('获取面试详情失败，请重试')
  }
}

onMounted(() => {
  loadInterviewData()
})
</script>

<style scoped>
.interview-edit-container {
  min-height: 100vh;
  background: #F8FAFC;
}

.edit-header {
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

.save-btn {
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

.save-btn:hover:not(:disabled) {
  background: #1D4ED8;
}

.save-btn:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
}

.edit-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
}

.edit-form {
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  padding: 2rem;
}

.form-section {
  margin-bottom: 2rem;
}

.form-section h3 {
  margin: 0 0 1rem 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1F2937;
  border-bottom: 1px solid #E5E7EB;
  padding-bottom: 0.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #374151;
  font-size: 0.9rem;
}

.form-input, .form-select, .form-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-size: 0.9rem;
  color: #374151;
  background: white;
  transition: all 0.3s ease;
}

.form-input:focus, .form-select:focus, .form-textarea:focus {
  outline: none;
  border-color: #2563EB;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.interview-preview {
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  padding: 2rem;
  position: sticky;
  top: 2rem;
  max-height: calc(100vh - 4rem);
  overflow-y: auto;
}

.preview-header {
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #E5E7EB;
}

.preview-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1F2937;
}

.preview-section {
  margin-bottom: 2rem;
}

.preview-section:last-child {
  margin-bottom: 0;
}

.preview-section h4 {
  margin: 0 0 1rem 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid #F3F4F6;
  padding-bottom: 0.5rem;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.preview-item:last-child {
  margin-bottom: 0;
}

.preview-item .label {
  font-weight: 500;
  color: #6B7280;
  font-size: 0.9rem;
}

.preview-item .value {
  color: #374151;
  font-size: 0.9rem;
  font-weight: 500;
}

.status-badge {
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

.preview-section p {
  margin: 0;
  color: #6B7280;
  line-height: 1.6;
  white-space: pre-line;
}

@media (max-width: 1024px) {
  .edit-content {
    grid-template-columns: 1fr;
  }
  
  .interview-preview {
    position: static;
    max-height: none;
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .edit-content {
    padding: 1rem;
  }
  
  .edit-form, .interview-preview {
    padding: 1.5rem;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
