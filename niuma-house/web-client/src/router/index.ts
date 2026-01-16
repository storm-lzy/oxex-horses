import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        component: () => import('@/views/Layout.vue'),
        children: [
            {
                path: '',
                name: 'Home',
                component: () => import('@/views/Home.vue'),
                meta: { title: '首页' }
            },
            {
                path: 'post/:id',
                name: 'PostDetail',
                component: () => import('@/views/PostDetail.vue'),
                meta: { title: '帖子详情', requiresAuth: true }
            },
            {
                path: 'post/create',
                name: 'CreatePost',
                component: () => import('@/views/CreatePost.vue'),
                meta: { title: '发布帖子', requiresAuth: true }
            },
            {
                path: 'companies',
                name: 'Companies',
                component: () => import('@/views/Companies.vue'),
                meta: { title: '坑逼公司墙' }
            },
            {
                path: 'company/:id',
                name: 'CompanyDetail',
                component: () => import('@/views/CompanyDetail.vue'),
                meta: { title: '公司详情' }
            },
            {
                path: 'company/create',
                name: 'CreateCompany',
                component: () => import('@/views/CreateCompany.vue'),
                meta: { title: '添加避雷公司', requiresAuth: true }
            },
            {
                path: 'messages',
                name: 'Messages',
                component: () => import('@/views/Messages.vue'),
                meta: { title: '私信', requiresAuth: true }
            },
            {
                path: 'profile',
                name: 'Profile',
                component: () => import('@/views/Profile.vue'),
                meta: { title: '个人中心', requiresAuth: true }
            }
        ]
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: { title: '登录' }
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('@/views/Register.vue'),
        meta: { title: '注册' }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
    document.title = `${to.meta.title || '牛马之家'} - 牛马之家`

    if (to.meta.requiresAuth) {
        const userStore = useUserStore()
        if (!userStore.isLoggedIn) {
            next({ name: 'Login', query: { redirect: to.fullPath } })
            return
        }
    }
    next()
})

export default router
