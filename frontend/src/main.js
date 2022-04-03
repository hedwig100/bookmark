import { createApp } from 'vue'
import App from './App.vue'
import { router } from './router.js'

// prepare app
const app = createApp(App)
app.use(router)
app.mount('#app')
