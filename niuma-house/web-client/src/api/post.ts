import request from '@/utils/request'
import type { User } from './user'

export interface Post {
    id: number
    user_id: number
    user?: User
    occupation_id: number
    occupation?: { id: number; name: string }
    title: string
    content: string
    likes_count: number
    views_count: number
    status: number
    created_at: string
}

export interface PostListResponse {
    list: Post[]
    total: number
    page: number
    size: number
}

export interface CreatePostRequest {
    title: string
    content: string
    occupation_id: number
}

// 获取帖子列表
export const getPosts = (params: { occupation_id?: number; page?: number; size?: number }): Promise<PostListResponse> => {
    return request.get('/posts', { params })
}

// 获取帖子详情
export const getPost = (id: number): Promise<{ post: Post; is_liked: boolean; is_favorited: boolean }> => {
    return request.get(`/posts/${id}`)
}

// 创建帖子
export const createPost = (data: CreatePostRequest): Promise<Post> => {
    return request.post('/posts', data)
}

// 更新帖子
export const updatePost = (id: number, data: Partial<CreatePostRequest>): Promise<void> => {
    return request.put(`/posts/${id}`, data)
}

// 删除帖子
export const deletePost = (id: number): Promise<void> => {
    return request.delete(`/posts/${id}`)
}

// 点赞帖子
export const likePost = (id: number): Promise<void> => {
    return request.post(`/posts/${id}/like`)
}

// 取消点赞
export const unlikePost = (id: number): Promise<void> => {
    return request.delete(`/posts/${id}/like`)
}

// 收藏帖子
export const favoritePost = (id: number): Promise<void> => {
    return request.post(`/posts/${id}/favorite`)
}

// 取消收藏
export const unfavoritePost = (id: number): Promise<void> => {
    return request.delete(`/posts/${id}/favorite`)
}

// 获取评论列表
export const getComments = (postId: number, params?: { page?: number; size?: number }): Promise<{ list: any[]; total: number }> => {
    return request.get(`/posts/${postId}/comments`, { params })
}

// 创建评论
export const createComment = (postId: number, data: { content: string; parent_id?: number }): Promise<any> => {
    return request.post(`/posts/${postId}/comments`, data)
}

// 删除评论
export const deleteComment = (id: number): Promise<void> => {
    return request.delete(`/comments/${id}`)
}
