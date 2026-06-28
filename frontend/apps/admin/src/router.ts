import { createRouter, createWebHistory } from 'vue-router'
import Login from './pages/Login.vue'
import Manage from './pages/Manage.vue'
import Write from './pages/Write.vue'
import ManageTags from './pages/ManageTags.vue'
import ManageComments from './pages/ManageComments.vue'
import Settings from './pages/Settings.vue'
import { useAuthStore } from './stores/auth'

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
