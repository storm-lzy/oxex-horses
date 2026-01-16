import request from '@/utils/request'

export interface User {
    id: number
    username: string
    occupation_id: number
    occupation?: { id: number; name: string }
    level: number
    exp: number
    role: string
    status: number
    created_at: string
}

export interface LoginRequest {
    username: string
    password: string
}

export interface RegisterRequest {
    username: string
    password: string
    occupation_id: number
}

export interface LoginResponse {
    token: string
    user: User
}

// 登录
export const login = (data: LoginRequest): Promise<LoginResponse> => {
    return request.post('/auth/login', data)
}

// 注册
export const register = (data: RegisterRequest): Promise<User> => {
    return request.post('/auth/register', data)
}

// 获取当前用户资料
export const getProfile = (): Promise<User> => {
    return request.get('/user/profile')
}

// 更新用户资料
export const updateProfile = (data: { occupation_id: number }): Promise<void> => {
    return request.put('/user/profile', data)
}

// 获取职业列表
export const getOccupations = (): Promise<{ id: number; name: string }[]> => {
    return request.get('/occupations')
}
