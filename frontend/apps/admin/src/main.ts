import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { router } from './router'
import '@pureblog/ui/tokens.css'
import './styles/admin.css'

createApp(App).use(createPinia()).use(router).mount('#app')
