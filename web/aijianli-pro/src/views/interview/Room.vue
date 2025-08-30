<template>
  <div class="interview-room">
    <!-- 顶部导航栏 -->
    <div class="room-header">
      <div class="header-left">
        <button class="back-btn" @click="goBack">
          <span class="icon">←</span>
          返回
        </button>
        <div class="interview-title">
          <h1>{{ interview.candidate }}</h1>
          <span class="position">{{ interview.position }}</span>
        </div>
    </div>

      <div class="header-right">
        <div class="interview-status" :class="getStatusClass(interview.status)">
          {{ getStatusText(interview.status) }}
          </div>
        <button class="info-toggle-btn" @click="toggleInterviewInfo" :class="{ 'active': !isInterviewInfoCollapsed }">
          <span class="icon">ℹ️</span>
        </button>
        </div>
          </div>

    <div class="room-content">
      <!-- 可伸缩的面试信息面板 -->
      <div class="info-panel" :class="{ 'collapsed': isInterviewInfoCollapsed }">
        <div class="panel-content">
          <div class="info-section">
            <h3>面试信息</h3>
            <div class="info-grid">
              <div class="info-item">
                <span class="label">候选人</span>
                <span class="value">{{ interview.candidate }}</span>
              </div>
              <div class="info-item">
                <span class="label">应聘职位</span>
                <span class="value">{{ interview.position }}</span>
              </div>
              <div class="info-item" v-if="interview.job_description">
                <span class="label">职位描述</span>
                <span class="value">{{ interview.job_description }}</span>
              </div>
              <div class="info-item" v-if="interview.time">
                <span class="label">面试时间</span>
                <span class="value">{{ formatTime(interview.time) }}</span>
              </div>

            </div>
          </div>
        </div>
      </div>

      <!-- 主面试区域 -->
      <div class="main-interview-area">
        <!-- 聊天区域 -->
      <div class="chat-container">
          <div class="chat-header">
            <h3>面试对话</h3>
            <div class="chat-actions">
              <button class="action-btn" @click="testSpeechAPI" :disabled="isSubmitting || isTestingSpeech">
                <span class="icon">🎤</span>
                {{ isTestingSpeech ? '测试中…' : '测试语音' }}
              </button>
              <button class="action-btn" @click="showAllEvaluations" :disabled="isSubmitting">
                <span class="icon">📊</span>
                评价总结
              </button>

            </div>
          </div>
          
          <div class="messages-container">
            <div class="messages">
              <div v-for="msg in messages" :key="msg.id" class="message" :class="msg.type">
                <div class="message-avatar">
                  <span class="avatar-icon">{{ msg.type === 'ai' ? '🤖' : '👤' }}</span>
        </div>
                <div class="message-content">
                  <!-- AI思考中的动态效果 -->
                  <div v-if="msg.type === 'ai' && msg.isThinking" class="thinking-container">
                    <div class="thinking-text">AI正在思考中</div>
                    <div class="thinking-dots">
                      <span class="dot"></span>
                      <span class="dot"></span>
                      <span class="dot"></span>
                    </div>
                    <div class="thinking-animation">
                      <div class="brain-waves">
                        <div class="wave"></div>
                        <div class="wave"></div>
                        <div class="wave"></div>
                        <div class="wave"></div>
                      </div>
        </div>
      </div>

                  <!-- 正常消息内容 -->
                  <div v-else>
                    <div class="message-text">{{ msg.content }}</div>
                    
                    <!-- AI消息的评价标签 -->
                    <div v-if="msg.type === 'ai' && msg.evaluation && msg.evaluation.summary" class="evaluation-tag">
                      <span 
                        @click="openEvaluationModal(msg.evaluation)"
                        class="eval-link"
                      >
                        📊 查看评价
                      </span>
                    </div>
        </div>
      </div>
    </div>

              
      </div>
    </div>

          <div class="input-container" :class="{ 'slide-up': isInputSlideUp }">
            <div class="input-wrapper">
              <textarea 
                v-model="inputMessage" 
                placeholder="请输入您的回答..." 
                :disabled="isSubmitting"
                class="message-input"
                @focus="slideUpInput"
                @blur="slideDownInput"
              ></textarea>
              <div class="input-actions">
                <label style="display:flex;align-items:center;gap:6px;margin-right:auto;user-select:none;">
                  <input type="checkbox" v-model="autoSendAfterASR" />
                  <span>自动发送</span>
                </label>
                <button 
                  @click="toggleRecording" 
                  :class="{ 'recording': isRecording }"
                  :disabled="isProcessing || isSubmitting"
                  class="voice-btn"
                >
                  <span class="icon">{{ isRecording ? '⏹️' : '🎤' }}</span>
                  {{ isRecording ? '停止' : '语音' }}
                </button>
                <button 
                  @click="sendMessage" 
                  :disabled="isSubmitting || !inputMessage.trim()"
                  class="send-btn"
                >
                  <span class="icon">📤</span>
                  {{ isSubmitting ? '发送中...' : '发送' }}
                </button>
              </div>
            </div>
          </div>
    </div>

        <!-- 可伸缩的简历侧边栏 -->
        <div class="resume-sidebar" :class="{ 'collapsed': isResumeCollapsed }">
          <div class="sidebar-header">
            <h3>简历预览</h3>
            <button @click="toggleResumeSidebar" class="toggle-btn" :title="isResumeCollapsed ? '展开简历' : '收起简历'">
              <span class="icon">{{ isResumeCollapsed ? '▶' : '◀' }}</span>
            </button>
          </div>

          <div class="sidebar-content" v-show="!isResumeCollapsed">
            <div v-if="resume.content" class="resume-content">
              <div class="markdown-content" v-html="renderMarkdown(resume.content)"></div>
            </div>
            <div v-else class="no-resume">
              <p>暂无简历内容</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 简历选择器 -->
    <div v-if="showResumeSelector" class="resume-selector-overlay">
      <div class="resume-selector">
        <div class="selector-header">
          <h3>选择简历</h3>
          <button @click="closeResumeSelector" class="close-btn">&times;</button>
        </div>
        <p>检测到当前面试没有关联简历，请选择一个简历继续面试：</p>
        
        <div v-if="availableResumes.length > 0" class="resume-list">
          <div 
            v-for="resume in availableResumes" 
            :key="resume.ID"
            class="resume-item"
            @click="selectResume(resume.ID)"
            :class="{ 'selected': selectedResumeId === resume.ID }"
          >
            <div class="resume-name">{{ resume.name }}</div>
            <div class="resume-preview">{{ resume.content.substring(0, 100) }}...</div>
            <div class="resume-meta">
              <span class="template-id">模板: {{ resume.template_id }}</span>
              <span class="status">状态: {{ resume.status === 1 ? '已发布' : '草稿' }}</span>
            </div>
          </div>
        </div>
        
        <div v-else class="no-resumes">
          <p>没有可用的简历，请先创建简历。</p>
          <button @click="goToCreateResume" class="create-resume-btn">创建简历</button>
        </div>

        <div v-if="isUploadingResume" class="uploading-status">
          <span class="spinner"></span>
          正在上传简历...
        </div>
        
        <div class="selector-actions">
          <button @click="closeResumeSelector" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>

    <!-- 评价弹窗 -->
    <div v-if="showEvaluationModal" class="evaluation-modal-overlay" @click="closeEvaluationModal">
      <div class="evaluation-modal" @click.stop>
        <div class="modal-header">
          <h3>面试评价</h3>
          <button @click="closeEvaluationModal" class="close-btn">&times;</button>
        </div>
        
        <div class="modal-content">
          <!-- 评价总结 -->
          <div v-if="currentEvaluationModal?.summary" class="eval-section">
            <h4>📋 评价总结</h4>
            <div class="eval-text">{{ currentEvaluationModal.summary }}</div>
          </div>
          
          <!-- 优点 -->
          <div v-if="currentEvaluationModal?.goodPoints && currentEvaluationModal.goodPoints.length > 0" class="eval-section">
            <h4>✅ 优点</h4>
            <div class="eval-tags">
              <span 
                v-for="point in currentEvaluationModal.goodPoints" 
                :key="point"
                class="eval-tag good-tag"
                @click="copyTagContent(point)"
                title="点击复制"
              >
                {{ point }}
              </span>
            </div>
          </div>
          
          <!-- 待改进 -->
          <div v-if="currentEvaluationModal?.badPoints && currentEvaluationModal.badPoints.length > 0" class="eval-section">
            <h4>❌ 待改进</h4>
            <div class="eval-tags">
              <span 
                v-for="point in currentEvaluationModal.badPoints" 
                :key="point"
                class="eval-tag bad-tag"
                @click="copyTagContent(point)"
                title="点击复制"
              >
                {{ point }}
              </span>
            </div>
          </div>
          
          <!-- 可追问知识点 -->
          <div v-if="currentEvaluationModal?.followUpQuestions && currentEvaluationModal.followUpQuestions.length > 0" class="eval-section">
            <h4>🔍 可追问知识点</h4>
            <div class="eval-tags">
              <span 
                v-for="question in currentEvaluationModal.followUpQuestions" 
                :key="question"
                class="eval-tag followup-tag"
                @click="copyTagContent(point)"
                title="点击复制"
              >
                {{ question }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getMeetingDetail, uploadResume } from '@/service/meeting'
