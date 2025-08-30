<template>
  <div class="resume-create-container">
    <!-- 模板选择阶段 -->
    <div v-if="!templateSelected" class="template-selection">
      <div class="selection-header">
        <button class="back-btn" @click="goBack">
          <span class="icon">←</span>
          返回
        </button>
        <h1>选择简历模板</h1>
      </div>
      
      <div class="templates-grid">
        <!-- 加载状态 -->
        <div v-if="loadingTemplates" class="loading-templates">
          <div class="loading-spinner"></div>
          <p>正在加载模板...</p>
        </div>
        
        <!-- 模板列表 -->
        <div 
          v-else-if="templates.length > 0"
          v-for="template in templates" 
          :key="template.ID" 
          class="template-card"
          :class="{ selected: selectedTemplate === template.ID }"
          @click="selectTemplate(template.ID)"
        >
          <div class="template-info">
            <div class="template-header">
              <h4>{{ template.name }}</h4>
              <span class="template-badge" :class="getTemplateBadgeClass(template.name)">
                {{ getTemplateType(template.name) }}
              </span>
            </div>
            <div class="template-description" v-html="renderMarkdown(template.show_content)"></div>
          </div>
        </div>
        
        <!-- 无模板状态 -->
        <div v-else class="no-templates">
          <p>暂无可用模板</p>
        </div>
      </div>
      
      <div class="selection-actions">
        <button 
          class="start-btn" 
          @click="startCreation"
          :disabled="!selectedTemplate"
        >
          开始创建简历
        </button>
      </div>
    </div>

    <!-- 聊天创建阶段 -->
    <div v-else class="chat-creation">
      <div class="chat-header">
        <button class="back-btn" @click="backToTemplateSelection">
          <span class="icon">←</span>
          返回选择模板
        </button>
        <h1>AI智能简历助手</h1>
        <div class="progress-indicator">
          <span class="progress-text">{{ currentStep }}/{{ totalSteps }}</span>
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: progressPercentage + '%' }"></div>
          </div>
        </div>
      </div>

      <div class="chat-container">
        <div class="chat-messages" ref="chatContainer">
          <!-- 欢迎消息 -->
          <div class="message ai-message">
            <div class="message-avatar">
              <img src="https://via.placeholder.com/40x40/2563EB/ffffff?text=AI" alt="AI助手" />
            </div>
            <div class="message-content">
              <div class="message-text">
                <p>您好！我是简历助手小A，我将通过几个简单问题帮您收集简历信息。</p>
                <p>每个回答请尽量简洁清晰，现在请告诉我您的全名？</p>
              </div>
              <div class="message-time">{{ getCurrentTime() }}</div>
            </div>
          </div>

          <!-- 合并后的时间序消息流（从旧到新，直接遍历） -->
          <div v-for="(msg, i) in messages" :key="i" :class="['message', msg.role === 'user' ? 'user-message' : 'ai-message']">
            <div class="message-avatar">
              <img v-if="msg.role==='ai'" src="https://via.placeholder.com/40x40/2563EB/ffffff?text=AI" alt="AI助手" />
              <img v-else src="https://via.placeholder.com/40x40/10B981/ffffff?text=我" alt="我" />
            </div>
            <div class="message-content">
              <div class="message-text" v-html="msg.html || msg.content"></div>
              <div class="message-time">{{ msg.time }}</div>
            </div>
          </div>

          <!-- 当前问题 -->
          <div v-if="currentQuestion && !isCompleted" class="message ai-message">
            <div class="message-avatar">
              <img src="https://via.placeholder.com/40x40/2563EB/ffffff?text=AI" alt="AI助手" />
            </div>
            <div class="message-content">
              <div class="message-text">
                <p>{{ currentQuestion.text }}</p>
                <div v-if="currentQuestion.example" class="question-example">
                  <p class="example-label">示例格式：</p>
                  <p class="example-text">{{ currentQuestion.example }}</p>
                </div>
              </div>
              <div class="message-time">{{ getCurrentTime() }}</div>
            </div>
          </div>

          <!-- 完成消息 -->
          <div v-if="isCompleted" class="message ai-message">
            <div class="message-avatar">
              <img src="https://via.placeholder.com/40x40/2563EB/ffffff?text=AI" alt="AI助手" />
            </div>
            <div class="message-content">
              <div class="message-text">
                <p>🎉 太好了！我已经收集了所有必要的信息。</p>
                <p>现在可以保存您的简历信息了！</p>
              </div>
              <div class="message-time">{{ getCurrentTime() }}</div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div v-if="!isCompleted" class="chat-input">
          <div class="input-container">
            <input
              v-model="inputMessage"
              type="text"
              :placeholder="getInputPlaceholder()"
              @keyup.enter="sendAnswer"
              class="message-input"
              :disabled="isProcessing"
            />
            <button 
              @click="sendAnswer" 
              class="send-btn"
              :disabled="!inputMessage.trim() || isProcessing"
            >
              <span v-if="isProcessing" class="loading-spinner"></span>
              <span v-else>发送</span>
            </button>
          </div>
        </div>

        <!-- 完成后的操作按钮 -->
        <div v-if="isCompleted" class="completion-actions">
          <button class="save-btn" @click="saveResume" :disabled="saving">
            <span v-if="saving" class="loading-spinner"></span>
            {{ saving ? '保存中...' : '保存简历' }}
          </button>
          <button class="edit-btn" @click="editResume">重新编辑</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { createResume, getResumeTemplates, type ResumeTemplate } from '@/service/resume'

