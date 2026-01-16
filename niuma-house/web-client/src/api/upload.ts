import request from '@/utils/request'

// 获取上传预签名 URL
export const getPresignedURL = (data: { filename: string; file_type?: string }): Promise<{
    upload_url: string
    access_url: string
    object_key: string
}> => {
    return request.post('/upload/presign', data)
}
