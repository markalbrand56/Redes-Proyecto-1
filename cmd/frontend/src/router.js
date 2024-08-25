// src/router.js
import { createRouter, createWebHistory } from 'vue-router'
import Chat from './pages/Chat.vue'
import Login from './pages/Login.vue'
import Register from "./pages/Register.vue";

const routes = [
    { path: '/', component: Login },
    {path: '/register', component: Register},
    { path: '/chat', component: Chat }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router