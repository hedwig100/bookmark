import { createRouter, createWebHashHistory } from 'vue-router'
import LoginPage from './components/LoginPage.vue'
import HelloWorld from './components/HelloWorld.vue'
import UserCreate from './components/UserCreate.vue'
import TheHome from './components/TheHome.vue'

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
    },
    {
        path: "/users/:username",
        component: TheHome
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export { router }
