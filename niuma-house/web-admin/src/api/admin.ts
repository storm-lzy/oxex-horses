import request from '@/utils/request'

// 登录
export const login = (data: { username: string; password: string }) => {
    return request.post('/api/auth/login', data)
}

// 获取统计数据
export const getDashboardStats = () => {
    return request.get('/api/admin/dashboard/stats')
}

// 获取用户列表
export const getUsers = (params?: { page?: number; size?: number }) => {
    return request.get('/api/admin/users', { params })
}

// 封禁用户
export const banUser = (id: number) => {
    return request.post(`/api/admin/users/${id}/ban`)
}

// 解封用户
export const unbanUser = (id: number) => {
    return request.post(`/api/admin/users/${id}/unban`)
}

// 获取帖子列表
export const getPosts = (params?: { page?: number; size?: number }) => {
    return request.get('/api/admin/posts', { params })
}

// 删除帖子
export const deletePost = (id: number) => {
    return request.delete(`/api/admin/posts/${id}`)
}

// 置顶帖子
export const topPost = (id: number) => {
    return request.post(`/api/admin/posts/${id}/top`)
}

// 获取公司列表
export const getCompanies = (params?: { page?: number; size?: number }) => {
    return request.get('/api/admin/companies', { params })
}

// 删除公司
export const deleteCompany = (id: number) => {
    return request.delete(`/api/admin/companies/${id}`)
}
