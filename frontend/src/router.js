import { createRouter, createWebHashHistory } from 'vue-router'
import LoginPage from './components/LoginPage.vue'
import HelloWorld from './components/HelloWorld.vue'

const routes = [
    {
        path: "/",
        component: HelloWorld
    },
    {
        path: "/login",
        component: LoginPage
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export { router }