import { getResumeDetail, getResumeList, type ResumeListItem } from '@/service/resume'
import { recognizeSpeech, previewPcmAudio, testSpeechAPI as testSpeechService } from '@/service/speech'
import { aiInterview } from '@/service/meeting'

const router = useRouter()
const route = useRoute()

const interview = ref<{
  candidate: string
  position: string
  resume_id?: number
  job_description?: string
  time?: number
  status?: string

}>({ candidate: '张三', position: '软件工程师' })
const resume = ref<{
  name: string
  basic_info: string
  content: string
}>({ name: '简历', basic_info: '基本信息', content: '' })
const messages = ref<ChatMessage[]>([
  { 
    id: 1, 
    type: 'ai', 
    content: '您好！欢迎参加面试。请介绍一下您的工作经验。',
    evaluation: null,
    isExpanded: false
  }
])
const inputMessage = ref('')

// 面试记录状态
const interviewRecord = ref('')
const isSubmitting = ref(false)

// 简历相关状态
const hasResume = ref(false)
const selectedResumeId = ref<number | null>(null)
const availableResumes = ref<ResumeListItem[]>([])
const showResumeSelector = ref(false)
const isUploadingResume = ref(false)

// 评价相关状态
const currentEvaluation = ref<{
  summary: string
  goodPoints: string[]
  badPoints: string[]
  followUpQuestions: string[]
  nextQuestion: string
} | null>(null)
const showEvaluation = ref(false)

// 重新设计：每条消息的评价数据
const messageEvaluations = ref<Map<number, {
  summary: string
  goodPoints: string[]
  badPoints: string[]
  followUpQuestions: string[]
  evaluation: string
}>>(new Map())

// 定义消息类型
interface ChatMessage {
  id: number
  type: 'ai' | 'user' | 'system'
  content: string
  evaluation: any | null
  isExpanded: boolean
  isThinking?: boolean
}

// 语音识别相关状态
const isRecording = ref(false)
const isProcessing = ref(false)
const recordingTime = ref(0)
const mediaRecorder = ref<MediaRecorder | null>(null)
const audioChunks = ref<Blob[]>([])
const recordingTimer = ref<number | null>(null)
const isTestingSpeech = ref(false)
const autoSendAfterASR = ref(false)

// 复制标签内容
const copyTagContent = (content: string) => {
  navigator.clipboard.writeText(content).then(() => {
    // 可以添加一个提示
    console.log('已复制到剪贴板:', content);
  });
};

// 输入框滑动效果
const slideUpInput = () => {
  isInputSlideUp.value = true;
};

const slideDownInput = () => {
  isInputSlideUp.value = false;
};

// 返回上一页
const goBack = () => {
  router.back();
};

