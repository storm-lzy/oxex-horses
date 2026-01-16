import request from '@/utils/request'
import type { User } from './user'

export interface Company {
    id: number
    name: string
    city: string
    tags: string[]
    risk_level: number
    evidence: string[]
    content: string
    creator_id: number
    creator?: User
    status: number
    view_count: number
    created_at: string
}

export interface CompanyListResponse {
    list: Company[]
    total: number
    page: number
    size: number
}

export interface CreateCompanyRequest {
    name: string
    city: string
    tags: string[]
    risk_level: number
    evidence: string[]
    content: string
}

// 获取公司列表
export const getCompanies = (params?: { page?: number; size?: number }): Promise<CompanyListResponse> => {
    return request.get('/companies', { params })
}

// 搜索公司
export const searchCompanies = (params: { keyword: string; page?: number; size?: number }): Promise<CompanyListResponse & { keyword: string }> => {
    return request.get('/companies/search', { params })
}

// 获取公司详情
export const getCompany = (id: number): Promise<Company> => {
    return request.get(`/companies/${id}`)
}

// 创建公司
export const createCompany = (data: CreateCompanyRequest): Promise<Company> => {
    return request.post('/companies', data)
}
