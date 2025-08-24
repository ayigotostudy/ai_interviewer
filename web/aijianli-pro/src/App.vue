<template>
  <div id="app">
    <!-- 全局Header -->
    <header class="app-header" v-if="showHeader">
      <div class="header-container">
        <!-- Logo区域 -->
        <div class="header-logo">
          <router-link to="/" class="logo-link">
            <span class="logo-icon">🎯</span>
            <span class="logo-text">AI简历助手</span>
          </router-link>
        </div>

        <!-- 导航菜单 -->
        <nav class="header-nav">
          <router-link to="/" class="nav-item" active-class="active">
            <span class="nav-icon">🏠</span>
            首页
          </router-link>
          <router-link to="/resume" class="nav-item" active-class="active">
            <span class="nav-icon">📄</span>
            简历管理
          </router-link>
          <router-link to="/interview" class="nav-item" active-class="active">
            <span class="nav-icon">💼</span>
            面试管理
          </router-link>
        </nav>

        <!-- 用户区域 -->
        <div class="header-user">
          <div v-if="isLoggedIn" class="user-info">
            <span class="user-avatar">👤</span>
            <span class="username">{{ username }}</span>
            <button class="logout-btn" @click="handleLogout">
              <span class="logout-icon">🚪</span>
              退出
            </button>
          </div>
          <div v-else class="auth-buttons">
            <router-link to="/login" class="login-btn">
              <span class="login-icon">🔑</span>
              登录
            </router-link>
            <router-link to="/register" class="register-btn">
              <span class="register-icon">📝</span>
              注册
            </router-link>
          </div>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="app-main" :class="{ 'with-header': showHeader }">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

// 用户状态
const isLoggedIn = ref(!!localStorage.getItem('token'))
const username = ref(localStorage.getItem('username') || '用户')

// 是否显示header（除了登录页面都显示）
const showHeader = computed(() => {
  // 登录页面不显示header
  return route.path !== '/login'
})

// 处理退出登录
const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('username')
  isLoggedIn.value = false
  username.value = '用户'
  router.push('/login')
}

// 初始化时更新一次状态
const updateUserStatus = () => {
  const token = localStorage.getItem('token')
  if (token) {
    isLoggedIn.value = true
    username.value = localStorage.getItem('username') || '用户'
  } else {
    isLoggedIn.value = false
    username.value = '用户'
  }
}

// 初始化
updateUserStatus()
</script>

<style scoped lang="scss">
.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
  backdrop-filter: blur(10px);
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 70px;
}

/* Logo样式 */
.header-logo {
  flex-shrink: 0;
  
  .logo-link {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    text-decoration: none;
    color: white;
    font-weight: 700;
    font-size: 1.5rem;
    transition: transform 0.3s ease;

    &:hover {
      transform: scale(1.05);
    }
  }

  .logo-icon {
    font-size: 2rem;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
  }

  .logo-text {
    background: linear-gradient(45deg, #fff, #f0f0f0);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
}

/* 导航菜单样式 */
.header-nav {
  display: flex;
  gap: 1.5rem;
  align-items: center;
  margin: 0 2rem;

  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: rgba(255, 255, 255, 0.9);
    text-decoration: none;
    font-weight: 500;
    font-size: 1rem;
    padding: 0.75rem 1.25rem;
    border-radius: 12px;
    transition: all 0.3s ease;
    position: relative;
    min-width: 100px;
    justify-content: center;

    &:hover {
      color: white;
      background: rgba(255, 255, 255, 0.15);
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
    }

    &.active {
      color: white;
      background: rgba(255, 255, 255, 0.25);
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
      font-weight: 600;

      &::after {
        content: '';
        position: absolute;
        bottom: -3px;
        left: 50%;
        transform: translateX(-50%);
        width: 30px;
        height: 3px;
        background: #ffd700;
        border-radius: 2px;
        box-shadow: 0 2px 8px rgba(255, 215, 0, 0.4);
      }
    }
  }

  .nav-icon {
    font-size: 1.2rem;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
  }
}

/* 用户区域样式 */
.header-user {
  flex-shrink: 0;
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 1rem;
    color: white;

    .user-avatar {
      font-size: 1.5rem;
      filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
    }

    .username {
      font-weight: 600;
      font-size: 0.95rem;
      min-width: 60px;
    }

    .logout-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      background: rgba(255, 255, 255, 0.15);
      color: white;
      border: 1px solid rgba(255, 255, 255, 0.25);
      padding: 0.5rem 1rem;
      border-radius: 8px;
      cursor: pointer;
      font-size: 0.9rem;
      transition: all 0.3s ease;

      &:hover {
        background: rgba(255, 255, 255, 0.25);
        transform: translateY(-1px);
        box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
      }
    }
  }

  .auth-buttons {
    display: flex;
    gap: 1rem;
    align-items: center;

    .login-btn, .register-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      text-decoration: none;
      padding: 0.75rem 1.5rem;
      border-radius: 10px;
      font-weight: 500;
      font-size: 0.9rem;
      transition: all 0.3s ease;
      min-width: 80px;
      justify-content: center;
    }

    .login-btn {
      color: white;
      background: rgba(255, 255, 255, 0.15);
      border: 1px solid rgba(255, 255, 255, 0.25);

      &:hover {
        background: rgba(255, 255, 255, 0.25);
        transform: translateY(-1px);
        box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
      }
    }

    .register-btn {
      color: #667eea;
      background: white;
      border: 1px solid white;

      &:hover {
        background: #f8f9fa;
        transform: translateY(-1px);
        box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
      }
    }
  }
}

/* 主内容区域 */
.app-main {
  flex: 1;
  
  &.with-header {
    padding-top: 0;
  }
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .header-container {
    padding: 0 1rem;
  }

  .header-nav {
    gap: 1rem;
    margin: 0 1rem;

    .nav-item {
      padding: 0.6rem 1rem;
      font-size: 0.95rem;
      min-width: 90px;
    }
  }
}

@media (max-width: 768px) {
  .header-container {
    flex-direction: column;
    height: auto;
    padding: 1rem;
    gap: 1rem;
  }

  .header-nav {
    flex-wrap: wrap;
    justify-content: center;
    gap: 0.75rem;
    margin: 0;

    .nav-item {
      padding: 0.5rem 0.75rem;
      font-size: 0.9rem;
      min-width: 80px;
    }
  }

  .header-user {
    .user-info, .auth-buttons {
      flex-direction: column;
      gap: 0.5rem;
    }
  }
}

@media (max-width: 480px) {
  .header-logo .logo-text {
    display: none;
  }

  .header-nav {
    .nav-item {
      min-width: 70px;
      padding: 0.5rem;
      
      span:not(.nav-icon) {
        display: none;
      }
    }
  }

  .header-user {
    .username {
      display: none;
    }
  }
}
</style>
