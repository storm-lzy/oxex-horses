import axios, { type AxiosInstance, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

const request: AxiosInstance = axios.create({
    baseURL: '/',
    timeout: 15000
})

request.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('admin_token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => Promise.reject(error)
)

request.interceptors.response.use(
    (response: AxiosResponse) => {
        const res = response.data
        if (res.code !== 0) {
            ElMessage.error(res.message || '请求失败')
            if (res.code === 401 || res.code === 10005 || res.code === 10006) {
                localStorage.removeItem('admin_token')
                window.location.href = '/login'
            }
            return Promise.reject(new Error(res.message))
        }
        return res.data
    },
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('admin_token')
            window.location.href = '/login'
        }
        ElMessage.error(error.message || '网络错误')
        return Promise.reject(error)
    }
)

export default request
