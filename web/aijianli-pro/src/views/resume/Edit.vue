<template>
  <div class="resume-edit-container">
    <div class="edit-header">
      <div class="header-content">
        <button class="back-btn" @click="goBack">
          <span class="icon">←</span>
          返回
        </button>
        <h1>编辑简历</h1>
        <div class="header-actions">
          <button class="save-btn" @click="saveResume" :disabled="saving">
            <span class="icon">{{ saving ? '⏳' : '💾' }}</span>
            {{ saving ? '保存中...' : '保存简历' }}
          </button>
        </div>
      </div>
    </div>

    <div class="edit-content">
      <div class="edit-form">
        <div class="form-section" v-if="false">
          <h3>基本信息</h3>
          <div class="form-row">
            <div class="form-group">
              <label>简历名称</label>
              <input 
                v-model="resumeData.name" 
                type="text" 
                placeholder="请输入简历名称"
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label>模板选择</label>
              <select v-model="resumeData.template_id" class="form-select">
                <option value="1">经典模板</option>
                <option value="2">现代简约</option>
                <option value="3">创意设计</option>
                <option value="4">专业商务</option>
              </select>
            </div>
          </div>
          
          <div class="form-group">
            <label>基本信息</label>
            <textarea 
              v-model="resumeData.basic_info" 
              placeholder="姓名、年龄、学历、联系方式、邮箱等"
              class="form-textarea"
              rows="3"
            ></textarea>
          </div>
        </div>

        <div class="form-section" v-if="false">
          <h3>自我评价</h3>
          <div class="form-group">
            <textarea 
              v-model="resumeData.self_eval" 
              placeholder="请描述您的个人特点、职业目标等"
              class="form-textarea"
              rows="4"
            ></textarea>
          </div>
        </div>

        <div class="form-section" v-if="false">
          <h3>工作经历</h3>
          <div class="form-group">
            <textarea 
              v-model="resumeData.work_exp" 
              placeholder="格式：公司名称-职位-时间-工作描述（每行一个经历）"
              class="form-textarea"
              rows="6"
            ></textarea>
            <div class="form-help">
              示例：腾讯科技-高级软件工程师-2020-2022-负责微信支付系统的开发和维护
            </div>
          </div>
        </div>

        <div class="form-section" v-if="false">
          <h3>项目经历</h3>
          <div class="form-group">
            <textarea 
              v-model="resumeData.project_exp" 
              placeholder="格式：项目名称-技术栈-时间-项目描述（每行一个项目）"
              class="form-textarea"
              rows="6"
            ></textarea>
            <div class="form-help">
              示例：微信支付系统-微服务架构-2020-2022-负责支付核心模块的设计和实现
            </div>
          </div>
        </div>

        <div class="form-section">
          <h3>内容（Markdown）</h3>
          <div class="form-group">
            <textarea 
              v-model="resumeData.content"
              placeholder="# 标题、## 分区、- 列表 等支持实时预览"
              class="form-textarea"
              rows="20"
            ></textarea>
            <div class="form-help">右侧实时预览，支持 Markdown 粗体、标题、列表等</div>
          </div>
        </div>

        <div class="form-section" v-if="false">
          <h3>专业技能</h3>
          <div class="form-group">
            <textarea 
              v-model="resumeData.skills" 
              placeholder="请列出您的专业技能，用逗号分隔"
              class="form-textarea"
              rows="3"
            ></textarea>
          </div>
        </div>

        <div class="form-section" v-if="false">
          <h3>目标职位</h3>
          <div class="form-group">
            <input 
              v-model="resumeData.target_job" 
              type="text" 
              placeholder="请输入目标职位"
              class="form-input"
            />
          </div>
        </div>

        <div class="form-section" v-if="false">
          <h3>获奖情况</h3>
          <div class="form-group">
            <textarea 
              v-model="resumeData.awards" 
              placeholder="请描述您的获奖情况、证书等"
              class="form-textarea"
              rows="3"
            ></textarea>
          </div>
        </div>

        <div class="form-section" v-if="false">
          <h3>简历状态</h3>
          <div class="form-group">
            <select v-model="resumeData.status" class="form-select">
              <option value="draft">草稿</option>
              <option value="completed">已完成</option>
              <option value="published">已发布</option>
            </select>
          </div>
        </div>
      </div>

      <div class="resume-preview">
        <div class="preview-header">
          <h3>实时预览</h3>
          <div class="template-badge">
            模板：{{ getTemplateName(resumeData.template_id) }}
          </div>
        </div>

        <div class="preview-content">
          <ResumeShow :content="resumeData.content" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { createResume, updateResume, getResumeDetail } from '@/service/resume'
import ResumeShow from '@/components/ResumeShow.vue'

const router = useRouter()
const route = useRoute()

// 简历数据
const resumeData = ref({
  name: '',
  basic_info: '',
  work_exp: '',
  project_exp: '',
  self_eval: '',
  skills: '',
  target_job: '',
  awards: '',
  content: '',
  template_id: 1,
  status: 'draft'
})

