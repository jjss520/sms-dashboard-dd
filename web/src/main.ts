import { createApp } from 'vue'
import axios from 'axios'
import './style.css'
import App from './App.vue'
import router from './router'

// 配置 API 基础 URL
// 双端口模式,使用相对路径,通过反代或同域访问
// 如果是跨域访问,需要配置 CORS
axios.defaults.baseURL = import.meta.env.VITE_API_URL || ''

createApp(App).use(router).mount('#app')
