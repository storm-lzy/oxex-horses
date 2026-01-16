import axios, { type AxiosInstance, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const request: AxiosInstance = axios.create({
    baseURL: '/api',
    timeout: 15000
})

// 请求拦截器
request.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    (response: AxiosResponse) => {
        const res = response.data
        if (res.code !== 0) {
            ElMessage.error(res.message || '请求失败')

            // Token 过期或无效
            if (res.code === 401 || res.code === 10005 || res.code === 10006) {
                localStorage.removeItem('token')
                window.location.href = '/login'
            }

            return Promise.reject(new Error(res.message))
        }
        return res.data
    },
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token')
            window.location.href = '/login'
        }
        ElMessage.error(error.message || '网络错误')
        return Promise.reject(error)
    }
)

export default request