const sendMessage = async () => {
  if (!inputMessage.value.trim() || isSubmitting.value) return
  
  // 检查是否有简历
  if (!hasResume.value) {
    alert('请先选择或上传简历后再开始面试！')
    return
  }
  
  const userMessage = inputMessage.value.trim()
  
  // 设置提交状态
  isSubmitting.value = true
  
  // 添加用户消息到聊天记录
  messages.value.push({
    id: Date.now(),
    type: 'user',
    content: userMessage,
    evaluation: null,
    isExpanded: false
  })
  
  // 清空输入框
  inputMessage.value = ''
  
  // 添加AI思考中的消息
  const thinkingMessageId = Date.now() + 1
  messages.value.push({
    id: thinkingMessageId,
    type: 'ai',
    content: '',
    evaluation: null,
    isExpanded: false,
    isThinking: true
  })
  
  try {
    // 调用AI面试接口
    const response = await aiInterview({
      meeting_id: parseInt(route.params.id as string),
      answer: userMessage
    })
    
    console.log('AI面试响应:', response)
    
    // 找到思考中的消息并更新为正常消息
    const thinkingIndex = messages.value.findIndex(msg => msg.id === thinkingMessageId)
    if (thinkingIndex !== -1) {
      // API返回格式：{ code: 1000, msg: "success", data: {...} }
      if (response.data.code === 1000) {
        // 解析评价数据
        const evaluation = parseEvaluation(response)
        
        // 检查是否有评价数据
        if (evaluation && evaluation.summary) {
          currentEvaluation.value = evaluation
          
          // 只显示下一个问题，不显示评价总结
          if (evaluation.nextQuestion && evaluation.nextQuestion.trim()) {
            console.log('使用解析到的问题:', evaluation.nextQuestion)
            messages.value[thinkingIndex] = {
              id: thinkingMessageId,
              type: 'ai',
              content: evaluation.nextQuestion,
              evaluation: evaluation,
              isExpanded: false,
              isThinking: false
            }
          } else {
            // 如果没有解析到问题，显示完整的reply内容
            const aiResponse = response.data.data?.reply || '感谢您的回答，请继续下一个问题。'
            console.log('使用完整reply内容:', aiResponse)
            messages.value[thinkingIndex] = {
              id: thinkingMessageId,
              type: 'ai',
              content: aiResponse,
              evaluation: evaluation, // 仍然保存评价数据用于标签显示
              isExpanded: false,
              isThinking: false
            }
          }
          
          // 更新面试记录
          if (response.data.data?.interview_record) {
            interviewRecord.value = response.data.data.interview_record
            console.log('面试记录更新:', interviewRecord.value)
          }
          
        } else {
          // 没有评价数据，显示完整的reply内容
          const aiResponse = response.data.data?.reply || '感谢您的回答，请继续下一个问题。'
          console.log('没有评价数据，显示完整reply:', aiResponse)
          messages.value[thinkingIndex] = {
            id: thinkingMessageId,
            type: 'ai',
            content: aiResponse,
            evaluation: null,
            isExpanded: false,
            isThinking: false
          }
        }
        
      } else {
        // 显示错误消息
        messages.value[thinkingIndex] = {
          id: thinkingMessageId,
          type: 'ai',
          content: `抱歉，处理您的回答时出现了问题: ${response.data.msg || '未知错误'}`,
          evaluation: null,
          isExpanded: false,
          isThinking: false
        }
      }
    }
    
  } catch (error) {
    console.error('AI面试接口调用失败:', error)
    
    // 找到思考中的消息并更新为错误消息
    const thinkingIndex = messages.value.findIndex(msg => msg.id === thinkingMessageId)
    if (thinkingIndex !== -1) {
      messages.value[thinkingIndex] = {
        id: thinkingMessageId,
        type: 'ai',
        content: '抱歉，网络连接出现问题，请稍后重试。',
        evaluation: null,
        isExpanded: false,
        isThinking: false
      }
    }
  } finally {
    // 重置提交状态
    isSubmitting.value = false
  }
}

// 开始/停止录音
const toggleRecording = async () => {
  if (isRecording.value) {
    stopRecording()
  } else {
    await startRecording()
  }
}

// 开始录音
const startRecording = async () => {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    
    mediaRecorder.value = new MediaRecorder(stream, {
      mimeType: 'audio/webm;codecs=opus'
    })
    
    audioChunks.value = []
    
    mediaRecorder.value.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.value.push(event.data)
      }
    }
    
    mediaRecorder.value.onstop = async () => {
      await processAudio()
      stream.getTracks().forEach(track => track.stop())
    }
    
    // 开始录音
    mediaRecorder.value.start()
    isRecording.value = true
    recordingTime.value = 0
    
    // 启动计时器
    recordingTimer.value = window.setInterval(() => {
      recordingTime.value++
      
      // 60秒后自动停止
      if (recordingTime.value >= 60) {
        stopRecording()
      }
    }, 1000)
    
  } catch (error) {
    console.error('无法访问麦克风:', error)
    alert('无法访问麦克风，请检查权限设置')
  }
}

// 停止录音
const stopRecording = () => {
  if (mediaRecorder.value && isRecording.value) {
    mediaRecorder.value.stop()
    isRecording.value = false
    
    // 清除计时器
    if (recordingTimer.value) {
      clearInterval(recordingTimer.value)
      recordingTimer.value = null
    }
  }
}

// 处理音频数据
const processAudio = async () => {
  if (audioChunks.value.length === 0) return
  
  isProcessing.value = true
  
  try {
    // 合并音频块
    const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' })
    
    // 转换为PCM格式
    const pcmData = await convertToPCM(audioBlob)
    
    // 调用语音识别接口
    const result = await recognizeSpeech(pcmData)
    
    if (result.success) {
      const text = (result.text || '').trim()
      if (autoSendAfterASR.value && text) {
        inputMessage.value = text
        // 直接发送
        await sendMessage()
      } else {
        // 仅写入输入框
        inputMessage.value = text
      }
    } else {
      alert('语音识别失败: ' + result.error)
    }
    
  } catch (error) {
    console.error('处理音频失败:', error)
    alert('处理音频失败，请重试')
  } finally {
    isProcessing.value = false
    audioChunks.value = []
  }
}

// 转换为 16kHz 单声道 16-bit PCM（提升兼容性与识别质量）
const convertToPCM = async (audioBlob: Blob): Promise<ArrayBuffer> => {
  const arrayBuffer = await blobToArrayBuffer(audioBlob)
  const ac = new (window.AudioContext || (window as any).webkitAudioContext)()
  try {
    const decoded = await ac.decodeAudioData(arrayBuffer)
    const mono = mixToMono(decoded)
    const resampled = await resampleToRate(mono, 16000)
    const pcm = pcmInt16FromAudioBuffer(resampled)
    // 自动静音裁剪（阈值 -40dB，20ms 窗口，首尾各保留 100ms 作为缓冲）
    const trimmed = trimInt16PcmSilence(pcm, 16000, {
      thresholdDb: -40,
      windowMs: 20,
      padMs: 100
    })
    return trimmed.byteLength > 0 ? trimmed : pcm
  } finally {
    ac.close()
  }
}

