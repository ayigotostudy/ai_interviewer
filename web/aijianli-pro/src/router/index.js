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
    component: () => import('../views/resume/Detail.vue')
  },
  {
    path: '/resume/edit/:id',
    name: 'ResumeEdit',
    component: () => import('../views/resume/Edit.vue')
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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
