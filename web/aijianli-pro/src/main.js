import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'

// 创建Pinia实例
const pinia = createPinia()

// 创建Vue应用实例
const app = createApp(App)

// 安装插件
app.use(router)
app.use(pinia)

// 挂载应用
app.mount('#app')