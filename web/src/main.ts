import { createApp } from 'vue'
import axios from 'axios'
import './style.css'
import App from './App.vue'
import router from './router'

// 配置 API 基础 URL
axios.defaults.baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

createApp(App).use(router).mount('#app')
