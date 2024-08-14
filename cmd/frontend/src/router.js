// src/router.js
import { createRouter, createWebHistory } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import Login from './components/Login.vue'

const routes = [
    { path: '/', component: Login },
    { path: '/hello', component: HelloWorld }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router