const saving = ref(false)

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

const workExperience = computed(() => {
  if (!resumeData.value.work_exp) return []
  return resumeData.value.work_exp.split('\n').filter(line => line.trim()).map((exp: string) => {
    const parts = exp.split('-')
    return {
      title: parts[0] || '',
      description: parts.slice(1).join('-') || ''
    }
  })
})

const projectExperience = computed(() => {
  if (!resumeData.value.project_exp) return []
  return resumeData.value.project_exp.split('\n').filter(line => line.trim()).map((project: string) => {
    const parts = project.split('-')
    return {
      title: parts[0] || '',
      description: parts.slice(1).join('-') || ''
    }
  })
})

const goBack = () => {
  if (confirm('确定要离开吗？未保存的更改将丢失。')) {
    router.go(-1)
  }
}

const saveResume = async () => {
  if (!resumeData.value.name.trim()) {
    alert('请输入简历名称')
    return
  }

  saving.value = true
  try {
    // 从token中获取用户ID
    const token = localStorage.getItem('token')
    let userId = 1 // 默认值
    
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        userId = payload.ID || 1
      } catch (e) {
        console.warn('无法解析token中的用户ID，使用默认值:', e)
      }
    }

    const apiData = {
      user_id: userId,
      name: resumeData.value.name,
      basic_info: resumeData.value.basic_info,
      work_exp: resumeData.value.work_exp,
      project_exp: resumeData.value.project_exp,
      self_eval: resumeData.value.self_eval,
      skills: resumeData.value.skills,
      target_job: resumeData.value.target_job,
      awards: resumeData.value.awards,
      template_id: resumeData.value.template_id,
      status: resumeData.value.status,
      content: resumeData.value.content
    }
    
    console.log('保存简历数据:', apiData)
    console.log('content字段内容:', resumeData.value.content)

    const id = route.params.id as string | undefined
    if (id) {
      // 编辑模式：调用更新接口
      await updateResume(Number(id), apiData)
      alert('简历更新成功！')
    } else {
      // 新建模式：调用创建接口
      await createResume(apiData)
      alert('简历创建成功！')
    }
    
    router.push('/resume')
  } catch (error) {
    console.error('保存简历失败:', error)
    alert('保存简历失败，请重试')
  } finally {
    saving.value = false
  }
}

const loadResumeData = async () => {
  try {
    const id = route.params.id as string | undefined
    if (!id) {
      return
    }
    const response = await getResumeDetail(id)
    // 期望后端返回格式：{ code, data }
    const data = response?.data?.data || response?.data
    if (data) {
      console.log('从API获取的简历数据:', data)
      resumeData.value = {
        name: data.name || '',
        basic_info: data.basic_info || '',
        work_exp: data.work_exp || '',
        project_exp: data.project_exp || '',
        self_eval: data.self_eval || '',
        skills: data.skills || '',
        target_job: data.target_job || '',
        awards: data.awards || '',
        content: data.content || '',
        template_id: data.template_id || 1,
        status: data.status || 'draft'
      }
      console.log('设置后的resumeData:', resumeData.value)
    }
  } catch (error) {
    console.error('获取简历详情失败:', error)
  }
}

onMounted(() => {
  loadResumeData()
})

// —— 将左侧表单与右侧 Markdown 内容联动 ——
// 只显示并编辑接口返回的 content，不再自动拼接覆盖
</script>

<style scoped>
.resume-edit-container {
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

.form-help {
  margin-top: 0.5rem;
  font-size: 0.8rem;
  color: #6B7280;
  font-style: italic;
}

.resume-preview {
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
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.template-badge {
  padding: 0.5rem 1rem;
  background: #DBEAFE;
  color: #1E40AF;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 500;
}

.preview-content {
  font-size: 0.9rem;
}

.preview-section {
  margin-bottom: 1.5rem;
}

.preview-section h2 {
  margin: 0 0 1rem 0;
  font-size: 1.5rem;
  font-weight: bold;
  color: #1F2937;
  text-align: center;
}

.preview-section h3 {
  margin: 0 0 0.75rem 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid #E5E7EB;
  padding-bottom: 0.25rem;
}

.preview-section p {
  margin: 0 0 0.75rem 0;
  color: #374151;
  line-height: 1.5;
}

.basic-info {
  text-align: center;
  color: #6B7280;
  line-height: 1.6;
}

.experience-item {
  margin-bottom: 1rem;
  padding-left: 1rem;
  border-left: 2px solid #E5E7EB;
}

.experience-title {
  font-weight: 600;
  color: #1F2937;
  margin-bottom: 0.25rem;
}

.experience-desc {
  color: #6B7280;
  line-height: 1.4;
}

@media (max-width: 1024px) {
  .edit-content {
    grid-template-columns: 1fr;
  }
  
  .resume-preview {
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
  
  .edit-form, .resume-preview {
    padding: 1.5rem;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
