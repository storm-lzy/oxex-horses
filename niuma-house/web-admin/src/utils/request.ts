import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const axiosInstance = axios.create({
    baseURL: '/',
    timeout: 15000
})

axiosInstance.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('admin_token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => Promise.reject(error)
)

axiosInstance.interceptors.response.use(
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

// 包装请求方法，正确处理返回类型
const request = {
    get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
        return axiosInstance.get(url, config) as Promise<T>
    },
    post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        return axiosInstance.post(url, data, config) as Promise<T>
    },
    put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        return axiosInstance.put(url, data, config) as Promise<T>
    },
    delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
        return axiosInstance.delete(url, config) as Promise<T>
    },
    patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        return axiosInstance.patch(url, data, config) as Promise<T>
    }
}

export default request
