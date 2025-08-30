import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/auth/Register.vue')
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue')
  },
  {
    path: '/resume',
    name: 'ResumeIndex',
    component: () => import('../views/resume/Index.vue')
  },
  {
    path: '/resume/create',
    name: 'ResumeCreate',
    component: () => import('../views/resume/Create.vue')
  },
  {
    path: '/resume/:id',
    name: 'ResumeDetail',
    // 统一跳转到工具页预览/编辑，保留原详情页路由名避免破坏引用
    redirect: to => ({ path: '/resume/tool', query: { id: to.params.id } })
  },
  {
    path: '/resume/edit/:id',
    name: 'ResumeEdit',
    component: () => import('../views/resume/Edit.vue')
  },
  {
    path: '/resume/tool',
    name: 'ResumeTool',
    component: () => import('../views/resume/Tool.vue')
  },
  {
    path: '/interview',
    name: 'InterviewIndex',
    component: () => import('../views/interview/Index.vue')
  },
  {
    path: '/interview/:id',
    name: 'InterviewDetail',
    component: () => import('../views/interview/Detail.vue')
  },
  {
    path: '/interview/room/:id',
    name: 'InterviewRoom',
    component: () => import('../views/interview/Room.vue')
  },
  {
    path: '/interview/edit/:id',
    name: 'InterviewEdit',
    component: () => import('../views/interview/Edit.vue')
  },
  {
    path: '/interview/evaluation/:id',
    name: 'InterviewEvaluation',
    component: () => import('../views/interview/Evaluation.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
