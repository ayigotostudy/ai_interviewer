<template>
  <div class="chat-room">
    <div class="chat-header">
      <div class="user-info">
        <div class="avatar">
          <img :src="avatar" :alt="title" />
        </div>
        <div class="user-details">
          <h3>{{ title }}</h3>
          <p>{{ subtitle }}</p>
        </div>
      </div>
      <div class="header-actions">
        <button class="action-btn" @click="handleCall" title="语音通话">
          📞
        </button>
        <button class="action-btn" @click="handleVideoCall" title="视频通话">
          📹
        </button>
      </div>
    </div>

    <div class="chat-messages" ref="messagesContainer">
      <div v-for="(message, index) in messages" :key="index" class="message-group">
        <div v-if="message.date" class="date-separator">
          {{ message.date }}
        </div>
        <div class="message" :class="message.type">
          <div class="message-avatar">
            <img :src="message.type === 'user' ? userAvatar : avatar" :alt="message.type === 'user' ? '我' : title" />
          </div>
          <div class="message-content">
            <div class="message-text">{{ message.content }}</div>
            <div class="message-time">{{ message.time }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="chat-input">
      <div class="input-container">
        <input
          v-model="inputMessage"
          type="text"
          placeholder="请输入消息..."
          @keyup.enter="sendMessage"
          :disabled="isProcessing"
        />
        <button class="voice-btn" @click="toggleVoiceInput" :class="{ active: isVoiceActive }">
          {{ isVoiceActive ? '🎤' : '🎙️' }}
        </button>
        <button class="send-btn" @click="sendMessage" :disabled="!inputMessage.trim() || isProcessing">
          发送
        </button>
      </div>
      
      <!-- 语音录制区域 -->
      <div v-if="isVoiceActive" class="voice-recorder">
        <div class="recorder-status">
          <span v-if="!isRecording">点击开始录音</span>
          <span v-else-if="isRecording">正在录音... 点击停止</span>
          <span v-else-if="isProcessing">正在识别...</span>
        </div>
        <button 
          class="record-btn" 
          @click="toggleRecording"
          :class="{ recording: isRecording, processing: isProcessing }"
        >
          {{ isRecording ? '⏹️' : isProcessing ? '⏳' : '🔴' }}
        </button>
      </div>
    </div>

    <!-- 隐藏的文件输入用于语音录制 -->
    <input
      ref="audioInput"
      type="file"
      accept="audio/*"
      style="display: none"
      @change="handleAudioFile"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'

interface Message {
  content: string
  type: 'user' | 'assistant'
  time: string
  date?: string
}

interface Props {
  title: string
  subtitle: string
  avatar: string
  userAvatar: string
  initialMessages?: Message[]
}

const props = withDefaults(defineProps<Props>(), {
  initialMessages: () => []
})

const emit = defineEmits<{
  sendMessage: [message: string]
  voiceMessage: [audioFile: File]
}>()

const messages = ref<Message[]>(props.initialMessages)
const inputMessage = ref('')
const isProcessing = ref(false)
const isVoiceActive = ref(false)
const isRecording = ref(false)
const messagesContainer = ref<HTMLElement>()
const audioInput = ref<HTMLInputElement>()

// 添加消息
const addMessage = (content: string, type: 'user' | 'assistant') => {
  const now = new Date()
  const time = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  
  // 检查是否需要添加日期分隔符
  const today = now.toLocaleDateString('zh-CN')
  const lastMessage = messages.value[messages.value.length - 1]
  const shouldAddDate = !lastMessage || !lastMessage.date || lastMessage.date !== today
  
  if (shouldAddDate) {
    messages.value.push({
      content: today,
      type: 'assistant',
      time: '',
      date: today
    })
  }
  
  messages.value.push({
    content,
    type,
    time
  })
  
  // 滚动到底部
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// 发送消息
const sendMessage = () => {
  if (!inputMessage.value.trim() || isProcessing.value) return
  
  const message = inputMessage.value.trim()
  addMessage(message, 'user')
  emit('sendMessage', message)
  inputMessage.value = ''
}

// 接收消息
const receiveMessage = (content: string) => {
  addMessage(content, 'assistant')
}

// 切换语音输入
const toggleVoiceInput = () => {
  isVoiceActive.value = !isVoiceActive.value
}

// 切换录音
const toggleRecording = () => {
  if (isRecording.value) {
    stopRecording()
  } else {
    startRecording()
  }
}

// 开始录音
const startRecording = () => {
  isRecording.value = true
  // 这里可以集成实际的录音功能
  // 暂时模拟录音过程
  setTimeout(() => {
    if (isRecording.value) {
      // 模拟录音完成，打开文件选择器
      audioInput.value?.click()
    }
  }, 2000)
}

// 停止录音
const stopRecording = () => {
  isRecording.value = false
}

// 处理音频文件
const handleAudioFile = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (!file) return
  
  isProcessing.value = true
  
  try {
    // 调用语音识别接口
    const formData = new FormData()
    formData.append('audio', file)
    
    const response = await fetch('/api/v1/speech/recognize', {
      method: 'POST',
      body: formData
    })
    
    if (response.ok) {
      const result = await response.json()
      const recognizedText = result.text
      
      if (recognizedText) {
        inputMessage.value = recognizedText
        // 自动发送识别后的文本
        sendMessage()
      }
    }
  } catch (error) {
    console.error('语音识别失败:', error)
    addMessage('语音识别失败，请重试', 'assistant')
  } finally {
    isProcessing.value = false
    // 重置文件输入
    if (target) target.value = ''
  }
}

// 通话功能
const handleCall = () => {
  addMessage('语音通话功能开发中...', 'assistant')
}

const handleVideoCall = () => {
  addMessage('视频通话功能开发中...', 'assistant')
}

// 暴露方法给父组件
defineExpose({
  addMessage,
  receiveMessage
})

// 监听消息变化，自动滚动
watch(messages, () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}, { deep: true })
</script>

