import { createApp } from 'vue'
import App from './App.vue'
import { router } from './router.js'
import axios from 'axios';

// prepare axios
axios.defaults.baseURL = 'http://localhost:8081';

// prepare app
const app = createApp(App)
app.use(router)
app.mount('#app')
