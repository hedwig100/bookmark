import { createRouter, createWebHashHistory } from 'vue-router'
import LoginPage from './components/LoginPage.vue'
import HelloWorld from './components/HelloWorld.vue'
import UserCreate from './components/UserCreate.vue'

const routes = [
    {
        path: "/",
        component: HelloWorld
    },
    {
        path: "/login",
        component: LoginPage
    },
    {
        path: "/users",
        component: UserCreate
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export { router }