const blobToArrayBuffer = (blob: Blob): Promise<ArrayBuffer> => {
  return new Promise((resolve, reject) => {
    const fr = new FileReader()
    fr.onload = () => resolve(fr.result as ArrayBuffer)
    fr.onerror = reject
    fr.readAsArrayBuffer(blob)
  })
}

// 混合到单声道（平均各通道）
const mixToMono = (buffer: AudioBuffer): AudioBuffer => {
  const channels = buffer.numberOfChannels
  if (channels === 1) return buffer
  const length = buffer.length
  const sampleRate = buffer.sampleRate
  const out = new AudioBuffer({ length, numberOfChannels: 1, sampleRate })
  const outData = out.getChannelData(0)
  const channelData: Float32Array[] = []
  for (let c = 0; c < channels; c++) channelData.push(buffer.getChannelData(c))
  for (let i = 0; i < length; i++) {
    let sum = 0
    for (let c = 0; c < channels; c++) sum += channelData[c][i]
    outData[i] = sum / channels
  }
  return out
}

// 使用 OfflineAudioContext 高质量重采样到目标采样率
const resampleToRate = async (buffer: AudioBuffer, targetRate: number): Promise<AudioBuffer> => {
  if (buffer.sampleRate === targetRate && buffer.numberOfChannels === 1) return buffer
  const length = Math.ceil(buffer.duration * targetRate)
  const offline = new OfflineAudioContext(1, length, targetRate)
  const source = offline.createBufferSource()
  source.buffer = buffer
  source.connect(offline.destination)
  source.start(0)
  return await offline.startRendering()
}

// 将单声道 AudioBuffer 转为 Int16 PCM
const pcmInt16FromAudioBuffer = (buffer: AudioBuffer): ArrayBuffer => {
  const data = buffer.getChannelData(0)
  const out = new Int16Array(data.length)
  for (let i = 0; i < data.length; i++) {
    const s = Math.max(-1, Math.min(1, data[i]))
    out[i] = s < 0 ? s * 0x8000 : s * 0x7fff
  }
  return out.buffer
}

// 基于 RMS 的前后静音裁剪
const trimInt16PcmSilence = (
  pcmBuffer: ArrayBuffer,
  sampleRate: number,
  options?: { thresholdDb?: number; windowMs?: number; padMs?: number }
): ArrayBuffer => {
  const int16 = new Int16Array(pcmBuffer)
  if (int16.length === 0) return pcmBuffer

  const thresholdDb = options?.thresholdDb ?? -40
  const windowMs = options?.windowMs ?? 20
  const padMs = options?.padMs ?? 100

  const threshold = Math.pow(10, thresholdDb / 20) // 线性幅度
  const windowSize = Math.max(1, Math.floor((sampleRate * windowMs) / 1000))
  const padSamples = Math.floor((sampleRate * padMs) / 1000)

  // 计算滑动 RMS（简化：使用平方和前缀和提高性能）
  const norm = 1 / 32768
  const squares = new Float32Array(int16.length)
  for (let i = 0; i < int16.length; i++) {
    const s = int16[i] * norm
    squares[i] = s * s
  }
  const prefix = new Float32Array(int16.length + 1)
  for (let i = 0; i < int16.length; i++) {
    prefix[i + 1] = prefix[i] + squares[i]
  }

  const rmsAt = (center: number): number => {
    const start = Math.max(0, center - Math.floor(windowSize / 2))
    const end = Math.min(int16.length, start + windowSize)
    const sum = prefix[end] - prefix[start]
    const n = end - start
    return n > 0 ? Math.sqrt(sum / n) : 0
  }

  // 找前后第一个超过阈值的帧
  let startIdx = 0
  while (startIdx < int16.length && rmsAt(startIdx) < threshold) startIdx += windowSize

  let endIdx = int16.length - 1
  while (endIdx > startIdx && rmsAt(endIdx) < threshold) endIdx -= windowSize

  // 应用缓冲
  startIdx = Math.max(0, startIdx - padSamples)
  endIdx = Math.min(int16.length - 1, endIdx + padSamples)

  if (endIdx <= startIdx) return new Int16Array(0).buffer

  const out = int16.subarray(startIdx, endIdx + 1)
  return out.slice().buffer
}

// 测试语音：录制3秒并自动播放本地音频，同时不调用后端
const testSpeechAPI = async () => {
  if (isTestingSpeech.value) return
  try {
    isTestingSpeech.value = true
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    // 选择最优可用的编码
    const preferredTypes = [
      'audio/webm;codecs=opus',
      'audio/webm',
      'audio/mp4;codecs=mp4a.40.2',
      'audio/mp4'
    ]
    const mimeType = preferredTypes.find((t) => (window as any).MediaRecorder?.isTypeSupported?.(t)) || ''
    const mr = new MediaRecorder(stream, mimeType ? { mimeType } : undefined)
    const chunks: Blob[] = []

    mr.ondataavailable = (e) => {
      if (e.data.size > 0) chunks.push(e.data)
    }

    mr.onstop = async () => {
      try {
        const type = mimeType || 'audio/webm'
        const blob = new Blob(chunks, { type })
        const url = URL.createObjectURL(blob)
        const audio = new Audio(url)
        audio.autoplay = true
        audio.controls = true
        audio.play().catch((e) => {
          console.warn('自动播放被阻止，附加到页面，请点击播放', e)
          try {
            document.body.appendChild(audio)
          } catch {}
        })
      } catch (err) {
        console.error('测试语音处理失败:', err)
      } finally {
        stream.getTracks().forEach(t => t.stop())
        isTestingSpeech.value = false
      }
    }

    mr.start()
    // 3 秒后自动停止
    setTimeout(() => {
      if (mr.state !== 'inactive') mr.stop()
    }, 3000)
  } catch (error) {
    console.error('无法访问麦克风:', error)
    alert('无法访问麦克风，请检查权限设置')
    isTestingSpeech.value = false
  }
}

// 检测简历是否存在
const checkResume = async () => {
  try {
    const interviewId = route.params.id as string
    if (!interviewId) return
    
    const response = await getMeetingDetail(interviewId)
    if (response.data.code === 1000) {
      const meetingData = response.data.data
      // 检查面试是否有关联的简历
      hasResume.value = !!(meetingData.resume && meetingData.resume.trim())
      
      if (!hasResume.value) {
        // 如果没有简历，加载可用的简历列表
        await loadAvailableResumes()
        showResumeSelector.value = true
      } else {
        // 如果有简历，不需要额外加载，面试数据中已包含简历内容
        console.log('面试已有简历内容')
      }
    }
  } catch (error) {
    console.error('检测简历失败:', error)
  }
}

