import { createApp } from 'vue'

import './style.css'

import App from './App.vue'
import router from '@/router/router.ts'

import ElementPlus from 'element-plus'

import 'element-plus/dist/index.css'
// 导入 ElementPlus 暗黑模式的样式
import 'element-plus/theme-chalk/dark/css-vars.css'

import { axiosInstance } from '@/api/instance'

// 测试用
if (import.meta.env.DEV) {
    // ✅ 仅在 npm run dev 时执行
    axiosInstance.defaults.baseURL = 'http://172.19.237.114:8080/api'
    axiosInstance.defaults.headers.common['Authorization'] = `Basic ${btoa('admin:admin')}`
}
// 生产环境
if (import.meta.env.PROD) {
    axiosInstance.defaults.baseURL = '/api'
}

const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.mount('#app')