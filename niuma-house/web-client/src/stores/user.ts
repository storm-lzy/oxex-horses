import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, register, getProfile, type User, type LoginRequest, type RegisterRequest } from '@/api/user'

export const useUserStore = defineStore('user', () => {
    const token = ref<string>(localStorage.getItem('token') || '')
    const user = ref<User | null>(null)

    const isLoggedIn = computed(() => !!token.value)

    const setToken = (newToken: string) => {
        token.value = newToken
        localStorage.setItem('token', newToken)
    }

    const clearToken = () => {
        token.value = ''
        user.value = null
        localStorage.removeItem('token')
    }

    const loginAction = async (data: LoginRequest) => {
        const res = await login(data)
        setToken(res.token)
        user.value = res.user
        return res
    }

    const registerAction = async (data: RegisterRequest) => {
        const res = await register(data)
        return res
    }

    const fetchProfile = async () => {
        if (!token.value) return
        try {
            const res = await getProfile()
            user.value = res
        } catch (error) {
            clearToken()
        }
    }

    const logout = () => {
        clearToken()
    }

    // 等级名称映射
    const levelName = computed(() => {
        const names: Record<number, string> = {
            1: '普通牛马',
            2: '内卷牛马',
            3: '精英牛马',
            4: '天选牛马',
            5: '核动力牛马'
        }
        return names[user.value?.level || 1] || '普通牛马'
    })

    return {
        token,
        user,
        isLoggedIn,
        levelName,
        loginAction,
        registerAction,
        fetchProfile,
        logout
    }
})