const router = useRouter()
const chatContainer = ref<HTMLElement | null>(null)
const inputMessage = ref('')
const isProcessing = ref(false)
const isCompleted = ref(false)
const saving = ref(false)
const currentStep = ref(1)

// 模板选择相关
const templateSelected = ref(false)
const selectedTemplate = ref<number | null>(null)

// 简历模板 - 从接口获取
const templates = ref<ResumeTemplate[]>([])
const loadingTemplates = ref(false)

// 简历数据
const resumeData = ref({
  name: '',
  basic_info: '',
  work_exp: '',
  project_exp: '',
  self_eval: '',
  skills: '',
  target_job: '',
  awards: '优秀员工'
})

// 用户回答记录
const userAnswers = ref<Array<{ content: string; time: string }>>([])
const aiConfirmations = ref<Array<{ content: string; nextQuestion: string; time: string }>>([])
// 统一消息流：按时间顺序渲染，避免用户消息“置顶”观感
const messages = computed(() => {
  const ua = userAnswers.value.map(a => ({ role: 'user', content: a.content, time: a.time }))
  const ac = aiConfirmations.value.map(c => ({ role: 'ai', content: `好的，我已经记录了：<strong>${c.content}</strong>${c.nextQuestion ? `<br/>${c.nextQuestion}` : ''}`, time: c.time, html: true }))
  // 简单交错：按出现顺序合并（假设每次回答后会有确认）
  const res: any[] = []
  const n = Math.max(ua.length, ac.length)
  for (let i = 0; i < n; i++) {
    if (ua[i]) res.push(ua[i])
    if (ac[i]) res.push(ac[i])
  }
  return res
})

// 自动滚动到底部（统一定义）
const scrollToBottom = () => {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  })
}

watch(messages, () => {
  nextTick(() => scrollToBottom())
})

// 预定义问题
const questions = [
  {
    id: 1,
    text: '请告诉我您的姓名：',
    example: '',
    field: 'name'
  },
  {
    id: 2,
    text: '您的年龄是多少？',
    example: '',
    field: 'age'
  },
  {
    id: 3,
    text: '您的最高学历是什么？',
    example: '本科、硕士、博士等',
    field: 'education'
  },
  {
    id: 4,
    text: '请提供您的联系电话：',
    example: '',
    field: 'phone'
  },
  {
    id: 5,
    text: '请提供您的邮箱地址：',
    example: '',
    field: 'email'
  },
  {
    id: 6,
    text: '请描述您的工作经历：',
    example: '公司名称-职位-时间段-工作内容（用换行分隔多个经历）',
    field: 'work_exp'
  },
  {
    id: 7,
    text: '请描述您的项目经历：',
    example: '项目名称-技术栈-时间段-项目描述（用换行分隔多个项目）',
    field: 'project_exp'
  },
  {
    id: 8,
    text: '请列出您的专业技能：',
    example: '用逗号分隔，如：Java, Spring Boot, MySQL等',
    field: 'skills'
  },
  {
    id: 9,
    text: '请简单描述一下您的自我评价：',
    example: '突出您的优势和特点',
    field: 'self_eval'
  },
  {
    id: 10,
    text: '您期望的目标职位是什么？',
    example: '如：高级软件工程师、产品经理等',
    field: 'target_job'
  }
]

const totalSteps = questions.length
const currentQuestion = ref<typeof questions[0] | null>(questions[0])
const currentQuestionIndex = ref(0)

