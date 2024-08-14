// src/router.js
import { createRouter, createWebHistory } from 'vue-router'
import Chat from './components/Chat.vue'
import Login from './components/Login.vue'

const routes = [
    { path: '/', component: Login },
    { path: '/chat', component: Chat }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router