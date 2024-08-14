// src/router.js
import { createRouter, createWebHistory } from 'vue-router'
import Chat from './pages/Chat.vue'
import Login from './pages/Login.vue'

const routes = [
    { path: '/', component: Login },
    { path: '/chat', component: Chat }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router