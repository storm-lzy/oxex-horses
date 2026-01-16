import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue')
    },
    {
        path: '/',
        component: () => import('@/views/Layout.vue'),
        redirect: '/dashboard',
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/Dashboard.vue'),
                meta: { title: '数据大屏' }
            },
            {
                path: 'users',
                name: 'Users',
                component: () => import('@/views/Users.vue'),
                meta: { title: '用户管理' }
            },
            {
                path: 'posts',
                name: 'Posts',
                component: () => import('@/views/Posts.vue'),
                meta: { title: '帖子管理' }
            },
            {
                path: 'companies',
                name: 'Companies',
                component: () => import('@/views/Companies.vue'),
                meta: { title: '公司管理' }
            }
        ]
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
    const token = localStorage.getItem('admin_token')
    if (to.name !== 'Login' && !token) {
        next({ name: 'Login' })
        return
    }
    next()
})

export default router