// 加载可用的简历列表
const loadAvailableResumes = async () => {
  try {
    const response = await getResumeList()
    if (response.data.code === 1000) {
      availableResumes.value = response.data.data || []
    } else {
      console.error('获取简历列表失败:', response.data.msg)
    }
  } catch (error) {
    console.error('获取简历列表失败:', error)
  }
}

// 选择简历
const selectResume = async (resumeId: number) => {
  selectedResumeId.value = resumeId
  await uploadSelectedResume()
}

// 上传选中的简历
const uploadSelectedResume = async () => {
  if (!selectedResumeId.value) return
  
  isUploadingResume.value = true
  
  try {
    // 获取简历详情
    const resumeResponse = await getResumeDetail(selectedResumeId.value.toString())
    if (resumeResponse.data.code === 1000) {
      const resumeData = resumeResponse.data.data
      
      // 上传简历到面试
      const uploadResponse = await uploadResume({
        meeting_id: parseInt(route.params.id as string),
        resume: resumeData.content
      })
      
      if (uploadResponse.data.code === 1000) {
        // 上传成功
        hasResume.value = true
        showResumeSelector.value = false
        selectedResumeId.value = null
        
        // 重新加载面试数据以获取最新的简历内容
        await loadInterviewData()
        
        // 显示成功消息
        messages.value.push({
          id: Date.now(),
          type: 'system',
          content: `简历"${resumeData.name}"已成功上传到面试`,
          evaluation: null,
          isExpanded: false
        })
        
        alert('简历上传成功！现在可以开始面试了。')
      } else {
        alert('简历上传失败: ' + uploadResponse.data.msg)
      }
    } else {
      alert('获取简历详情失败: ' + resumeResponse.data.msg)
    }
  } catch (error) {
    console.error('上传简历失败:', error)
    alert('上传简历失败，请重试')
  } finally {
    isUploadingResume.value = false
  }
}

// 加载面试数据
const loadInterviewData = async () => {
  try {
    const interviewId = route.params.id as string
    if (interviewId) {
      const response = await getMeetingDetail(interviewId)
      if (response.data.code === 1000) {
        const interviewData = response.data.data
        interview.value = {
          candidate: interviewData.candidate || '候选人',
          position: interviewData.position || '职位',
          resume_id: interviewData.resume_id,
          job_description: interviewData.job_description,
          time: interviewData.time,
          status: interviewData.status,

        }
        
        // 直接使用面试数据中的resume字段
        if (interviewData.resume) {
          resume.value = {
            name: interviewData.candidate || '简历',
            basic_info: interviewData.position || '基本信息',
            content: interviewData.resume
          }
          console.log('从面试数据加载简历内容成功:', resume.value.name)
        } else {
          console.log('面试数据中没有简历内容')
        }
      }
    }
  } catch (error) {
    console.error('加载面试数据失败:', error)
  }
}