const progressPercentage = computed(() => {
  return (currentStep.value / totalSteps) * 100
})

const getCurrentTime = () => {
  return new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

const getInputPlaceholder = () => {
  if (!currentQuestion.value) return '请输入...'
  return `请输入${currentQuestion.value.text.replace('：', '')}`
}

const workExperience = computed(() => {
  if (!resumeData.value.work_exp) return []
  return resumeData.value.work_exp.split('\n').map(exp => {
    const parts = exp.split('-')
    return {
      title: parts[0] || '',
      description: parts.slice(1).join('-') || ''
    }
  })
})

const projectExperience = computed(() => {
  if (!resumeData.value.project_exp) return []
  return resumeData.value.project_exp.split('\n').map(project => {
    const parts = project.split('-')
    return {
      title: parts[0] || '',
      description: parts.slice(1).join('-') || ''
    }
  })
})

const sendAnswer = async () => {
  if (!inputMessage.value.trim() || isProcessing.value) return

  const answer = inputMessage.value.trim()
  const time = getCurrentTime()

  // 添加用户回答
  userAnswers.value.push({
    content: answer,
    time
  })

  // 处理回答
  await processAnswer(answer)

  inputMessage.value = ''
  scrollToBottom()
}

const processAnswer = async (answer: string) => {
  isProcessing.value = true

  // 模拟AI处理时间
  await new Promise(resolve => setTimeout(resolve, 1000))

  // 根据当前问题更新简历数据
  if (!currentQuestion.value) return
  const field = currentQuestion.value.field
  switch (field) {
    case 'name':
      resumeData.value.name = answer
      break
    case 'age':
      resumeData.value.basic_info = `姓名：${resumeData.value.name} | 年龄：${answer}岁`
      break
    case 'education':
      resumeData.value.basic_info += ` | 学历：${answer}`
      break
    case 'phone':
      resumeData.value.basic_info += ` | 联系方式：${answer}`
      break
    case 'email':
      resumeData.value.basic_info += ` | 邮箱：${answer}`
      break
    case 'work_exp':
      resumeData.value.work_exp = answer
      break
    case 'project_exp':
      resumeData.value.project_exp = answer
      break
    case 'skills':
      resumeData.value.skills = answer
      break
    case 'self_eval':
      resumeData.value.self_eval = answer
      break
    case 'target_job':
      resumeData.value.target_job = answer
      break
  }

  // 添加AI确认消息
  const nextQuestion = currentQuestionIndex.value < questions.length - 1 
    ? questions[currentQuestionIndex.value + 1].text 
    : ''

  aiConfirmations.value.push({
    content: answer,
    nextQuestion,
    time: getCurrentTime()
  })

  // 移动到下一个问题
  currentQuestionIndex.value++
  currentStep.value++

  if (currentQuestionIndex.value < questions.length) {
    currentQuestion.value = questions[currentQuestionIndex.value]
  } else {
    // 所有问题都回答完了
    isCompleted.value = true
    currentQuestion.value = null
  }

  isProcessing.value = false
  scrollToBottom()
}

const saveResume = async () => {
  saving.value = true

  try {
    const apiData = {
      user_id: 1, // 这里需要从用户状态获取
      name: resumeData.value.name,
      basic_info: resumeData.value.basic_info,
      work_exp: resumeData.value.work_exp,
      project_exp: resumeData.value.project_exp,
      self_eval: resumeData.value.self_eval,
      awards: resumeData.value.awards,
      target_job: resumeData.value.target_job,
      template_id: selectedTemplate.value || 1 // 使用选中的模板ID
    }

    await createResume(apiData)

    // 保存成功提示
    alert('简历保存成功！')
    router.push('/resume')

  } catch (error) {
    console.error('保存简历失败:', error)
    alert('保存简历失败，请重试')
  } finally {
    saving.value = false
  }
}

const editResume = () => {
  isCompleted.value = false
  currentQuestionIndex.value = 0
  currentStep.value = 1
  currentQuestion.value = questions[0]
  userAnswers.value = []
  aiConfirmations.value = []
  scrollToBottom()
}

const goBack = () => {
  router.go(-1)
}

const backToTemplateSelection = () => {
  templateSelected.value = false
  selectedTemplate.value = null
  currentStep.value = 1
  currentQuestionIndex.value = 0
  currentQuestion.value = questions[0]
  userAnswers.value = []
  aiConfirmations.value = []
  scrollToBottom()
}

const startCreation = () => {
  templateSelected.value = true
  currentStep.value = 1
  currentQuestionIndex.value = 0
  currentQuestion.value = questions[0]
  userAnswers.value = []
  aiConfirmations.value = []
  scrollToBottom()
}

const selectTemplate = (id: number) => {
  selectedTemplate.value = id
}

// 已统一定义在上方

onMounted(() => {
  loadTemplates()
  scrollToBottom()
})

const loadTemplates = async () => {
  loadingTemplates.value = true
  try {
    const response = await getResumeTemplates()
    // API返回格式：{ code: 1000, msg: "success", data: [...] }
    if (response.data.code === 1000) {
      templates.value = response.data.data || []
    } else {
      console.error('获取模板失败:', response.data.msg)
      // 使用默认模板作为备选
      templates.value = [
        {
          ID: 1,
          CreatedAt: '',
          UpdatedAt: '',
          DeletedAt: null,
          name: '经典模板',
          show_content: '**适用场景：** 传统行业、金融、教育等\n\n**特点：**\n- 布局清晰，层次分明\n- 重点突出，易于阅读\n- 专业稳重，适合正式场合\n\n**推荐职位：** 管理、行政、销售等'
        },
        {
          ID: 2,
          CreatedAt: '',
          UpdatedAt: '',
          DeletedAt: null,
          name: '现代简约',
          show_content: '**适用场景：** 互联网、创意、科技等\n\n**特点：**\n- 设计简洁，视觉清爽\n- 重点突出，信息层次清晰\n- 现代感强，适合年轻职场\n\n**推荐职位：** 设计、开发、产品等'
        }
      ]
    }
  } catch (error) {
    console.error('加载简历模板失败:', error)
    // 使用默认模板作为备选
    templates.value = [
      {
        ID: 1,
        CreatedAt: '',
        UpdatedAt: '',
        DeletedAt: null,
        name: '经典模板',
        show_content: '**适用场景：** 传统行业、金融、教育等\n\n**特点：**\n- 布局清晰，层次分明\n- 重点突出，易于阅读\n- 专业稳重，适合正式场合\n\n**推荐职位：** 管理、行政、销售等'
      },
      {
        ID: 2,
        CreatedAt: '',
        UpdatedAt: '',
        DeletedAt: null,
        name: '现代简约',
        show_content: '**适用场景：** 互联网、创意、科技等\n\n**特点：**\n- 设计简洁，视觉清爽\n- 重点突出，信息层次清晰\n- 现代感强，适合年轻职场\n\n**推荐职位：** 设计、开发、产品等'
      }
    ]
  } finally {
    loadingTemplates.value = false
  }
}

const getTemplateType = (name: string) => {
  if (name.includes('经典') || name.includes('传统')) return '经典'
  if (name.includes('现代') || name.includes('简约')) return '现代'
  if (name.includes('创意') || name.includes('设计')) return '创意'
  if (name.includes('商务') || name.includes('专业')) return '商务'
  if (name.includes('技术') || name.includes('开发')) return '技术'
  return '通用'
}



const getTemplateBadgeClass = (name: string) => {
  const type = getTemplateType(name);
  if (type === '经典') return 'badge-classic';
  if (type === '现代') return 'badge-modern';
  if (type === '创意') return 'badge-creative';
  if (type === '商务') return 'badge-professional';
  if (type === '技术') return 'badge-technical';
  return '';
};

// Markdown渲染函数
const renderMarkdown = (markdown: string) => {
  if (!markdown) return '';
  
  console.log('原始Markdown:', markdown);
  
  let content = markdown;
  
  // 处理粗体文本
  content = content.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
  
  // 处理列表项 - 简化版本
  content = content.replace(/^- (.+)$/gm, '<li>$1</li>');
  
  // 将连续的li标签包装成ul
  content = content.replace(/(<li>.*<\/li>)/s, '<ul class="template-list">$1</ul>');
  
  // 处理换行
  content = content.replace(/\n\n/g, '</p><p>');
  content = content.replace(/\n/g, '<br>');
  
  // 包装段落
  if (!content.startsWith('<')) {
    content = '<p>' + content + '</p>';
  }
  
  console.log('渲染后的HTML:', content);
  return content;
};
</script>

<style scoped>
.resume-create-container {
  min-height: 100vh;
  background: #F8FAFC;
}

.template-selection {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.selection-header {
  display: flex;
  align-items: center;
  gap: 2rem;
  margin-bottom: 2rem;
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

.selection-header h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1F2937;
  flex: 1;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.loading-templates {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.loading-templates .loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #EA0202;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

.loading-templates p {
  color: #6B7280;
  font-size: 1rem;
  margin: 0;
}

.no-templates {
  grid-column: 1 / -1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 15px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.no-templates p {
  color: #9CA3AF;
  font-size: 1.1rem;
  margin: 0;
  font-style: italic;
}

.template-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  height: 160px;
  max-height: 160px;
}

.template-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.1);
}

.template-card.selected {
  border: 2px solid #EA0202;
  box-shadow: 0 8px 25px rgba(234, 2, 2, 0.3);
}

/* 模板徽章样式 */
.badge-classic {
  background: #F59E0B;
  color: white;
}

.badge-modern {
  background: #3B82F6;
  color: white;
}

.badge-creative {
  background: #8B5CF6;
  color: white;
}

.badge-professional {
  background: #10B981;
  color: white;
}

.badge-technical {
  background: #EF4444;
  color: white;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 0.5rem 1rem;
  background: #E8F5E9; /* Lighter green header */
  border-radius: 8px;
  margin-bottom: 0.5rem;
}

.preview-header h3 {
  margin: 0;
  font-size: 1.1rem;
  color: #2E7D32; /* Darker green text */
}

.template-badge {
  background: #4CAF50; /* Green badge */
  color: white;
  padding: 0.2rem 0.6rem;
  border-radius: 10px;
  font-size: 0.7rem;
  font-weight: 600;
  white-space: nowrap;
}



.template-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  width: 100%;
}

.template-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.template-header h4 {
  margin: 0;
  font-size: 1rem;
  color: #1F2937;
  font-weight: 600;
}

.template-info p {
  font-size: 0.9rem;
  color: #6B7280;
  margin-top: 0.3rem;
}

.template-description {
  font-size: 0.85rem;
  color: #6B7280;
  margin-top: 0.3rem;
  line-height: 1.4;
  max-height: 100px;
  overflow-y: auto;
  padding-right: 6px;
  flex: 1;
}

/* 自定义滚动条样式 */
.template-description::-webkit-scrollbar {
  width: 3px;
}

.template-description::-webkit-scrollbar-track {
  background: #f5f5f5;
  border-radius: 2px;
}

.template-description::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 2px;
}

