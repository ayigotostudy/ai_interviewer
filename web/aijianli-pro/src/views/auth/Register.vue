<template>
  <div class="register-container">
    <div class="register-card">
      <div class="logo-section">
        <div class="logo">
          <img src="https://via.placeholder.com/80x80/2563EB/ffffff?text=AI" alt="AI简历助手" />
        </div>
        <h1>创建账号</h1>
        <p>加入AI简历助手，开启智能求职之旅</p>
      </div>

      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-group">
          <label for="email">邮箱地址</label>
          <div class="email-input-group">
            <input
              id="email"
              v-model="form.email"
              type="email"
              placeholder="请输入您的邮箱"
              required
              :disabled="loading || step > 1"
            />
            <button 
              type="button" 
              class="verify-btn" 
              @click="sendVerifyCode"
              :disabled="loading || !form.email || step > 1"
            >
              {{ verifyBtnText }}
            </button>
          </div>
        </div>

        <div v-if="step >= 2" class="form-group">
          <label for="authCode">验证码</label>
          <input
            id="authCode"
            v-model="form.auth_code"
            type="text"
            placeholder="请输入6位验证码"
            maxlength="6"
            required
            :disabled="loading"
          />
          <p class="help-text">验证码已发送至 {{ form.email }}</p>
        </div>

        <div v-if="step >= 2" class="form-group">
          <label for="password">设置密码</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            placeholder="请设置6位以上密码"
            required
            :disabled="loading"
          />
        </div>

        <div v-if="step >= 2" class="form-group">
          <label for="confirmPassword">确认密码</label>
          <input
            id="confirmPassword"
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            required
            :disabled="loading"
          />
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <div v-if="success" class="success-message">
          {{ success }}
        </div>

        <button type="submit" class="register-btn" :disabled="loading || step < 2">
          <span v-if="loading" class="loading-spinner"></span>
          {{ loading ? '注册中...' : '完成注册' }}
        </button>
      </form>

      <div class="form-footer">
        <p>已有账号？</p>
        <router-link to="/login" class="login-link">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { register, getAuthCode } from '@/service/user'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const success = ref('')
const step = ref(1)
const countdown = ref(0)

const form = ref({
  email: '',
  auth_code: '',
  password: '',
  confirmPassword: ''
})

const verifyBtnText = computed(() => {
  if (countdown.value > 0) {
    return `${countdown.value}s后重试`
  }
  return '获取验证码'
})

const sendVerifyCode = async () => {
  if (!form.value.email) {
    error.value = '请先输入邮箱地址'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await getAuthCode(form.value.email)
    
    if (response.code === 0 || response.code === '0') {
      success.value = '验证码已发送，请查收邮件'
      step.value = 2
      startCountdown()
    } else {
      error.value = response.msg || '发送验证码失败，请重试'
    }
  } catch (err: any) {
    console.error('发送验证码失败:', err)
    error.value = err.response?.data?.msg || '发送验证码失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

const startCountdown = () => {
  countdown.value = 60
  const timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
    }
  }, 1000)
}

const handleRegister = async () => {
  if (step.value < 2) {
    error.value = '请先完成邮箱验证'
    return
  }

  if (!form.value.email || !form.value.auth_code || !form.value.password) {
    error.value = '请填写完整的注册信息'
    return
  }

  if (form.value.password !== form.value.confirmPassword) {
    error.value = '两次输入的密码不一致'
    return
  }

  if (form.value.password.length < 6) {
    error.value = '密码长度至少6位'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await register({
      email: form.value.email,
      auth_code: form.value.auth_code,
      password: form.value.password
    })

    if (response.code === 0 || response.code === '0') {
      success.value = '注册成功！正在跳转到登录页面...'
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    } else {
      error.value = response.msg || '注册失败，请重试'
    }
  } catch (err: any) {
    console.error('注册失败:', err)
    error.value = err.response?.data?.msg || '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #2563EB 0%, #1E40AF 50%, #1E3A8A 100%);
  padding: 2rem;
}

.register-card {
  background: white;
  border-radius: 20px;
  padding: 3rem;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  width: 100%;
  max-width: 450px;
}

.logo-section {
  text-align: center;
  margin-bottom: 2rem;
}

.logo {
  margin-bottom: 1rem;
}

.logo img {
  width: 80px;
  height: 80px;
  border-radius: 20px;
}

.logo-section h1 {
  margin: 0 0 0.5rem 0;
  font-size: 1.8rem;
  font-weight: 700;
  color: #1F2937;
}

.logo-section p {
  margin: 0;
  color: #6B7280;
  font-size: 0.9rem;
}

.register-form {
  margin-bottom: 2rem;
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

.email-input-group {
  display: flex;
  gap: 0.5rem;
}

.email-input-group input {
  flex: 1;
}

.verify-btn {
  padding: 0.75rem 1rem;
  background: #10B981;
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
  min-width: 100px;
}

.verify-btn:hover:not(:disabled) {
  background: #059669;
}

.verify-btn:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
}

.form-group input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 2px solid #E5E7EB;
  border-radius: 10px;
  font-size: 1rem;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #2563EB;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.form-group input:disabled {
  background-color: #F9FAFB;
  cursor: not-allowed;
}

.help-text {
  margin: 0.5rem 0 0 0;
  font-size: 0.8rem;
  color: #6B7280;
}

.error-message {
  background: #FEF2F2;
  color: #DC2626;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  border: 1px solid #FECACA;
}

.success-message {
  background: #F0FDF4;
  color: #16A34A;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  border: 1px solid #BBF7D0;
}

.register-btn {
  width: 100%;
  padding: 0.875rem 1rem;
  background: #2563EB;
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.register-btn:hover:not(:disabled) {
  background: #1D4ED8;
  transform: translateY(-1px);
}

.register-btn:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
  transform: none;
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

.form-footer {
  text-align: center;
  padding-top: 1.5rem;
  border-top: 1px solid #E5E7EB;
}

.form-footer p {
  margin: 0 0 0.5rem 0;
  color: #6B7280;
  font-size: 0.9rem;
}

.login-link {
  color: #2563EB;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.9rem;
  transition: color 0.3s ease;
}

.login-link:hover {
  color: #1D4ED8;
  text-decoration: underline;
}

@media (max-width: 480px) {
  .register-container {
    padding: 1rem;
  }
  
  .register-card {
    padding: 2rem 1.5rem;
  }
  
  .email-input-group {
    flex-direction: column;
  }
  
  .verify-btn {
    min-width: auto;
  }
}
</style>
