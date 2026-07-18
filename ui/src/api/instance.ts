// src/api/axiosInstance.ts
import axios from 'axios'

// 可以在这里改统一前缀，例如 '/api'、'/backend/api' 等
export const axiosInstance = axios.create()

// ✅ 请求拦截器（可选）
axiosInstance.interceptors.request.use(
    config => { return config },
    error => { return Promise.reject(error) }
)

// ✅ 响应拦截器（可选）
axiosInstance.interceptors.response.use(
    response => response,
    error => Promise.reject(error)
)