.template-description::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

.template-description h3 {
  font-size: 0.9rem;
  color: #1F2937;
  font-weight: 600;
  margin: 0.4rem 0 0.25rem 0;
}

.template-description h4 {
  font-size: 0.85rem;
  color: #374151;
  font-weight: 600;
  margin: 0.3rem 0 0.15rem 0;
}

.template-description strong {
  color: #EA0202;
  font-weight: 600;
}

.template-description ul {
  margin: 0.25rem 0;
  padding-left: 1rem;
}

.template-description li {
  margin-bottom: 0.15rem;
  color: #4B5563;
  font-size: 0.8rem;
}

.selection-actions {
  text-align: center;
}

.start-btn {
  padding: 1rem 2.5rem;
  background: #EA0202;
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  justify-content: center;
}

.start-btn:hover:not(:disabled) {
  background: #D10000;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(234, 2, 2, 0.3);
}

.start-btn:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.chat-creation {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100vh;
  padding: 0;
}

.chat-header {
  background: white;
  border-bottom: 1px solid #E5E7EB;
  padding: 1rem 2rem;
  border-radius: 0;
  display: flex;
  align-items: center;
  gap: 2rem;
  width: 75%;
  max-width: 1000px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
}

.chat-header .back-btn {
  margin-left: 0;
}

