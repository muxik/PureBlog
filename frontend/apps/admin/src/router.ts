import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from './stores/auth'

const Login = () => import('./pages/Login.vue')
const Manage = () => import('./pages/Manage.vue')
const Write = () => import('./pages/Write.vue')
const ManageTags = () => import('./pages/ManageTags.vue')
const ManageComments = () => import('./pages/ManageComments.vue')
const Settings = () => import('./pages/Settings.vue')

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/manage' },
    { path: '/login', component: Login },
    { path: '/manage', component: Manage, meta: { auth: true } },
    { path: '/write', component: Write, meta: { auth: true } },
    { path: '/write/:id', component: Write, meta: { auth: true } },
    { path: '/tags', component: ManageTags, meta: { auth: true } },
    { path: '/comments', component: ManageComments, meta: { auth: true } },
    { path: '/settings', component: Settings, meta: { auth: true } },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.auth && !auth.isAuthed) return '/login'
  return true
})