// 导出面试记录
const exportInterviewRecord = () => {
  const blob = new Blob([interviewRecord.value], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `面试记录_${interview.value.candidate}_${new Date().toISOString().slice(0, 10)}.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

// 跳转到创建简历页面
const goToCreateResume = () => {
  router.push({ name: 'CreateResume' })
}

// 关闭简历选择器
const closeResumeSelector = () => {
  showResumeSelector.value = false
  selectedResumeId.value = null
}

// 解析AI面试响应中的评价数据
const parseEvaluation = (response: any) => {
  const evaluation = {
    summary: '',
    goodPoints: [] as string[],
    badPoints: [] as string[],
    followUpQuestions: [] as string[],
    nextQuestion: ''
  }
  
  // 提取评价总结 - 适配新的数据结构
  const replyText = response.data?.data?.reply || response.data?.data?.evaluation_summary
  if (replyText) {
    console.log('原始回复文本:', replyText)
    
    // 更精确的解析逻辑
    let currentText = replyText
    
    // 解析✅标记的优点
    const goodMatches = currentText.match(/✅\s*([^❌\n]+)/g)
    if (goodMatches) {
      evaluation.goodPoints = goodMatches.map((match: string) => 
        match.replace('✅', '').trim()
      )
      console.log('提取到优点:', evaluation.goodPoints)
    }
    
    // 解析❌标记的缺点
    const badMatches = currentText.match(/❌\s*([^✅\n]+)/g)
    if (badMatches) {
      evaluation.badPoints = badMatches.map((match: string) => 
        match.replace('❌', '').trim()
      )
      console.log('提取到缺点:', evaluation.badPoints)
    }
    
    // 提取可追问的知识点 - 适配新格式
    const followUpMatch = currentText.match(/可追问的知识点：\s*([^\n]+)/)
    if (followUpMatch) {
      const followUpText = followUpMatch[1].trim()
      evaluation.followUpQuestions = followUpText
        .split(/[,，、\n]/)
        .map((item: string) => item.trim())
        .filter((item: string) => item.length > 0)
      console.log('提取到可追问知识点:', evaluation.followUpQuestions)
    }
    
    // 提取下一个问题 - 适配新格式
    const questionMatch = currentText.match(/问题：\s*([^\n]+)/)
    if (questionMatch) {
      evaluation.nextQuestion = questionMatch[1].trim()
      console.log('提取到问题:', evaluation.nextQuestion)
    } else {
      // 尝试匹配 **问题**： 格式
      const boldQuestionMatch = currentText.match(/\*\*问题\*\*：\s*([^\n]+)/)
      if (boldQuestionMatch) {
        evaluation.nextQuestion = boldQuestionMatch[1].trim()
        console.log('提取到问题(粗体格式):', evaluation.nextQuestion)
      } else {
        // 如果没有找到 "问题：" 格式，尝试其他格式
        const altQuestionMatch = currentText.match(/问题:\s*([^\n]+)/)
        if (altQuestionMatch) {
          evaluation.nextQuestion = altQuestionMatch[1].trim()
          console.log('提取到问题(备用格式):', evaluation.nextQuestion)
        }
      }
    }
    
    evaluation.summary = replyText
    
    console.log('解析后的评价数据:', evaluation)
  }
  
  return evaluation
}

// 显示评价标签
const toggleEvaluation = () => {
  showEvaluation.value = !showEvaluation.value
}

// 切换消息的评价展开/收起
const toggleMessageEvaluation = (messageId: number) => {
  const msg = messages.value.find(msg => msg.id === messageId)
  if (msg) {
    msg.isExpanded = !msg.isExpanded
  }
}

// 评价弹窗相关状态
const showEvaluationModal = ref(false)
const currentEvaluationModal = ref<{
  summary: string
  goodPoints: string[]
  badPoints: string[]
  followUpQuestions: string[]
  nextQuestion: string
} | null>(null)

// 显示评价弹窗
const openEvaluationModal = (evaluation: any) => {
  currentEvaluationModal.value = evaluation
  showEvaluationModal.value = true
}

// 关闭评价弹窗
const closeEvaluationModal = () => {
  showEvaluationModal.value = false
  currentEvaluationModal.value = null
  
  // 清理动态添加的评价内容
    setTimeout(() => {
    const modalContent = document.querySelector('.modal-content')
    if (modalContent) {
      // 保留前4个基础section（评价总结、优点、待改进、可追问知识点）
      const sections = modalContent.querySelectorAll('.eval-section')
      sections.forEach((section, index) => {
        if (index >= 4) {
          section.remove()
        }
      })
    }
  }, 100)
}

// 收集所有消息的评价数据
const showAllEvaluations = () => {
  const allEvaluations: any[] = []
  
  messages.value.forEach(msg => {
    if (msg.evaluation && msg.evaluation.summary) {
      allEvaluations.push({
        id: msg.id,
        content: msg.content,
        evaluation: msg.evaluation
      })
    }
  })
  
  if (allEvaluations.length > 0) {
    // 显示评价汇总弹窗
    currentEvaluationModal.value = {
      summary: `共收集到 ${allEvaluations.length} 条评价数据`,
      goodPoints: [],
      badPoints: [],
      followUpQuestions: [],
      nextQuestion: ''
    }
    
    // 将汇总数据存储到全局状态
    currentEvaluation.value = {
      summary: `评价汇总 (${allEvaluations.length}条)`,
      goodPoints: [],
      badPoints: [],
      followUpQuestions: [],
      nextQuestion: ''
    }
    
    // 显示弹窗
    showEvaluationModal.value = true
    
    // 在弹窗中显示所有评价
    setTimeout(() => {
      const modalContent = document.querySelector('.modal-content')
      if (modalContent) {
        allEvaluations.forEach((item, index) => {
          const evalSection = document.createElement('div')
          evalSection.className = 'eval-section'
          evalSection.innerHTML = `
            <h4>📝 第${index + 1}条评价</h4>
            <div class="eval-text">问题: ${item.content}</div>
            <div class="eval-text">评价: ${item.evaluation.summary}</div>
          `
          modalContent.appendChild(evalSection)
        })
      }
    }, 100)
    
  } else {
    alert('暂无评价数据，请先进行面试对话。')
  }
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

// 格式化时间
const formatTime = (timestamp: number): string => {
  if (!timestamp) return '未设置'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取状态样式类
const getStatusClass = (status: string): string => {
  switch (status) {
    case 'pending':
      return 'status-pending'
    case 'in_progress':
      return 'status-progress'
    case 'completed':
      return 'status-completed'
    case 'cancelled':
      return 'status-cancelled'
    default:
      return 'status-default'
  }
}

// 获取状态文本
const getStatusText = (status: string): string => {
  switch (status) {
    case 'pending':
      return '待开始'
    case 'in_progress':
      return '进行中'
    case 'completed':
      return '已完成'
    case 'cancelled':
      return '已取消'
    default:
      return status || '未知'
  }
}

// 简历预览侧边栏状态
const isResumeCollapsed = ref(false)
const toggleResumeSidebar = () => {
  isResumeCollapsed.value = !isResumeCollapsed.value
}

// 面试信息侧边栏状态
const isInterviewInfoCollapsed = ref(true)
const toggleInterviewInfo = () => {
  isInterviewInfoCollapsed.value = !isInterviewInfoCollapsed.value
}

const isInputSlideUp = ref(false)

onMounted(() => {
  loadInterviewData()
  checkResume() // 在组件挂载时检查简历
})

onUnmounted(() => {
  // 清理录音资源
  if (recordingTimer.value) {
    clearInterval(recordingTimer.value)
  }
  if (mediaRecorder.value && isRecording.value) {
    mediaRecorder.value.stop()
  }
})
</script>

<style scoped>
/* 面试室整体布局 */
.interview-room {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 顶部导航栏 */
.room-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.back-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 0.75rem 1.25rem;
  border-radius: 12px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.back-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.interview-title h1 {
  margin: 0;
  font-size: 1.8rem;
  font-weight: 700;
  color: #1a202c;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.position {
  display: block;
  font-size: 0.9rem;
  color: #718096;
  font-weight: 500;
  margin-top: 0.25rem;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.interview-status {
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.interview-status.status-pending {
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: white;
}

.interview-status.status-progress {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
}

.interview-status.status-completed {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.info-toggle-btn {
  background: linear-gradient(135deg, #f7fafc 0%, #edf2f7 100%);
  border: 1px solid #e2e8f0;
  padding: 0.75rem;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.info-toggle-btn:hover,
.info-toggle-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

/* 主内容区域 */
.room-content {
  flex: 1;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
  min-height: 0;
  overflow: hidden;
}

/* 面试信息面板 */
.info-panel {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin-bottom: 1rem;
  flex-shrink: 0;
}

.info-panel.collapsed {
  max-height: 60px;
}

.info-panel:not(.collapsed) {
  max-height: 180px;
}

.panel-content {
  padding: 1.5rem;
}

.info-section h3 {
  margin: 0 0 1.5rem 0;
  font-size: 1.3rem;
  font-weight: 600;
  color: #2d3748;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.info-section h3::before {
  content: '📋';
  font-size: 1.2rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item .label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #718096;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-item .value {
  font-size: 1rem;
  color: #2d3748;
  font-weight: 500;
  line-height: 1.5;
}

/* 主面试区域 */
.main-interview-area {
  display: flex;
  gap: 1.5rem;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

/* 聊天容器 */
.chat-container {
  flex: 1;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  flex-direction: column;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  min-height: 0;
}

.chat-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
  border-radius: 20px 20px 0 0;
}

.chat-header h3 {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 600;
}

.chat-actions {
  display: flex;
  gap: 0.75rem;
}

.action-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  padding: 0.6rem 1.2rem;
  border-radius: 12px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  backdrop-filter: blur(10px);
}

.action-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(255, 255, 255, 0.2);
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 消息容器 */
.messages-container {
  flex: 1 1 0;
  padding: 2rem;
  overflow-y: auto;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  min-height: 0;
  scrollbar-width: thin;
  scrollbar-color: #cbd5e0 #f1f5f9;
}

.messages-container::-webkit-scrollbar {
  width: 8px;
}

.messages-container::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 4px;
}

.messages-container::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 4px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}

.messages {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding-bottom: 1rem;
}

.message {
  display: flex;
  gap: 1rem;
  align-items: flex-start;
}

.message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 1.3rem;
  flex-shrink: 0;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.message.user .message-avatar {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  box-shadow: 0 4px 15px rgba(16, 185, 129, 0.3);
}

.message-content {
  flex: 1;
  max-width: 75%;
}

.message.user .message-content {
  text-align: right;
}

.message-text {
  background: white;
  padding: 1.2rem 1.5rem;
  border-radius: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  line-height: 1.6;
  color: #2d3748;
  font-size: 1rem;
  border: 1px solid rgba(226, 232, 240, 0.8);
}

.message.user .message-text {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: none;
}

.message.ai .message-text {
  background: white;
  border: 1px solid #e2e8f0;
}

.evaluation-tag {
  margin-top: 1rem;
  text-align: left;
}

.eval-link {
  display: inline-block;
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: white;
  padding: 0.6rem 1.2rem;
  border-radius: 25px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(251, 191, 36, 0.3);
  border: none;
}

.eval-link:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(251, 191, 36, 0.4);
}

/* 输入容器 */
.input-container {
  padding: 1rem;
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border-top: 1px solid rgba(226, 232, 240, 0.8);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  transform: translateY(0);
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.08);
  flex-shrink: 0;
  min-height: 40px;
  max-height: 150px;
  border-radius: 0 0 20px 20px;
  resize: vertical;
  overflow: hidden;
}

.input-container.slide-up {
  transform: translateY(-20px);
  box-shadow: 0 -6px 25px rgba(0, 0, 0, 0.12);
  background: linear-gradient(135deg, #ffffff 0%, #f0f4ff 100%);
}

.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  height: 100%;
}

.message-input {
  width: 100%;
  min-height: 40px;
  max-height: 120px;
  padding: 0.75rem;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  font-size: 0.95rem;
  line-height: 1.5;
  resize: vertical;
  transition: all 0.3s ease;
  font-family: inherit;
  background: #f8fafc;
  color: #2d3748;
  overflow-y: auto;
}

.message-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  background: white;
  transform: translateY(-1px);
}

.message-input::placeholder {
  color: #a0aec0;
  font-weight: 400;
}

.input-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  align-items: center;
  margin-top: auto;
}

.voice-btn,
.send-btn {
  padding: 0.6rem 1.4rem;
  border: none;
  border-radius: 12px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  box-shadow: 0 3px 12px rgba(0, 0, 0, 0.1);
}

.voice-btn {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.voice-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(245, 158, 11, 0.3);
}

.voice-btn.recording {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  animation: pulse 2s infinite;
}

.send-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.send-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.send-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 简历侧边栏 */
.resume-sidebar {
  width: 380px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  flex-direction: column;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  min-height: 0;
  flex-shrink: 0;
}

.resume-sidebar.collapsed {
  width: 80px;
}

.resume-sidebar.collapsed .sidebar-content {
  display: none;
}

.resume-sidebar.collapsed .sidebar-header h3 {
  display: none;
}

.sidebar-header {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  padding: 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
  border-radius: 20px 20px 0 0;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 600;
}

.toggle-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  padding: 0.5rem;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  backdrop-filter: blur(10px);
}

.toggle-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(255, 255, 255, 0.2);
}

.sidebar-content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  scrollbar-width: thin;
  scrollbar-color: #cbd5e0 #f1f5f9;
  min-height: 0;
}

.sidebar-content::-webkit-scrollbar {
  width: 8px;
}

.sidebar-content::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 4px;
}

.sidebar-content::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 4px;
}

.sidebar-content::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}

/* 简历内容样式优化 */
.resume-content {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(226, 232, 240, 0.8);
  height: 100%;
  min-height: 100%;
}

.markdown-content {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  line-height: 1.7;
  color: #2d3748;
  height: 100%;
  min-height: 100%;
}

.markdown-content h1 {
  font-size: 1.8rem;
  font-weight: 700;
  color: #1a202c;
  margin: 0 0 1.5rem 0;
  padding-bottom: 0.75rem;
  border-bottom: 3px solid #667eea;
}

.markdown-content h2 {
  font-size: 1.4rem;
  font-weight: 600;
  color: #ea0202;
  margin: 2rem 0 1rem 0;
}

.markdown-content h3 {
  font-size: 1.2rem;
  font-weight: 600;
  color: #374151;
  margin: 1.5rem 0 0.75rem 0;
}

.markdown-content p {
  margin: 0 0 1rem 0;
  line-height: 1.8;
}

.markdown-content strong {
  font-weight: 600;
  color: #1a202c;
}

.markdown-content em {
  font-style: italic;
  color: #6b7280;
}

.markdown-content code {
  background: #f3f4f6;
  padding: 0.25rem 0.5rem;
  border-radius: 6px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9rem;
  color: #dc2626;
}

.markdown-content pre {
  background: #f8fafc;
  padding: 1.5rem;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  overflow-x: auto;
  margin: 1.5rem 0;
}

.markdown-content pre code {
  background: none;
  padding: 0;
  color: #374151;
}

.markdown-content a {
  color: #2563eb;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-color 0.2s;
}

.markdown-content a:hover {
  border-bottom-color: #2563eb;
}

.markdown-content li {
  margin: 0.5rem 0;
  line-height: 1.7;
}

/* 高亮块样式 */
.markdown-content .highlight-block {
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border: 1px solid #bae6fd;
  border-radius: 12px;
  padding: 1.5rem;
  margin: 1.5rem 0;
  position: relative;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.markdown-content .highlight-block::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6, #06b6d4);
  border-radius: 12px 12px 0 0;
}

.markdown-content .highlight-line {
  padding: 0.75rem 0;
  border-bottom: 1px solid rgba(186, 230, 253, 0.4);
  font-weight: 500;
  color: #1e40af;
}

.markdown-content .highlight-line:last-child {
  border-bottom: none;
}

.markdown-content .basic-info-line {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.75rem;
  padding: 1rem 0;
  margin: 1.5rem 0;
  background: #f8fafc;
  border-radius: 8px;
  padding: 1rem 1.5rem;
}

.markdown-content .info-item {
  color: #374151;
  font-weight: 500;
  white-space: nowrap;
}

.markdown-content .info-separator {
  color: #9ca3af;
  font-weight: 300;
  margin: 0 0.5rem;
}

.markdown-content .highlight-symbol {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  color: #92400e;
  padding: 0.25rem 0.5rem;
  border-radius: 6px;
  font-weight: 600;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9rem;
  border: 1px solid #fbbf24;
  display: inline-block;
  margin: 0 0.25rem;
}

.no-resume {
  text-align: center;
  color: #6b7280;
  font-style: italic;
  padding: 3rem 1.5rem;
  background: #f9fafb;
  border-radius: 12px;
  border: 2px dashed #d1d5db;
}

/* 动画效果 */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}

@keyframes shimmer {
  0% { left: -100%; }
  100% { left: 100%; }
}

@keyframes wave {
  0%, 100% {
    transform: scaleY(0.3);
    opacity: 0.6;
  }
  50% {
    transform: scaleY(1);
    opacity: 1;
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.7;
    transform: scale(1.1);
  }
}

/* 思考状态的机器人头像特殊效果 */
.message.thinking .message-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  animation: thinking-glow 2s infinite ease-in-out;
  box-shadow: 0 0 20px rgba(102, 126, 234, 0.5);
}

@keyframes thinking-glow {
  0%, 100% {
    box-shadow: 0 0 20px rgba(102, 126, 234, 0.5);
    transform: scale(1);
  }
  50% {
    box-shadow: 0 0 30px rgba(102, 126, 234, 0.8);
    transform: scale(1.05);
  }
}

/* 弹框样式 */
.resume-selector-overlay,
.evaluation-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.resume-selector,
.evaluation-modal {
  background: white;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  max-width: 90vw;
  max-height: 90vh;
  overflow: hidden;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(30px) scale(0.95);
    opacity: 0;
  }
  to {
    transform: translateY(0) scale(1);
    opacity: 1;
  }
}

/* 简历选择器样式 */
.resume-selector {
  width: 600px;
  max-height: 80vh;
}

.selector-header {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  padding: 1.5rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.selector-header h3 {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 600;
}

.close-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.1);
}

.resume-selector > p {
  padding: 1.5rem 2rem 1rem;
  margin: 0;
  color: #6b7280;
  font-size: 1rem;
  line-height: 1.6;
}

.resume-list {
  padding: 0 2rem;
  max-height: 400px;
  overflow-y: auto;
}

.resume-item {
  background: #f8fafc;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.resume-item:hover {
  border-color: #10b981;
  background: #f0fdf4;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(16, 185, 129, 0.15);
}

.resume-item.selected {
  border-color: #10b981;
  background: #ecfdf5;
  box-shadow: 0 8px 25px rgba(16, 185, 129, 0.2);
}

.resume-name {
  font-size: 1.2rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 0.5rem;
}

.resume-preview {
  color: #6b7280;
  font-size: 0.9rem;
  line-height: 1.5;
  margin-bottom: 1rem;
}

.resume-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.8rem;
  color: #9ca3af;
}

.no-resumes {
  padding: 2rem;
  text-align: center;
  color: #6b7280;
}

.create-resume-btn {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  margin-top: 1rem;
  transition: all 0.3s ease;
}

.create-resume-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(16, 185, 129, 0.3);
}

.uploading-status {
  padding: 1rem 2rem;
  text-align: center;
  color: #10b981;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #e5e7eb;
  border-top: 2px solid #10b981;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.selector-actions {
  padding: 1.5rem 2rem;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: flex-end;
}

.cancel-btn {
  background: #6b7280;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.cancel-btn:hover {
  background: #4b5563;
  transform: translateY(-1px);
}

/* 评价弹窗样式 */
.evaluation-modal {
  width: 700px;
  max-height: 80vh;
}

.modal-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1.5rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 600;
}

.modal-content {
  padding: 2rem;
  max-height: 500px;
  overflow-y: auto;
}

.eval-section {
  margin-bottom: 2rem;
}

.eval-section h4 {
  margin: 0 0 1rem 0;
  font-size: 1.2rem;
  font-weight: 600;
  color: #1f2937;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.eval-text {
  background: #f8fafc;
  padding: 1rem;
  border-radius: 8px;
  border-left: 4px solid #667eea;
  color: #374151;
  line-height: 1.6;
}

.eval-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.eval-tag {
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.eval-tag.good-tag {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.eval-tag.bad-tag {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
}

.eval-tag.followup-tag {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.eval-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .resume-selector,
  .evaluation-modal {
    width: 95vw;
    margin: 1rem;
  }
  
  .selector-header,
  .modal-header {
    padding: 1rem 1.5rem;
  }
  
  .resume-selector > p,
  .resume-list,
  .modal-content {
    padding-left: 1.5rem;
    padding-right: 1.5rem;
  }
}

/* 动态思考效果样式 */
.thinking-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 16px;
  border: 1px solid rgba(102, 126, 234, 0.2);
  position: relative;
  overflow: hidden;
}

.thinking-text {
  font-size: 1.1rem;
  font-weight: 600;
  color: #667eea;
  margin: 0;
}

.thinking-dots {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.dot {
  width: 8px;
  height: 8px;
  background: #667eea;
  border-radius: 50%;
  animation: pulse 1.5s ease-in-out infinite;
}

.dot:nth-child(2) {
  animation-delay: 0.2s;
}

.dot:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.6;
  }
  50% {
    transform: scale(1.2);
    opacity: 1;
  }
}

.thinking-animation {
  position: relative;
  width: 60px;
  height: 40px;
}

.brain-waves {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  gap: 2px;
}

.wave {
  width: 4px;
  height: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px;
  animation: wave 1.2s ease-in-out infinite;
}

.wave:nth-child(1) { animation-delay: 0s; }
.wave:nth-child(2) { animation-delay: 0.1s; }
.wave:nth-child(3) { animation-delay: 0.2s; }
.wave:nth-child(4) { animation-delay: 0.3s; }

@keyframes wave {
  0%, 100% {
    height: 8px;
    opacity: 0.6;
  }
  50% {
    height: 24px;
    opacity: 1;
  }
}

/* 消息头像发光效果 */
.message.ai .message-avatar {
  position: relative;
}

.message.ai .message-avatar::after {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  z-index: -1;
  animation: glow 2s ease-in-out infinite;
}

@keyframes glow {
  0%, 100% {
    opacity: 0.3;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(1.05);
  }
}
</style>