.chat-header h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1F2937;
  flex: 1;
}

.progress-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.progress-text {
  font-size: 0.9rem;
  color: #6B7280;
  font-weight: 500;
}

.progress-bar {
  width: 200px;
  height: 8px;
  background: #E5E7EB;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #EA0202, #FF6B6B);
  border-radius: 4px;
  transition: width 0.5s ease;
}

.chat-container {
  background: white;
  border-radius: 0;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  width: 75%;
  max-width: 1000px;
  height: calc(100vh - 80px); /* 减去header高度 */
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex: 1;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 统一聊天区域基础布局与间距，让消息更整齐 */
.chat-messages {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  padding: 1rem 1.25rem;
  height: 520px;
  overflow-y: auto;
}

.message {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  animation: fadeIn 0.3s ease;
}

.message.user-message { flex-direction: row-reverse; }

.message-avatar {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.message-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.message-content {
  flex: 1;
  max-width: 640px;
}

.message.user-message .message-content {
  max-width: 70%;
}

.message-text {
  background: #F6F7F9;
  padding: 12px 14px;
  border-radius: 12px;
  font-size: 0.95rem;
  line-height: 1.5;
  color: #374151;
  border: 1px solid #E5E7EB;
}

.message.user-message .message-text { background: #5b5ce2; color: #fff; border-color: transparent; }

.message-time {
  font-size: 0.8rem;
  color: #9CA3AF;
  margin-top: 6px;
  text-align: left;
}

.message.user-message .message-time { text-align: right; }

.question-example {
  margin-top: 1rem;
  padding: 1rem;
  background: #FEF3C7;
  border-radius: 12px;
  border-left: 4px solid #F59E0B;
}

.example-label {
  font-weight: 600;
  color: #92400E;
  margin: 0 0 0.5rem 0;
  font-size: 0.9rem;
}

.example-text {
  color: #92400E;
  margin: 0;
  font-size: 0.9rem;
  line-height: 1.5;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.chat-input {
  padding: 1.5rem 2rem;
  border-top: 1px solid #E5E7EB;
  background: #FAFAFA;
  border-radius: 0 0 15px 15px;
}

.input-container {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.message-input {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 2px solid #E5E7EB;
  border-radius: 25px;
  font-size: 0.95rem;
  transition: border-color 0.3s ease;
  background: white;
}

.message-input:focus {
  outline: none;
  border-color: #EA0202;
  box-shadow: 0 0 0 3px rgba(234, 2, 2, 0.1);
}

.send-btn {
  padding: 0.75rem 1.5rem;
  background: #EA0202;
  color: white;
  border: none;
  border-radius: 25px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
  min-width: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.send-btn:hover:not(:disabled) {
  background: #D10000;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(234, 2, 2, 0.3);
}

.send-btn:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
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

/* 完成后的操作按钮 */
.completion-actions {
  padding: 2rem;
  border-top: 1px solid #E5E7EB;
  background: #FAFAFA;
  border-radius: 0 0 15px 15px;
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.save-btn, .edit-btn {
  padding: 1rem 2rem;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 120px;
  justify-content: center;
}

.save-btn {
  background: #EA0202;
  color: white;
}

.save-btn:hover:not(:disabled) {
  background: #D10000;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(234, 2, 2, 0.3);
}

.save-btn:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.edit-btn {
  background: #F3F4F6;
  color: #374151;
  border: 2px solid #D1D5DB;
}

.edit-btn:hover {
  background: #E5E7EB;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .chat-header,
  .chat-container {
    width: 90%;
  }
  
  .templates-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  }
}

@media (max-width: 768px) {
  .template-selection {
    padding: 1rem;
  }
  
  .chat-creation {
    padding: 0;
  }
  
  .chat-header,
  .chat-container {
    width: 95%;
  }
  
  .chat-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
    padding: 1rem;
  }
  
  .chat-header h1 {
    font-size: 1.25rem;
  }
  
  .progress-bar {
    width: 150px;
  }
  
  .chat-container {
    height: calc(100vh - 120px); /* 移动端header更高 */
  }
  
  .templates-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  .template-card {
    padding: 1rem;
  }
  
  .template-preview {
    height: 120px;
  }
  
  .chat-messages {
    padding: 1rem;
  }
  
  .message-content {
    max-width: 85%;
  }
  
  .message.user-message .message-content {
    max-width: 85%;
  }
  
  .completion-actions {
    flex-direction: column;
    padding: 1rem;
  }
  
  .save-btn, .edit-btn {
    width: 100%;
  }
  
  .input-container {
    flex-direction: column;
    gap: 1rem;
  }
  
  .send-btn {
    width: 100%;
    border-radius: 12px;
  }
}

@media (max-width: 480px) {
  .chat-header,
  .chat-container {
    width: 100%;
    border-radius: 0;
  }
  
  .chat-header {
    border-radius: 0;
  }
  
  .chat-container {
    border-radius: 0;
    height: calc(100vh - 140px); /* 小屏幕header更高 */
  }
  
  .message {
    gap: 0.5rem;
  }
  
  .message-avatar {
    width: 35px;
    height: 35px;
  }
  
  .message-text {
    padding: 0.75rem 1rem;
    font-size: 0.9rem;
  }
}
</style>