<style scoped>
.chat-room {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: linear-gradient(135deg, #ff6b35 0%, #f7931e 50%, #8e44ad 100%);
  color: white;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.avatar img {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.user-details h3 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 600;
}

.user-details p {
  margin: 0;
  font-size: 0.9rem;
  opacity: 0.8;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.action-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  color: white;
  font-size: 1.2rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.1);
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.date-separator {
  text-align: center;
  margin: 1rem 0;
  font-size: 0.9rem;
  opacity: 0.7;
  background: rgba(255, 255, 255, 0.1);
  padding: 0.5rem 1rem;
  border-radius: 20px;
  align-self: center;
}

.message {
  display: flex;
  gap: 0.5rem;
  max-width: 70%;
}

.message.user {
  flex-direction: row-reverse;
  align-self: flex-end;
}

.message.assistant {
  align-self: flex-start;
}

.message-avatar img {
  width: 35px;
  height: 35px;
  border-radius: 50%;
  object-fit: cover;
}

.message-content {
  background: rgba(255, 255, 255, 0.1);
  padding: 0.8rem 1rem;
  border-radius: 18px;
  backdrop-filter: blur(10px);
}

.message.user .message-content {
  background: rgba(76, 175, 80, 0.8);
}

.message-text {
  margin-bottom: 0.3rem;
  line-height: 1.4;
}

.message-time {
  font-size: 0.8rem;
  opacity: 0.7;
}

.chat-input {
  padding: 1rem;
  background: rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.input-container {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.input-container input {
  flex: 1;
  padding: 0.8rem 1rem;
  border: none;
  border-radius: 25px;
  background: rgba(255, 255, 255, 0.9);
  color: #333;
  font-size: 1rem;
}

.input-container input:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.5);
}

.voice-btn, .send-btn {
  padding: 0.8rem;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1.2rem;
}

.voice-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.voice-btn.active {
  background: #e74c3c;
  transform: scale(1.1);
}

.send-btn {
  background: #3498db;
  color: white;
  width: 45px;
  height: 45px;
}

.send-btn:hover:not(:disabled) {
  background: #2980b9;
  transform: scale(1.05);
}

.send-btn:disabled {
  background: #95a5a6;
  cursor: not-allowed;
}

.voice-recorder {
  margin-top: 1rem;
  text-align: center;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
}

.recorder-status {
  margin-bottom: 1rem;
  font-size: 0.9rem;
  opacity: 0.8;
}

.record-btn {
  width: 60px;
  height: 60px;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.5rem;
  transition: all 0.3s ease;
}

.record-btn {
  background: #e74c3c;
  color: white;
}

.record-btn.recording {
  background: #c0392b;
  animation: pulse 1s infinite;
}

.record-btn.processing {
  background: #f39c12;
}

@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}

/* 滚动条样式 */
.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
}

.chat-messages::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}
</style>
